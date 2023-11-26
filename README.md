![image](https://github.com/deggja/netfetch/assets/15778492/b9a93dce-a09a-4823-be99-dcda5dbf6dc7)

## Using Netfetch

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, checking whether implicit default deny policies are in place and examining if there are any other policies targeting the pods.

This document guides you on how to use `netfetch` to perform these scans.

## Installation via Homebrew for Mac

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

For specific Linux distros, Windows and other install binaries, check the latest release.

### Prerequisites

Before you begin, ensure you have the following:

- `netfetch` binary installed in your system.
- Access to a Kubernetes cluster with configured `kubectl`.
- Permissions to read and create network policies in at least one namespace.

### Usage

The primary command provided by `netfetch` is `scan`. This command scans all non-system Kubernetes namespaces for network policies.

You can also scan specific namespaces by specifying the name of that namespace.

Scan entire cluster.

```sh
netfetch scan
```

Scan a namespace called production.

```sh
netfetch scan production
```

Launch dashboard.

```sh
netfetch dash
```

![Netfetch Dashboard](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/netfetch_dash.png)

## Netfetch score

The `netfetch` tool provides a basic score at the end of each scan. The score ranges from 1 to 42, with 1 being the lowest and 42 being the highest possible score.

This score reflects the security posture of your Kubernetes namespaces based on network policies and general policy coverage. If changes are made based on recommendations from the initial scan, rerunning `netfetch` will likely result in a higher score.

## Contribute
You are welcome to contribute!

1. Fork the Project
2. Create your Feature Branch (git checkout -b feature/AmazingFeature)
3. Commit your Changes (git commit -m 'Add some AmazingFeature')
4. Push to the Branch (git push origin feature/AmazingFeature)
5. Open a Pull Request

## License

Netfetch is distributed under the MIT License. See the LICENSE file for more information. See the [LICENSE](LICENSE) for more information.
