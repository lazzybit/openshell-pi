# OpenShell Pi

Two OCI packages are built together for linux/amd64 and linux/arm64:

- `runtime/` → `ghcr.io/lazzybit/openshell-pi-runtime`: Pi workload, pinned npm dependencies, Hindsight and persistent configuration initialization.
- `interceptor/` → `ghcr.io/lazzybit/openshell-pi-interceptor`: Go implementation of the OpenShell v0.1.8 Gateway Interceptor contract.

```sh
docker build -t openshell-pi-runtime:local runtime
docker build -t openshell-pi-interceptor:local interceptor
cd interceptor && go test ./...
```

The interceptor adds the workload Pod DNS capture opt-out during CreateSandbox. For v0.1.8 named templates, it expands an explicitly configured template snapshot into an inline template before adding annotations. Unknown names are denied. This loses native named-template provenance; it is a compatibility workaround, not native named-template annotation support. Deployment configuration owns the approved image and template snapshot, outside the binary. Gateway registration and mesh access policies belong in the infrastructure repository.
