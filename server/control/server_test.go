package control

import (
	"context"
	"testing"

	controlpb "example.com/remote/proto"
)

// TestRegisterDevice verifies successful enrollment.
func TestRegisterDevice(t *testing.T) {
	srv := &Server{Store: &MemoryStore{}}
	resp, err := srv.RegisterDevice(context.Background(), &controlpb.EnrollReq{Token: "abc"})
	if err != nil {
		t.Fatalf("RegisterDevice returned error: %v", err)
	}
	if resp.DeviceId == "" {
		t.Fatalf("expected device id, got empty")
	}
}

// TestRegisterDevice_EmptyToken ensures validation is enforced.
func TestRegisterDevice_EmptyToken(t *testing.T) {
	srv := &Server{Store: &MemoryStore{}}
	_, err := srv.RegisterDevice(context.Background(), &controlpb.EnrollReq{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
