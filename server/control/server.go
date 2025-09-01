package control

import (
	"context"
	"strconv"
	"sync"

	controlpb "example.com/remote/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeviceStore persists device registrations.
type DeviceStore interface {
	Save(ctx context.Context, token string) (string, error)
}

// MemoryStore is an in-memory implementation of DeviceStore.
type MemoryStore struct {
	mu   sync.Mutex
	next int
}

// Save stores a device and returns a new identifier.
func (m *MemoryStore) Save(ctx context.Context, token string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.next++
	return strconv.Itoa(m.next), nil
}

// Server implements the DeviceRegistry gRPC service.
type Server struct {
	controlpb.UnimplementedDeviceRegistryServer
	Store DeviceStore
}

// RegisterDevice stores the enrolling device and returns its identifier.
func (s *Server) RegisterDevice(ctx context.Context, req *controlpb.EnrollReq) (*controlpb.EnrollResp, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}
	id, err := s.Store.Save(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "save failed")
	}
	return &controlpb.EnrollResp{DeviceId: id}, nil
}
