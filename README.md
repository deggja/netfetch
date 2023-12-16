![image](https://github.com/deggja/netfetch/assets/15778492/b9a93dce-a09a-4823-be99-dcda5dbf6dc7)

## Contents
- [What is this project?](#⭐-what-is-this-project-⭐)
  - [Install with brew](#installation-via-homebrew-for-mac-💻)
  - [Install in Kubernetes](#installation-via-helm-🎩)
  - [How to use](#usage)
  - [Dashboard](#using-the-dashboard-📟)
  - [Score](#netfetch-score)
  - [Uninstalling](#uninstalling-netfetch)
- [Contribute](#contribute-🔨)

## ⭐ What is this project ⭐

This project aims to demystify network policies in Kubernetes. It's a work in progress!

The `netfetch` tool is designed to scan Kubernetes namespaces for network policies, checking if your workloads are targeted by a network policy or not.

What can I use `netfetch` for? 🤔

CLI:
- Scan your Kubernetes cluster or namespace to identify pods running with no ingress and egress restrictions.
- Save the output of your scans in a text file to analyze.
- Create implicit default deny network policies in namespaces that do not have one.
- Get a score calculated for your cluster or namespace based on the findings of the scans.

Dashboard:
- Scan your cluster or namespace and list pods running without network restrictions in a table.
- Visualise all existing network policies and pods in your cluster or namespace in a network map you can interact with.
- Double click a network policy in a network map to preview the YAML of that policy.
- Create implicit default deny network policies in namespaces that do not have one.
- Get suggestions for network policies that you can edit & apply to your namespaces by analysing existing pods.
- Get a score calculated for your cluster or namespace based on the findings of the scans.

## Installation via Homebrew for Mac 💻

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

For specific Linux distros, Windows and other install binaries, check the latest release.

## Installation via Helm 🎩

You can deploy the `netfetch` dashboard in your Kubernetes clusters using Helm.

```sh
helm repo add deggja https://deggja.github.io/netfetch/
helm repo update
helm install netfetch deggja/netfetch --namespace netfetch --create-namespace
```

Follow the instructions after deployment to access the dashboard.

### Prerequisites 🌌

- Installed `netfetch` via homebrew or a release binary.
- Access to a Kubernetes cluster with `kubectl` configured.
- Permissions to read and create network policies.

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

### Using the dashboard 📟

Launch the dashboard:

```sh
netfetch dash
```

While in the dashboard, you have a couple of options.

You can use the `Scan cluster` button, which is the equivalent to the CLI `netfetch scan` command. This will populate the table view with all pods not targeted by a network policy.

Scanning a specific namespace is done by selecting the namespace of choice from the `Select a namespace` dropdown and using the `Scan namespace` button. This is the equivalent to the CLI `netfetch scan namespace` command. 

This will populate the table view with all pods not targeted by a network policy in that specific namespace. In addition to this, if there are any pods in the cluster already targeted by a network policy - it will create a visualisation of this in a network map rendered using [D3](https://d3-graph-gallery.com/network.html) below the table view.

![Netfetch Dashboard](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/netfetch_new_dash.png)

You can click the `Create cluster map` button to do exactly that. This will render a network map with D3, fetching all pods and policies in all the namespaces you have access to in the cluster.

Inside the network map visualisations, you can double click the network policy nodes to preview the YAML of that policy.

![Network map](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/netfetch_network_map.png)

## Netfetch score 🥇

The `netfetch` tool provides a basic score at the end of each scan. The score ranges from 1 to 42, with 1 being the lowest and 42 being the highest possible score.

As of today, your score will decrease if you are missing implicit default deny all network policies in your namespace or cluster. It will also decrease based on the amount of pods not targeted by a network policy.

The score reflects the security posture of your Kubernetes namespaces based on network policies and general policy coverage. If changes are made based on recommendations from the initial scan, rerunning `netfetch` will likely result in a higher score.

## Uninstalling netfetch

If you want to uninstall the application - you can do so by running the following commands.

```
brew uninstall netfetch
brew cleanup -s netfetch
brew untap deggja/netfetch https://github.com/deggja/netfetch
```

## Contribute 🔨
You are welcome to contribute!

See [CONTRIBUTING](CONTRIBUTING.md) for instructions on how to proceed.

## Tools 🧰
Netfetch uses other tools for a plethora of different things. It would not be possible without the following:

- [statik](https://github.com/rakyll/statik)
- [D3](https://d3-graph-gallery.com/network.html)
- [Helm](https://helm.sh/docs/)
- [Brew](https://brew.sh/)

## License

Netfetch is distributed under the MIT License. See the [LICENSE](LICENSE) for more information.
