![image](https://github.com/deggja/netfetch/assets/15778492/b9a93dce-a09a-4823-be99-dcda5dbf6dc7)

## Using Netfetch Tool

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, checking whether implicit default deny policies are in place and examining any other policies targeting the pods.

This document guides you on how to use `netfetch` to perform these scans.

## Contribute
You are welcome to contribute!
 
Open an issue or create a pull request if there is some functionality missing that you would like.

## Installation via Homebrew for Mac

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

For specific Linux distros, Windows etc. Check the latest release for a downloadable binary.

### Prerequisites

Before you begin, ensure you have the following:

- `netfetch` binary installed in your system.
- Access to a Kubernetes cluster with configured `kubectl`.
- Sufficient permissions to list namespaces and network policies in the cluster.

### Basic usage

The primary command provided by `netfetch` is `scan`. This command scans all non-system Kubernetes namespaces for network policies.

#### Command structure

```sh
netfetch scan
```

You can also specifiy namespaces when running netfetch.

```sh
netfetch scan default
```

## Netfetch score

The `netfetch` tool provides a score at the end of each scan. The score ranges from 1 to 42, with 1 being the lowest and 42 being the highest possible score.

This score reflects the security posture of your Kubernetes namespaces based on network policies and pod coverage. If changes are made based on recommendations from the initial scan, rerunning `netfetch` will likely result in an improved score.

## License

[MIT License], see [LICENSE](LICENSE).