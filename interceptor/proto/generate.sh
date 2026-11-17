#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
protoc -I proto \
  --go_out=. --go_opt=module=github.com/lazzybit/openshell-pi-interceptor \
  --go_opt=Mdatamodel.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/datamodelv1 \
  --go_opt=Mextension.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/extensionv1 \
  --go_opt=Mgateway_interceptor.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1 \
  --go_opt=Mopenshell.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/openshellv1 \
  --go_opt=Moptions.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/optionsv1 \
  --go_opt=Msandbox.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/sandboxv1 \
  --go-grpc_out=. --go-grpc_opt=module=github.com/lazzybit/openshell-pi-interceptor \
  --go-grpc_opt=Mdatamodel.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/datamodelv1 \
  --go-grpc_opt=Mextension.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/extensionv1 \
  --go-grpc_opt=Mgateway_interceptor.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1 \
  --go-grpc_opt=Mopenshell.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/openshellv1 \
  --go-grpc_opt=Moptions.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/optionsv1 \
  --go-grpc_opt=Msandbox.proto=github.com/lazzybit/openshell-pi-interceptor/internal/pb/sandboxv1 \
  proto/datamodel.proto proto/extension.proto proto/gateway_interceptor.proto \
  proto/openshell.proto proto/options.proto proto/sandbox.proto
