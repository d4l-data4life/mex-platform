# FAQs

## How to configure MEx?

The recommended configuration setup is described [here](./config.md).

## How can MEx be extended?

The MEx backend services are implemented in Golang and expose an HTTP RESTful API.
The actual API endpoints are defined via [gRPC](https://grpc.io/).
We then use a [gRPC-to-REST gateway library](https://github.com/grpc-ecosystem/grpc-gateway) to derive the code for the HTTP endpoints.
Hence, if you want to amend or add API endpoints, you need to amend the respective `proto` files and re-generate the stubs.
The `proto` files for each service reside in their respective `endpoints` packages.
We checked in all generated stub files, so you do not have to run any gRPC-related tasks to run MEx with the standard API.

## Are there integration und UI test available?

We have used an extensive number of integration and UI tests to ensure a working codebase while developing.
However, almost all of these tests need to dynamically prepare specific MEx configurations in order to run.
We partly use Kirby CMS content to maintain the configuration (see the above question on configuration), and you will need to have a Kirby setup or equivalent to run such tests.
Since the test setup depends on the details of how configurations are managed, we have not included the test suite here.
