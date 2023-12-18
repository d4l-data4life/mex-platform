# Run MEx demo locally

Everything is running in a container.
All relative paths are given w.r.t. the root folder of this repository.

## Prerequisites

- Docker
- Make
- Node.js (min. v18)

Note that if you are working on a machine with an Apple silicon CPU (ARM-architecture), you may need to change the setting of the environmental variable `GOARCH` from `amd64` to `arm64` in the [backend Dockerfile](../backend/Dockerfile) to make the build work.

## Start containers

In `setup/`:

```sh
make dcup
```

This may take a while during which all images are built and started.
Note that this local setup includes a service that mimics a JWT provider (pairgen service) to allow JWT authentication.
**Never deploy the pairgen service! It is only to be used during development and testing!**

## Prepare system

In `setup/tools/mexctl/`:

```sh
npm install
npm run build
```

In `setup/`:

```sh
make demo-load-mesh
make demo-load-config
make demo-load-items
```

## Use MEx

In browser: https://localhost:53000
Accept the HTTPS error and also wait some seconds until the initial error message disappears (which is due to the absent authentication token).

## Shut down MEx

In `setup/`:

```sh
make dcdownv
```

This stops all containers and removes the volumes.
