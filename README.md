# kubeStacker



## Description

## Usage

```
KST_CONFIG='config.json'
KUBECTL_PATH='/usr/bin/kubectl'
HELM_PATH='/usr/bin/helm'

./kubeStacker add --name=my.great.site.com  --type=wordpress --size=10G
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/valentin2105/kubeStacker
```

## Contribution

1. Fork ([https://github.com/valentin2105/kubeStacker/fork](https://github.com/valentin2105/kubeStacker/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[valentin2105](https://github.com/valentin2105)
