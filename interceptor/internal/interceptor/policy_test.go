package interceptor

import (
	"context"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/protobuf/types/known/structpb"
	"reflect"
	"testing"
)

func TestDeploymentPolicyReplacesCallerPolicy(t *testing.T) {
	approved := map[string]any{"version": float64(1), "filesystem": map[string]any{"readWrite": []any{"/tmp", "/dev/null", "/dev/pts"}}}
	for _, caller := range []map[string]any{nil, {"filesystem": map[string]any{"readWrite": []any{"/"}}}} {
		spec := map[string]any{}
		if caller != nil {
			spec["policy"] = caller
		}
		op, _ := structpb.NewStruct(map[string]any{"workloadTemplate": "pi", "spec": spec})
		s := &Server{Templates: map[string]map[string]any{"pi": {"image": "approved"}}, Policy: approved}
		result, err := s.Evaluate(context.Background(), &pb.InterceptorEvaluation{Service: "openshell.v1.OpenShell", Method: "CreateSandbox", Phase: &pb.InterceptorEvaluation_ModifyOperation{ModifyOperation: &pb.ModifyOperationEvaluation{ProposedOperation: op}}})
		if err != nil || !result.GetAllowed() {
			t.Fatalf("result=%v err=%v", result, err)
		}
		found := false
		for _, p := range result.Patches {
			if p.Path == "/spec/policy" {
				found = true
				if !reflect.DeepEqual(p.Value.GetStructValue().AsMap(), approved) {
					t.Fatalf("wrong policy %v", p.Value)
				}
			}
		}
		if !found {
			t.Fatal("deployment policy was not injected")
		}
	}
}
