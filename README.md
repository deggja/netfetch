<div align="center">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-v1.21-brightgreen.svg" alt="go version">
  </a>
  <a href="https://d3js.org/">
    <img src="https://img.shields.io/badge/D3-v7.8.5-brightgreen.svg" alt="d3 version">
  </a>
  <a href="https://nodejs.org/">
    <img src="https://img.shields.io/badge/Node-v21.5.0-brightgreen.svg" alt="node version">
  </a>
  <a href="https://vuejs.org/">
    <img src="https://img.shields.io/badge/Vue-v5.0.8-brightgreen.svg" alt="vue version">
  </a>
</div>


<div align="center">

  <h1>Netfetch</h1>
  <h3>Scan your Kubernetes clusters to identifiy unprotected workloads and map your existing Network policies</h3>
  
  <img src="https://github.com/deggja/netfetch/assets/17279882/64306f2f-abbf-462c-97d6-a326ca70c2ad" width="130px" alt="Netfetch"/>

</div>

## Contents
- [**What is this project?**](#-what-is-this-project-)
  - [Support](#networkpolicy-type-support-in-netfetch)
- **[Installation](#installation)**
  - [Install with brew](#installation-via-homebrew-for-mac-)
  - [Install in Kubernetes](#installation-via-helm-)
- [**Usage**](#usage)
  - [Get started](#get-started)
  - [Dashboard](#using-the-dashboard-)
  - [Score](#netfetch-score-)
  - [Uninstalling](#uninstalling-netfetch)
- [**Contribute**](#contribute-)

## ⭐ What is this project ⭐

This project aims to demystify network policies in Kubernetes. It's a work in progress!

The `netfetch` tool will scan your Kubernetes cluster and let you know if you have any pods running without being targeted by network policies.

| Feature                                                                | CLI  | Dashboard |
|------------------------------------------------------------------------|------|-----------|
| Scan cluster identify pods without network policies                    | ✓    | ✓         |
| Save scan output to a text file                                        | ✓    |           |
| Visualize network policies and pods in a interactive network map       |      | ✓         |
| Create default deny network policies where this is missing             | ✓    | ✓         |
| Get suggestions for network policies based on existing workloads       |      | ✓         |
| Calculate a security score based on scan findings                      | ✓    | ✓         |
| Scan a specific policy by name to see what pods it  targets            | ✓    |           |

### NetworkPolicy type support in Netfetch

| Type      | CLI  | Dashboard |
|-----------|------|-----------|
| Kubernetes| ✓    | ✓         |
| Cilium    | ✓    |           |

Support for additional types of network policies is in the works. No support for the type you need? Check out [issues](https://github.com/deggja/netfetch/issues) for an existing request or create a new one if there is none.

## Installation
### Installation via Homebrew for Mac 💻

You can install `netfetch` using our Homebrew tap:

```sh
brew tap deggja/netfetch https://github.com/deggja/netfetch
brew install netfetch
```

For specific Linux distros, Windows and other install binaries, check the latest release.

### Installation via Helm 🎩

You can deploy the `netfetch` dashboard in your Kubernetes clusters using Helm.

```sh
helm repo add deggja https://deggja.github.io/netfetch/
helm repo update
helm install netfetch deggja/netfetch --namespace netfetch --create-namespace
```

Follow the instructions after deployment to access the dashboard.

#### Prerequisites 🌌

- Installed `netfetch` via homebrew or a release binary.
- Access to a Kubernetes cluster with `kubectl` configured.
- Permissions to read and create network policies.

## Usage

### Get started

The primary command provided by `netfetch` is `scan`. This command scans all non-system Kubernetes namespaces for network policies.

You can also scan specific namespaces by specifying the name of that namespace.

You may add the --dryrun or -d flag to run a dryrun of the scan. The application will not prompt you about adding network policies, but still give you the output of the scan.

Run `netfetch` in dryrun against a cluster.

```sh
netfetch scan --dryrun
```

Run `netfetch` in dryrun against a namespace

```sh
netfetch scan crossplane-system --dryrun
```

![netfetch-demo](https://github.com/deggja/netfetch/assets/15778492/015e9d9f-a678-4a14-a8bd-607f02c13d9f)

Scan entire cluster.

```sh
netfetch scan
```

Scan a namespace called crossplane-system.

```sh
netfetch scan crossplane-system
```

Scan entire cluster for Cilium Network Policies and or Cluster Wide Cilium Network Policies.

```sh
netfetch scan --cilium
```

Scan a namespace called production for regular Cilium Network Policies.

```sh
netfetch scan production --cilium
```

Scan a specific network policy.

```sh
netfetch scan --target my-policy-name
```

Scan a specific Cilium Network Policy.

```sh
netfetch scan --cilium --target default-cilium-default-deny-all
```

[![asciicast](https://asciinema.org/a/661200.svg)](https://asciinema.org/a/661200)

### Using the dashboard 📟

Launch the dashboard:

```sh
netfetch dash
```

You may also specify a port for the dashboard to run on (default is 8080).

```sh
netfetch dash --port 8081
```

### Dashboard functionality overview

The Netfetch Dashboard offers an intuitive interface for interacting with your Kubernetes cluster's network policies. Below is a detailed overview of the functionalities available through the dashboard:

| Action               | Description                                                                                                     | Screenshot Link                                                 |
|----------------------|-----------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------|
| Scan Cluster         | Initiates a cluster-wide scan to identify pods without network policies, similar to `netfetch scan`.            | ![Netfetch Dashboard](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/new-dash.png) |
| Scan Namespace       | Scans a selected namespace for pods not covered by network policies, equivalent to `netfetch scan namespace`.   | ![Cluster map](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/new-clustermap.png) |
| Create Cluster Map   | Generates a D3-rendered network map of all pods and policies across accessible namespaces.                      | ![Network map](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/new-ns.png) |
| Suggest Policy       | Provides network policy suggestions based on existing workloads within a selected namespace.                     | ![Suggested policies](https://github.com/deggja/netfetch/blob/main/frontend/dash/src/assets/new-suggestpolicy.png) |

### Interactive Features

- **Table View**: Shows pods not targeted by network policies. It updates based on the cluster or namespace scans.
- **Network Map Visualization**: Rendered using D3 to show how pods and policies interact within the cluster.
- **Policy Preview**: Double-click network policy nodes within the network map to view policy YAML.
- **Policy Editing**: Edit suggested policies directly within the dashboard or copy the YAML for external use.


### Netfetch score 🥇

The `netfetch` tool provides a basic score at the end of each scan. The score ranges from 1 to 42, with 1 being the lowest and 42 being the highest possible score.

Your score will decrease based on the amount of workloads in your cluster that are running without being targeted by a network policy.

The score reflects the security posture of your Kubernetes namespaces based on network policies and general policy coverage. If changes are made based on recommendations from the initial scan, rerunning `netfetch` will likely result in a higher score.

### Uninstalling netfetch

If you want to uninstall the application - you can do so by running the following commands.

```
brew uninstall netfetch
brew cleanup -s netfetch
brew untap deggja/netfetch https://github.com/deggja/netfetch
```

## Running Tests

To run tests for netfetch, follow these steps:

1. Navigate to the root directory of the project in your terminal.

2. Navigate to the backend directory within the project:

```
cd backend
```

3. Run the following command to execute all tests in the project:

```
go test ./...
```

This command will recursively search for tests in all subdirectories (./...) and run them.

4. After executing the command, you will see the test results in the terminal output.

## Contribute 🔨
Thank you to the following awesome people:

- [roopeshsn](https://github.com/roopeshsn)
- [s-rd](https://github.com/s-rd)
- [JJGadgets](https://github.com/JJGadgets)
- [Home Operations Discord](https://github.com/onedr0p/home-ops)
- [pehlicd](https://github.com/pehlicd)


You are welcome to contribute!

See [CONTRIBUTING](CONTRIBUTING.md) for instructions on how to proceed.

## Tools 🧰
Netfetch uses other tools for a plethora of different things. It would not be possible without the following:

- [statik](https://github.com/rakyll/statik)
- [D3](https://d3-graph-gallery.com/network.html)
- [Helm](https://helm.sh/docs/)
- [Brew](https://brew.sh/)
- [lipgloss](https://github.com/charmbracelet/lipgloss)

## License

Netfetch is distributed under the MIT License. See the [LICENSE](LICENSE) for more information.
