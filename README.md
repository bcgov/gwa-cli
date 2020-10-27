# gwa-cli

> **NOTE:** This tool is going through rapid development and could change often. Check back for updates.

GWA CLI is a tool for composing, validating and generating Kong Gateway configuration files from OpenAPI (aka Swagger) specs and managing Kong Plugins.

## Installation

```shell
$ npm install -g @bcgov/gwa-cli
```

#### Prerequisites

- Node 12+
- npm (latest, comes installed with most node installers)

Below are a few setting values that you should have ready before publishing workspaces to help the process go quicker.

##### A Namespace

A namespace represents a collections of Kong Services and Routes that are managed independently.

To create a new namespace, go to the [API Services Portal](https://gwa-qwzrwc-test.pathfinder.gov.bc.ca/int).

##### Client ID & Secret

While in the [API Services Portal](https://gwa-qwzrwc-test.pathfinder.gov.bc.ca/int), go to **Service Accounts** tab and click the **Create Service Account**. A new credential will be created - make a note of the `CLIENT ID` and `CLIENT SECRET`, they will only be shown once.

## Usage

For the following usage example we'll use a demonstration namespace of `sampler`, and an OpenAPI spec available at https://website/swagger.json

Initialize a new configuration

```shell
$ cd /to/project/dir
$ mkdir configs && cd configs
$ gwa init --namespace=sampler --client-id=<CLIENT ID> --client-secret=<CLIENT SECRET>
```

Import OpenAPI spec and convert to a Kong configuration file with the rate-limiting plugin.

```shell
$ gwa new https://website/swagger.json --service=sampler-service \
  --plugins=rate-limiting oidc --outfile=sampler-service.yaml
```

Note you can see the list of available plugins and their description by running `$ gwa plugins`. Copy the **Plugin ID** to use in the `new` command. Use a space to add multiple plugins.

**TIP** Any command that accepts a URL as its input (i.e. `$ gwa <command> <input> --options`) can also accept a file on your local disk.

The result of this command will be a `sampler-service.yaml` file in the current directory. The plugins property in the service config will have the fields needed to configure the service. You can fill these in using you IDE of choice.

After filling out the plugins settings, check that your entries are valid by running

```shell
$ gwa validate sampler-service.yaml
```

If the config file is valid you're ready to publish. Run

```shell
$ gwa publish-gateway sampler-service.yaml --dry-run
```

If successful the shell will print a success message. Next you can add admins and users via their IDIRs.

```shell
$ gwa acl --managers=acope@idir --users=jjones@idir
```

A success message will return added, removed and missing users. Also note that `acl` will replace the remote admin/user's list, not append.

Lastly, if your API routes change after publishing your API gateway config, you can update and republish by running

```shell
$ gwa update sampler-service.yaml -u https://website/swagger.json
```

You can re-validate and publish your gateway with the updated routes.

## Commands

### `gwa init`

Generates a `.env` file in the current working directory.

```shell
--namespace        The namespace of you routes collection
--client-id        Namespace's Client ID from API Services Portal
--client-secret    Namespace's Client Secret from API Services Portal
```

A `.env` file should have the following key/values

```
GWA_NAMESPACE=<your namespace>
CLIENT_ID=<gwa client id>
CLIENT_SECRET=<gwa client secret>
GWA_ENV=<dev, prod or test>
```

To create and work with configurations you don't require `CLIENT_ID` or `CLIENT_SECRET`, but to make any API requests you will.

Note you can copy this output above and paste the env keys from the sources mentioned in Prerequistes.

### `gwa new <input file or URL>`

Initialize a config file in the current directory. The input file must be an OpenAPI JSON file or URL

```shell
--service    Name of the service
--plugins    Space separated list of plugin IDs (dash separated, see plugins command)
--outfile    The file to write to write output to
```

### `gwa validate <input file>`

Validate a configuration file

### `gwa update <input file>`

```shell
--url        URL of OpenAPI/Swagger JSON to update
--file       Local file of OpenAPI/Swagger JSON to update. Not required if --url is set
```

Update a config with new OpenAPI specs

### `gwa plugins`

List all available plugins

### `gwa publish-gateway <config file>` Alias `pg`

Publish all YAML config files in current directory. Make sure your `.env` file is configured correctly.

```shell
--dry-run    true/false    Publish as a dry run only
```

### `gwa acl`

Update the full membership. Note that this command will overwrite the remote list of users, use with caution

```shell
--managers    A list of IDs to be giving admin roles
--users       A list of IDs to be giving read-only roles
```

## Help

Run `$ gwa --help` to see all available commands, `$ gwa <command> --help` to view an individual command's help content.

## Development

To install checkout this repo in the `dev` branch, then run the following:

```bash
$ npm i
$ npm start
$ npm link
```

To uninstall simply run `$ npm uninstall` from this directory.
