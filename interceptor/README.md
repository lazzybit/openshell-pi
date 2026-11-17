# OpenShell Pi interceptor

Go gRPC service implementing the official OpenShell v0.1.8 Gateway Interceptor protocol. The only binding is `openshell.v1.OpenShell/CreateSandbox` in `modify_operation`, with fail-closed behavior.

It preserves existing Pod annotations and forces `ambient.istio.io/dns-capture=false`. For `--template pi`, OpenShell v0.1.8 cannot combine a named workload template with inline annotations, so an operator must provide a JSON snapshot mapping approved template names to inline templates:

```json
{"pi":{"image":"registry.example/pi@sha256:<approved-digest>"}}
```

Start with `openshell-pi-interceptor --templates /config/templates.json`. The image is deployment configuration, never compiled into Go. The service removes `workloadTemplate`, supplies the configured inline template and adds the annotation. Unknown names fail closed. Native named-template provenance is lost. Keep this snapshot aligned with the registered template during image updates; this is a v0.1.8 compatibility workaround.

Pass `--policy /config/policy.json` to load the deployment-owned SandboxPolicy in protobuf JSON format. The service validates its protobuf schema at startup and replaces `/spec/policy` on CreateSandbox, including caller-supplied policies. Clients do not need a local policy file. Keep the policy in deployment configuration, outside the runtime image. The deployed Pi policy includes `/dev/pts` for tmux.

Without a named template, an inline `spec.template` is required. This service does not enforce command, tenant quotas or an image allowlist for inline requests. Do not describe it as complete workload admission governance.

Build and test: `go test ./...`, `CGO_ENABLED=0 go build ./cmd/openshell-pi-interceptor`, or `docker build .`. Generated bindings and original Apache-2.0 proto files are vendored from NVIDIA/OpenShell v0.1.8. The binary runs as uid 10001 and listens on port 50051. It exposes the standard gRPC health service.

Gateway registration must use `binding_policy=exact`, CreateSandbox `modify_operation`, and `failure_policy=fail_closed`. This server uses plaintext gRPC and expects deployment-level caller authentication such as Istio ambient mTLS with an exact Gateway ServiceAccount authorization rule. Never expose it publicly. Native extension JWT verification is not implemented in this package.
