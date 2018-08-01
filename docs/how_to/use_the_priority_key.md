---
version: v1.3.0-rc
---

# Using the priority key for Apps

The `priority` flag in Apps definition allows you to define the order at which apps operations will be applied. This is useful if you have dependencies between your apps/services. 

Priority is an optional flag and has a default value of 0 (zero). If set, it can only use a negative value. The lower the value, the higher the priority. 

## Example

```toml
[metadata]
org = "example.com"
description = "example Desired State File for demo purposes."


[settings]
kubeContext = "minikube" 

[namespaces]
  [namespaces.staging]
  protected = false
  [namespaces.production]
  prtoected = true

[helmRepos]
stable = "https://kubernetes-charts.storage.googleapis.com"
incubator = "http://storage.googleapis.com/kubernetes-charts-incubator"


[apps]

    [apps.jenkins]
    name = "jenkins" # should be unique across all apps
    description = "jenkins"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "stable/jenkins" # changing the chart name means delete and recreate this chart
    version = "0.14.3" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -2

    [apps.jenkins1]
    name = "jenkins1" # should be unique across all apps
    description = "jenkins"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "stable/jenkins" # changing the chart name means delete and recreate this chart
    version = "0.14.3" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    

    [apps.jenkins2]
    name = "jenkins2" # should be unique across all apps
    description = "jenkins"
    namespace = "production" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "stable/jenkins" # changing the chart name means delete and recreate this chart
    version = "0.14.3" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -3

    [apps.artifactory]
    name = "artifactory" # should be unique across all apps
    description = "artifactory"
    namespace = "staging" # maps to the namespace as defined in environments above
    enabled = true # change to false if you want to delete this app release [empty = false]
    chart = "stable/artifactory" # changing the chart name means delete and recreate this chart
    version = "7.0.6" # chart version
    valuesFile = "" # leaving it empty uses the default chart values
    priority= -2
```

The above example will generate the following plan:

```
DECISION: release [ jenkins2 ] is not present in the current k8s context. Will install it in namespace [[ production ]] -- priority: -3
DECISION: release [ jenkins ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: -2
DECISION: release [ artifactory ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: -2
DECISION: release [ jenkins1 ] is not present in the current k8s context. Will install it in namespace [[ staging ]] -- priority: 0

```