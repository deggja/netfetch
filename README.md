![image](https://github.com/deggja/netfetch/assets/15778492/b9a93dce-a09a-4823-be99-dcda5dbf6dc7)

## Using Netfetch Tool

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, checking if there are implicit defautl deny policies in place or any other policy targetting the pods.

This document guides you on how to use `netfetch` to perform these scans.

## Contribute
You are welcome to contribute!
 
Open an issue or create a pull request if there is some functionality missing that you would like.

## Installation via Homebrew

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

### Prerequisites

Before you begin, ensure you have the following:

- `netfetch` binary installed in your system.
- Access to a Kubernetes cluster with configured `kubectl`.
- Sufficient permissions to list namespaces and network policies in the cluster.

### Basic Usage

The primary command provided by `netfetch` is `scan`. This command scans all non-system Kubernetes namespaces for network policies.

#### Command Structure

```sh
netfetch scan
```

You can also specifiy namespaces when running netfetch.

```sh
netfetch scan default
```