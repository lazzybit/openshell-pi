package interceptor

import (
	"context"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestNamedTemplateExpansion(t *testing.T) {
	op, _ := structpb.NewStruct(map[string]any{"workloadTemplate": "pi", "spec": map[string]any{"command": []any{"pi"}}})
	req := &pb.InterceptorEvaluation{Service: "openshell.v1.OpenShell", Method: "CreateSandbox", Phase: &pb.InterceptorEvaluation_ModifyOperation{ModifyOperation: &pb.ModifyOperationEvaluation{ProposedOperation: op}}}
	server := &Server{Templates: map[string]map[string]any{"pi": {"image": "example/pi@sha256:abc"}}}
	result, err := server.Evaluate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Allowed || len(result.Patches) != 3 {
		t.Fatalf("expected expansion, got %v", result)
	}
}

func TestInlineDNSPatchPreservesOtherAnnotations(t *testing.T) {
	op, _ := structpb.NewStruct(map[string]any{"spec": map[string]any{"template": map[string]any{"image": "example/image", "annotations": map[string]any{"keep": "value"}}}})
	req := &pb.InterceptorEvaluation{Service: "openshell.v1.OpenShell", Method: "CreateSandbox", Phase: &pb.InterceptorEvaluation_ModifyOperation{ModifyOperation: &pb.ModifyOperationEvaluation{ProposedOperation: op}}}
	result, err := new(Server).Evaluate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Allowed || len(result.Patches) != 1 {
		t.Fatalf("unexpected result: %v", result)
	}
	annotations := result.Patches[0].Value.GetStructValue().AsMap()
	if annotations["keep"] != "value" || annotations["ambient.istio.io/dns-capture"] != "false" {
		t.Fatalf("annotations: %v", annotations)
	}
}
