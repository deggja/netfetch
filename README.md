## Using Netfetch Tool

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, comparing them with a predefined standard to ensure they adhere to best practices. This document guides you on how to use `netfetch` to perform these scans.

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
