# `mexctl`: MEx CLI

This tool allows setting up a complete configuration and initial data for an entire MEx system.

## Parameters

For help on parameters:

```sh
$ npm start
```

## Building and running the tool

To use the tool first install and build it.

First build and link the required `mexlib` module:

```sh
# in ./tools/mexlib
npm install
npm build
npm link
```

Then build the tool:

```sh
# in ./tools/mexctl
npm install
npm link mexlib
npm run build
```

The tool needs the connection details to talk to MEx instances.
They are stored in a YAML file (similar to `kubectl`) which by default is located at `~/.mexctl/config.yaml`.
It may contain multiple instance configurations.
Please use [config.yaml.example](../../config.yaml.example) as a template.

You can also use the `--config` argument to override the config.
If specified as `--config @foo`, then `foo` is interpreted as a file name.
If specified as `--config foo`, then `foo` is interpreted as the YAML content of the config.
The effective instance is either the value of `defaultInstance` in the configuration, or can be overridden via `--instance`.

## Examples

See [Makefile](../../Makefile).
