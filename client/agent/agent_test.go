package agent

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	controlpb "example.com/remote/proto"
	svc "example.com/remote/server/control"
)

const bufSize = 1024 * 1024

func dialer(s *grpc.Server) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(bufSize)
	go s.Serve(lis)
	return func(ctx context.Context, s string) (net.Conn, error) { return lis.Dial() }
}

// TestAgentEnroll ensures the agent can register via gRPC.
func TestAgentEnroll(t *testing.T) {
	s := grpc.NewServer()
	controlpb.RegisterDeviceRegistryServer(s, &svc.Server{Store: &svc.MemoryStore{}})
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(s)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	a := New()
	id, err := a.Enroll(ctx, controlpb.NewDeviceRegistryClient(conn), "tok")
	if err != nil {
		t.Fatalf("enroll failed: %v", err)
	}
	if id == "" {
		t.Fatalf("expected id, got empty")
	}
}

// TestAgentEnroll_EmptyToken ensures validation failure on missing token.
func TestAgentEnroll_EmptyToken(t *testing.T) {
	s := grpc.NewServer()
	controlpb.RegisterDeviceRegistryServer(s, &svc.Server{Store: &svc.MemoryStore{}})
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(s)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	a := New()
	if _, err := a.Enroll(ctx, controlpb.NewDeviceRegistryClient(conn), ""); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestAgentRun verifies that Run enrolls the device and waits for context cancellation.
func TestAgentRun(t *testing.T) {
	s := grpc.NewServer()
	controlpb.RegisterDeviceRegistryServer(s, &svc.Server{Store: &svc.MemoryStore{}})
	ctx, cancel := context.WithCancel(context.Background())

	a := New()
	done := make(chan struct{})
	go func() {
		defer close(done)
		id, err := a.Run(ctx, "bufnet", "tok", grpc.WithContextDialer(dialer(s)), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Errorf("Run returned error: %v", err)
			return
		}
		if id == "" {
			t.Errorf("expected id, got empty")
		}
	}()
	time.AfterFunc(50*time.Millisecond, cancel)
	<-done
}

// TestAgentRun_EmptyToken verifies Run propagates validation errors.
func TestAgentRun_EmptyToken(t *testing.T) {
	s := grpc.NewServer()
	controlpb.RegisterDeviceRegistryServer(s, &svc.Server{Store: &svc.MemoryStore{}})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := New()
	if _, err := a.Run(ctx, "bufnet", "", grpc.WithContextDialer(dialer(s)), grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
