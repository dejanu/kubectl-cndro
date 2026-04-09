### CNDRO conference kubectl plugin

A `kubectl` plugin for [CloudNativeDays_RO](https://cloudnativedays.ro/)

The binary must be named `kubectl-cndro` and be on your `PATH` so `kubectl cndro` invokes it ([kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)).

#### Install

**Krew** (custom index — this repo [dejanu/kubectl-cndro](https://github.com/dejanu/kubectl-cndro)):

```bash
kubectl krew index add dejanu https://github.com/dejanu/kubectl-cndro.git
kubectl krew update
kubectl krew install dejanu/cndro
```

Requires a [GitHub Release](https://github.com/dejanu/kubectl-cndro/releases) for the tag in [plugins/cndro.yaml](plugins/cndro.yaml). Maintainer steps: [docs/KREW.md](docs/KREW.md). (Optional later: submit the same manifest to [krew-index](https://github.com/kubernetes-sigs/krew-index) for `kubectl krew install cndro` without a custom index.)

**Go**

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
