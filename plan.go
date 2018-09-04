package main

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// orderedDecision type representing a Decision and it's priority weight
type orderedDecision struct {
	Description string
	Priority    int
}

// orderedCommand type representing a Command and it's priority weight and the targeted release from the desired state
type orderedCommand struct {
	Command       command
	Priority      int
	targetRelease *release
}

// plan type representing the plan of actions to make the desired state come true.
type plan struct {
	Commands  []orderedCommand
	Decisions []orderedDecision
	Created   time.Time
}

// createPlan initializes an empty plan
func createPlan() plan {

	p := plan{
		Commands:  []orderedCommand{},
		Decisions: []orderedDecision{},
		Created:   time.Now().UTC(),
	}
	return p
}

// addCommand adds a command type to the plan
func (p *plan) addCommand(cmd command, priority int, r *release) {
	oc := orderedCommand{
		Command:       cmd,
		Priority:      priority,
		targetRelease: r,
	}

	p.Commands = append(p.Commands, oc)
}

// addDecision adds a decision type to the plan
func (p *plan) addDecision(decision string, priority int) {
	od := orderedDecision{
		Description: decision,
		Priority:    priority,
	}
	p.Decisions = append(p.Decisions, od)
}

// execPlan executes the commands (actions) which were added to the plan.
func (p plan) execPlan() {
	p.sortPlan()
	log.Println("INFO: Executing the plan ... ")
	for _, cmd := range p.Commands {
		if exitCode, msg := cmd.Command.exec(debug, verbose); exitCode != 0 {
			logError("Command returned with exit code: " + string(exitCode) + ". And error message: " + msg)
		} else {
			if cmd.targetRelease != nil {
				labelResource(cmd.targetRelease)
			}
			if _, err := url.ParseRequestURI(s.Settings.SlackWebhook); err == nil {
				notifySlack(cmd.Command.Description+" ... SUCCESS!", s.Settings.SlackWebhook, false, true)
			}
		}
	}
}

// printPlanCmds prints the actual commands that will be executed as part of a plan.
func (p plan) printPlanCmds() {
	fmt.Println("Printing the commands of the current plan ...")
	for _, cmd := range p.Commands {
		fmt.Println(cmd.Command.Args[1])
	}
}

// printPlan prints the decisions made in a plan.
func (p plan) printPlan() {
	fmt.Println("----------------------")
	log.Printf("INFO: Plan generated at: %s \n", p.Created.Format("Mon Jan _2 2006 15:04:05"))
	for _, decision := range p.Decisions {
		fmt.Println(decision.Description + " -- priority: " + strconv.Itoa(decision.Priority))
	}
}

// sendPlanToSlack sends the description of plan commands to slack if a webhook is provided.
func (p plan) sendPlanToSlack() {
	if _, err := url.ParseRequestURI(s.Settings.SlackWebhook); err == nil {
		str := ""
		for _, c := range p.Commands {
			str = str + c.Command.Description + "\n"
		}

		notifySlack(strings.TrimRight(str, "\n"), s.Settings.SlackWebhook, false, false)
	}

}

// sortPlan sorts the slices of commands and decisions based on priorities
// the lower the priority value the earlier a command should be attempted
func (p plan) sortPlan() {
	log.Println("INFO: sorting the commands in the plan based on priorities (order flags) ... ")

	sort.SliceStable(p.Commands, func(i, j int) bool {
		return p.Commands[i].Priority < p.Commands[j].Priority
	})

	sort.SliceStable(p.Decisions, func(i, j int) bool {
		return p.Decisions[i].Priority < p.Decisions[j].Priority
	})
}
