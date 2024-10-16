package main

import (
	"os"

	"github.com/IrekArcher/helmsman/v3/internal/app"
)

func main() {
	exitCode := app.Main()
	os.Exit(exitCode)
}
