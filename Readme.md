### CNDRO conference kubectl plugin

A `kubectl` plugin for [CloudNativeDays_RO](https://cloudnativedays.ro/)

The binary must be named `kubectl-cndro` and be on your `PATH` so `kubectl cndro` invokes it ([kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)).

#### Install

```bash
go install ./cmd/kubectl-cndro

# published version:
go install github.com/dejanu/cndro/cmd/kubectl-cndro@latest
```


* Verify the plugin is visible:

```bash
go build -o kubectl-cndro ./cmd/kubectl-cndro
export PATH="$PATH:$(pwd)"
kubectl plugin list
```

You should see `cndro` in the list.

#### Usage

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
