package interceptor

import (
	"context"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestRejectMixedTemplates(t *testing.T) {
	op, _ := structpb.NewStruct(map[string]any{"workloadTemplate": "pi", "spec": map[string]any{"template": map[string]any{"image": "unapproved"}}})
	req := &pb.InterceptorEvaluation{Service: "openshell.v1.OpenShell", Method: "CreateSandbox", Phase: &pb.InterceptorEvaluation_ModifyOperation{ModifyOperation: &pb.ModifyOperationEvaluation{ProposedOperation: op}}}
	_, err := (&Server{Templates: map[string]map[string]any{"pi": {"image": "approved"}}}).Evaluate(context.Background(), req)
	if err == nil {
		t.Fatal("mixed named and inline templates accepted")
	}
}
