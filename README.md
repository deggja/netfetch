![image](https://github.com/deggja/netfetch/assets/15778492/b9a93dce-a09a-4823-be99-dcda5dbf6dc7)

## Contents
- [What is this project?](#using-netfetch)
  - [Installation via Homebrew for Mac](#installation-via-homebrew-for-mac)
  - [How to use](#usage)
  - [Dashboard](#dashboard)
  - [Score](#netfetch-score)
  - [Updating](#update-netfetch)
  - [Uninstalling](#uninstalling-netfetch)
- [Contribute](#contribute)

## What is this project?

This project aims to simplify the mapping of network policies in a Kubernetes cluster. It's a work in progress!

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, checking whether implicit default deny policies are in place and examining if there are any other policies targeting the pods.

This document guides you on how to use `netfetch` to perform these scans.

## Installation via Homebrew for Mac

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

For specific Linux distros, Windows and other install binaries, check the latest release.

## Update netfetch

If you are running an older version of the application - you can update `netfetch` by running the following commands.

```sh
brew update
netfetch version # verify that you are running the latest version
```

This will fetch the latest updates from the Homebrew tap.

If you are experiencing issues trying to fetch the latest version, you can run the following to commands:

```sh
brew uninstall netfetch
brew cleanup -s netfetch
brew update
brew install netfetch
netfetch version # verify that you are running the latest version
```

This should clean up any traces of the old version, update the Homebrew tap and refresh the binary.

## Uninstalling netfetch

If you want to uninstall the application - you can do so by running the following commands.

```
brew uninstall netfetch
brew cleanup -s netfetch
brew untap deggja/netfetch https://github.com/deggja/netfetch
```

### Prerequisites

Before you begin, ensure you have the following:

- `netfetch` binary installed in your system.
- Access to a Kubernetes cluster with configured `kubectl`.
- Permissions to read and create network policies in at least one namespace.

### Usage

The primary command provided by `netfetch` is `scan`. This command scans all non-system Kubernetes namespaces for network policies.

You can also scan specific namespaces by specifying the name of that namespace.

You may add the --dryrun or -d flag to run a dryrun of the scan. The application will not prompt you about adding network policies, but still give you the output of the scan.

Run `netfetch` in dryrun against a cluster.

```sh
netfetch scan --dryrun
```

Run `netfetch` in dryrun against a namespace

```sh
netfetch scan production --dryrun
```

Scan entire cluster.

```sh
netfetch scan
```

Scan a namespace called production.

```sh
netfetch scan production
```

### Dashboard

Launch dashboard.

```sh
netfetch dash
```

While in the dashboard, you have a couple of options.

You can use the `Scan cluster` button, which is the equivalent to the CLI `netfetch scan` command. This will populate the table view with all pods not targeted by a network policy.

Scanning a specific namespace is done by selecting the namespace of choice from the `Select a namespace` dropdown and using the `Scan namespace` button. This is the equivalent to the CLI `netfetch scan namespace` command. 

This will populate the table view with all pods not targeted by a network policy in that specific namespace. In addition to this, if there are any pods in the cluster already targeted by a network policy - it will create a visualisation of this in a network map rendered using [D3](https://d3-graph-gallery.com/network.html) below the table view.

You can click the `Create cluster map` button to do exactly that. This will render a network map with D3, fetching all pods and policies in all the namespaces you have access to in the cluster.

Inside the network map visualisations, you can double click the network policy nodes to preview the YAML of that policy.

![Netfetch Dashboard](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/netfetch_new_dash.png)


## Netfetch score

The `netfetch` tool provides a basic score at the end of each scan. The score ranges from 1 to 42, with 1 being the lowest and 42 being the highest possible score.

As of today, your score will decrease if you are missing implicit default deny all network policies in your namespace or cluster. It will also decrease based on the amount of pods not targeted by a network policy.

The score reflects the security posture of your Kubernetes namespaces based on network policies and general policy coverage. If changes are made based on recommendations from the initial scan, rerunning `netfetch` will likely result in a higher score.

## Contribute
You are welcome to contribute!

1. Fork the Project
2. Create your Feature Branch (git checkout -b feature/AmazingFeature)
3. Commit your Changes (git commit -m 'Add some AmazingFeature')
4. Push to the Branch (git push origin feature/AmazingFeature)
5. Open a Pull Request

## License

Netfetch is distributed under the MIT License. See the LICENSE file for more information. See the [LICENSE](LICENSE) for more information.
