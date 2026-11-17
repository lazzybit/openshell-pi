package interceptor

import (
	"context"
	ext "github.com/lazzybit/openshell-pi-interceptor/internal/pb/extensionv1"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Server implements OpenShell's version 1 Gateway Interceptor contract.
type Server struct {
	pb.UnimplementedGatewayInterceptorServer
	Templates map[string]map[string]any
	Policy    map[string]any
}

func (*Server) Describe(_ context.Context, req *pb.DescribeRequest) (*pb.InterceptorManifest, error) {
	if req.GetGateway().GetProtocolVersion().GetMajor() != 1 {
		return nil, status.Error(codes.FailedPrecondition, "extension protocol major must be 1")
	}
	return &pb.InterceptorManifest{Name: "openshell-pi-interceptor", FailurePolicy: "fail_closed", Extension: &ext.PeerMetadata{ProtocolVersion: &ext.ProtocolVersion{Major: 1}, ImplementationName: "openshell-pi-interceptor", ImplementationVersion: "1", SupportedCapabilities: []string{"openshell.gateway-interceptor.contract"}}, Bindings: []*pb.InterceptorBinding{{Id: "pi-dns", Selector: &pb.InterceptorSelector{Rpc: "openshell.v1.OpenShell/CreateSandbox"}, Phases: []pb.GatewayInterceptorPhase{pb.GatewayInterceptorPhase_GATEWAY_INTERCEPTOR_PHASE_MODIFY_OPERATION}}}}, nil
}

func (s *Server) Evaluate(_ context.Context, req *pb.InterceptorEvaluation) (*pb.InterceptorResult, error) {
	if req.Service != "openshell.v1.OpenShell" || req.Method != "CreateSandbox" || req.GetModifyOperation() == nil {
		return nil, status.Error(codes.InvalidArgument, "unexpected RPC or phase")
	}
	op := req.GetModifyOperation().GetProposedOperation().AsMap()
	spec, ok := op["spec"].(map[string]any)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "spec is required")
	}
	environment, _ := spec["environment"].(map[string]any)
	if environment == nil {
		environment = map[string]any{}
	}
	environment["LANG"] = "C.UTF-8"
	environmentValue, err := structpb.NewValue(environment)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid environment")
	}
	patches := []*pb.JsonPatch{{Op: "add", Path: "/spec/environment", Value: environmentValue}}
	if s.Policy != nil {
		value, err := structpb.NewValue(s.Policy)
		if err != nil {
			return nil, status.Error(codes.Internal, "invalid deployment policy")
		}
		patches = append(patches, &pb.JsonPatch{Op: "add", Path: "/spec/policy", Value: value})
	}
	if name, ok := op["workloadTemplate"].(string); ok && name != "" {
		configured, exists := s.Templates[name]
		if !exists {
			return &pb.InterceptorResult{Allowed: false, StatusCode: "PERMISSION_DENIED", Reason: "template is not configured"}, nil
		}
		if _, exists := spec["template"]; exists {
			return nil, status.Error(codes.InvalidArgument, "named and inline templates cannot be combined")
		}
		value, err := structpb.NewValue(configured)
		if err != nil {
			return nil, status.Error(codes.Internal, "invalid template configuration")
		}
		patches = append(patches, &pb.JsonPatch{Op: "remove", Path: "/workloadTemplate"}, &pb.JsonPatch{Op: "add", Path: "/spec/template", Value: value})
		spec["template"] = value.GetStructValue().AsMap()
	}
	template, ok := spec["template"].(map[string]any)
	if !ok {
		return &pb.InterceptorResult{Allowed: false, StatusCode: "FAILED_PRECONDITION", Reason: "inline template or configured named template required"}, nil
	}
	annotations, _ := template["annotations"].(map[string]any)
	if annotations == nil {
		annotations = map[string]any{}
	}
	annotations["ambient.istio.io/dns-capture"] = "false"
	value, err := structpb.NewValue(annotations)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid annotations")
	}
	return &pb.InterceptorResult{Allowed: true, Patches: append(patches, &pb.JsonPatch{Op: "add", Path: "/spec/template/annotations", Value: value})}, nil
}
