package interceptor

import (
	"context"
	ext "github.com/lazzybit/openshell-pi-interceptor/internal/pb/extensionv1"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/structpb"
	"net"
	"testing"
)

func TestGRPCContract(t *testing.T) {
	listener := bufconn.Listen(1048576)
	server := grpc.NewServer()
	pb.RegisterGatewayInterceptorServer(server, &Server{})
	go server.Serve(listener)
	t.Cleanup(server.Stop)
	conn, err := grpc.NewClient("passthrough:///buffer", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGatewayInterceptorClient(conn)
	manifest, err := client.Describe(context.Background(), &pb.DescribeRequest{Gateway: &ext.PeerMetadata{ProtocolVersion: &ext.ProtocolVersion{Major: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	if manifest.FailurePolicy != "fail_closed" || len(manifest.Bindings) != 1 {
		t.Fatalf("invalid manifest: %v", manifest)
	}
	op, _ := structpb.NewStruct(map[string]any{"spec": map[string]any{"template": map[string]any{"image": "pi"}}})
	result, err := client.Evaluate(context.Background(), &pb.InterceptorEvaluation{Service: "openshell.v1.OpenShell", Method: "CreateSandbox", Phase: &pb.InterceptorEvaluation_ModifyOperation{ModifyOperation: &pb.ModifyOperationEvaluation{ProposedOperation: op}}})
	if err != nil || !result.GetAllowed() {
		t.Fatalf("result=%v error=%v", result, err)
	}
}

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
