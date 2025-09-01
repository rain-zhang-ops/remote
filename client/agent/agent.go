package agent

import (
	"context"
	"errors"
	"time"

	controlpb "example.com/remote/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Agent encapsulates core modules of the client daemon.
type Agent struct{}

// New creates a new Agent instance.
func New() *Agent {
	return &Agent{}
}

// Run connects to the control plane, enrolls the device, and blocks until the
// context is canceled. It returns the assigned device identifier.
func (a *Agent) Run(ctx context.Context, addr, token string, opts ...grpc.DialOption) (string, error) {
	if len(opts) == 0 {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
	dialOpts := append(opts, grpc.WithBlock())
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, addr, dialOpts...)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	id, err := a.Enroll(ctx, controlpb.NewDeviceRegistryClient(conn), token)
	if err != nil {
		return "", err
	}
	<-ctx.Done()
	return id, nil
}

// Enroll registers the device with the control plane using the provided client.
func (a *Agent) Enroll(ctx context.Context, client controlpb.DeviceRegistryClient, token string) (string, error) {
	if token == "" {
		return "", errors.New("token required")
	}
	resp, err := client.RegisterDevice(ctx, &controlpb.EnrollReq{Token: token})
	if err != nil {
		return "", err
	}
	if resp.GetDeviceId() == "" {
		return "", errors.New("empty device id")
	}
	return resp.GetDeviceId(), nil
}
