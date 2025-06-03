// Package provider implements the Flexera Cost provider for Pulumi.
package provider

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type flexeraProvider struct {
	pulumirpc.UnimplementedResourceProviderServer

	host    *plugin.Host
	name    string
	version string
}

// NewProvider returns a new Flexera provider
func NewProvider(host *plugin.Host, name, version string) (pulumirpc.ResourceProviderServer, error) {
	return &flexeraProvider{
		host:    host,
		name:    name,
		version: version,
	}, nil
}

// GetSchema returns the schema for the provider
func (p *flexeraProvider) GetSchema(ctx context.Context,
	req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	// Return stub schema for now
	return &pulumirpc.GetSchemaResponse{
		Schema: "{}",
	}, nil
}

// Configure configures the provider
func (p *flexeraProvider) Configure(ctx context.Context,
	req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	// Stub configuration
	return &pulumirpc.ConfigureResponse{}, nil
}

// CheckConfig validates the configuration
func (p *flexeraProvider) CheckConfig(ctx context.Context,
	req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{
		Inputs: req.News,
	}, nil
}

// Check validates resource inputs
func (p *flexeraProvider) Check(ctx context.Context,
	req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{
		Inputs: req.News,
	}, nil
}

// Diff computes the diff between the old and new resource states
func (p *flexeraProvider) Diff(ctx context.Context,
	req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

// Create creates a new resource
func (p *flexeraProvider) Create(ctx context.Context,
	req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	return &pulumirpc.CreateResponse{
		Id:         req.Urn,
		Properties: req.Properties,
	}, nil
}

// Read reads the current state of a resource
func (p *flexeraProvider) Read(ctx context.Context,
	req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return &pulumirpc.ReadResponse{
		Id:         req.Id,
		Properties: req.Properties,
	}, nil
}

// Update updates a resource
func (p *flexeraProvider) Update(ctx context.Context,
	req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return &pulumirpc.UpdateResponse{
		Properties: req.News,
	}, nil
}

// Delete deletes a resource
func (p *flexeraProvider) Delete(ctx context.Context,
	req *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// GetPluginInfo returns plugin metadata
func (p *flexeraProvider) GetPluginInfo(ctx context.Context,
	req *empty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
}

// Cancel cancels the provider
func (p *flexeraProvider) Cancel(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Invoke invokes a function
func (p *flexeraProvider) Invoke(ctx context.Context,
	req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Invoke not implemented")
}

// StreamInvoke invokes a streaming function
func (p *flexeraProvider) StreamInvoke(req *pulumirpc.InvokeRequest,
	server interface{}) error {
	return status.Error(codes.Unimplemented, "StreamInvoke not implemented")
}

// Attach attaches to an existing provider
func (p *flexeraProvider) Attach(ctx context.Context,
	req *pulumirpc.PluginAttach) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
