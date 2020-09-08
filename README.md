# gwa-cli

> **NOTE:** This tool is going through rapid development and could change often. Check back for updates.

GWA CLI is a tool for composing, validating and generating Kong Gateway configuration files for OpenAPI specs and Kong Plugins.

## Installation

#### Prerequisites

- Node 12+
- npm (latest, comes installed with most node installers)

A Github Package that can easily be installed like an npm package is coming soon, in the meanttime see [Development](#development) below on how to run this tool.

## Usage

Create a new config file:

```bash
$ cd /to/an/empty/dir
$ gwa
```

Update an existing config

```bash
$ cd /to/a/config/dir
$ gwa medications.yaml
```

The `gwa` command will run the configuration wizard.


## Commands

##### Quit

Key: `ctrl + c`

This works the same as any command-line runtime

##### Export

Key: `ctrl + y`

This command will try to export the config as is.

##### Go Forward (if in history)

Key: `ctrl + k`

##### Go Back

Key: `ctrl + j`

#### Editor commands

##### Enable/Disable plugin

Key: `ctrl + e`

##### Save plugin config

Key: `ctrl + s`

##### Go to Next Plugin Config (when viewing config)

Key: `ctrl + n`

##### Go to Previous Plugin Config (when viewing config)

Key: `ctrl + p`


## Development

To install checkout this repo in the `dev` branch, then run the following:

```bash
$ npm i
$ npm run build
$ npm link
```

To run the TypeScript compiler, in another terminal run `$ npm start`.

To uninstall simply run `$ npm uninstall` from this directory.
