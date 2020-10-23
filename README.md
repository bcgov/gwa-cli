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

### `gwa init`

Generates a `.env` file in the current working directory.

A `.env` file should have the following key/values

```
GWA_NAMESPACE=<your namespace>
CLIENT_ID=<gwa client id>
CLIENT_SECRET=<gwa client secret>
GWA_ENV=<dev, prod or test>
```

To create and work with configurations you don't require `CLIENT_ID` or `CLIENT_SECRET`, but to make any API requests you will.

### `gwa new <input file>`

Initialize a config file in the current directory. The input file must be an OpenAPI JSON file or URL

### `gwa validate <input file>`

Validate a configuration file

### `gwa update <input file>`

Update a config with new OpenAPI specs

### `gwa plugins`

List all available plugins

### `gwa publish-gateway <config file>` Alias `pg`

Publish all YAML config files in current directory. Make sure your `.env` file is configured correctly.

### `gwa acl`

Update the full membership. Note that this command will overwrite the remote list of users, use with caution

## Help

Run `$ gwa --help` to see all available commands, `$ gwa <command> --help` to view an individual command's help content.

## Development

To install checkout this repo in the `dev` branch, then run the following:

```bash
$ npm i
$ npm run build
$ npm link
```

To run the TypeScript compiler, in another terminal run `$ npm start`.

To uninstall simply run `$ npm uninstall` from this directory.
