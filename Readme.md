### CNDRO conference kubectl plugin

* A `kubectl` plugin for [CloudNativeDays_RO](https://cloudnativedays.ro/)

### Install

* Using **Krew** custom index [dejanu/kubectl-cndro](https://github.com/dejanu/kubectl-cndro):

```bash
kubectl krew index add dejanu https://github.com/dejanu/kubectl-cndro.git
kubectl krew update
kubectl krew install dejanu/cndro
```

* Check the plugin is visible: `kubectl krew list`. You should see `cndro` in the list.

### Usage

```bash
# show help/usage
kubectl cndro -h

# show day1 schedule
kubectl cndro day1

# show day2 schedule
kubectl cndro day2

# show tickets URL and pricing table
kubectl cndro tickets
```