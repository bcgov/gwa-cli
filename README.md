# GWA CLI

<img src="https://github.com/bcgov/gwa-cli/workflows/Build/badge.svg"></img>
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gwa-cli&metric=alert_status)](https://sonarcloud.io/dashboard?id=gwa-cli)
[![img](https://img.shields.io/badge/Lifecycle-Stable-97ca00)](https://github.com/bcgov/repomountie/blob/master/doc/lifecycle-badges.md)

> **NOTE:** This tool is going through rapid development and could change often. Check back for updates.

`gwa` is a tool for composing, validating and generating Kong Gateway configuration files from OpenAPI (aka Swagger) specs and managing Kong Plugins.

## Documentation

Currently documentation is limited to the APS [USER-JOURNEY.md](https://github.com/bcgov/gwa-api/blob/dev/USER-JOURNEY.md) and via the installed executable's help commands. A hosted website option is in the works and will be available soon.

To run help on any command while using the CLI:

```sh
$ gwa-cli login --help
You can login via device login or by using client credentials

To use device login, simply run the command like so:
    $ gwa login

To use your credentials you must supply both a client-id and client-secret:
    $ gwa login --client-id <YOUR_CLIENT_ID> --client-secret <YOUR_CLIENT_SECRET>
...

```

To generate a Markdown-formatted string of every command's documentation to `stdout`, run `$ just docs`. You could then pipe it to a file like so `$ just docs > gwa-commands.md`, though this will be just a temporary solution.

## Development

Prerequisites:
- [Go](https://go.dev) 1.20 or higher
- [Just](https://github.com/casey/just) (alternative to `make`)
- [GoReleaser](https://goreleaser.com)

Tools:
- [Cobra] Command line argument parser
- [Viper](https://github.com/spf13/viper) Configuration file manager, integrates tightly with Cobra
- [Lipgloss](https://github.com/charmbracelet/lipgloss) Styles and colours

#### Steps to set up dev environment

1. Verify you have Go 1.20+ installed

   ```sh
   $ go version
   ```
   If you don't have `go` installed on your machine, follow instructions on [the Go website](https://golang.org/doc/install).

2. Clone this repository

   ```sh
   $ git clone git@github.com:bcgov/gwa-cli.git
   $ cd gwa-cli
   $ mv env.example .env
   ```
   **Note** Some local environments require Go projects are run from the `$HOME/go/src` directory. If any `module not found` errors are reported, try moving it.

   Also be sure to fill in the environment variables before running a command.

3. Run commands

   Test any commmands by running `just run` in the `cwd`. You can also use `$ just test` to run all tests.

   ```sh
   $ just run namespace current
   your-namespace
   $ just test
   ?   	github.com/bcgov/gwa-cli	[no test files]
   ok  	github.com/bcgov/gwa-cli/cmd
   ok  	github.com/bcgov/gwa-cli/pkg
   ```

4. Set up your IDE

   Go has great tooling which is required to ensure code contributed is formatted consitently and is type-safe.

   - **VSCode:** Install the [official Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
   - **NeoVim:** [go.nvim](https://github.com/ray-x/go.nvim) is a great plugin

#### Nix develop (optional)

Instead of installing Go on your system you can alternatively use Nix to manage Go and other dependencies for you. This option assumes you already running NixOS or have Nix installed.

Simply run `$ nix develop` in the project root. This will pull all the dependencies needed to work on this project. Remember to run `exit` to leave the Nix shell when done.

## Deployment

Details on how the deployment process work are coming soon as we migrate codebase languages.

Make sure you have GoReleaser installed to test locally.

###### Tip

A dry-run of the build process can be executed by running `$ just release`. This will compile your source code into a `/dist` folder where you can run any platform's executable.

## Installation

Currently `gwa` is only installable by building from source. The releases page will be updated one it's ready.

To install locally you can follow the first 2 steps in Development above, then run

```sh
$ just install
...
$  gwa-cli
gwa version 2.0.0-beta
```

#### Completions

Shell completions for all the commands ships with each version. Completions allow you to tab while entering commands to cycle though a list of possible commands.

To install completions, run this after installing, using `zsh` for example:

```sh
$ gwa-cli completion zsh --help
$ gwa-cli completion zsh | pbcopy
```

Then follow the instructions from the help command and paste the output where it needs to live. Bash, Fish and Powershell are also supported.
