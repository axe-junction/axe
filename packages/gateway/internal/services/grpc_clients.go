package services

import (
	"context"
	"fmt"
	"log"

	"github.com/axe-junction/axe-server/gateway/internal/config"
	pb "github.com/axe-junction/axe-server/gateway/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoutingServiceClient struct {
	client pb.RoutingServiceClient
	conn   *grpc.ClientConn
}

func NewRoutingServiceClient(cfg *config.Config) (*RoutingServiceClient, error) {
	address := fmt.Sprintf("%s:%s", cfg.Services.RoutingService.Host, cfg.Services.RoutingService.Port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to routing service: %w", err)
	}

	client := pb.NewRoutingServiceClient(conn)

	return &RoutingServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (r *RoutingServiceClient) GetBestRoute(ctx context.Context, req *pb.RouteRequest) (*pb.RouteResponse, error) {
	return r.client.GetBestRoute(ctx, req)
}

func (r *RoutingServiceClient) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

type VTCServiceClient struct {
	client pb.VTCServiceClient
	conn   *grpc.ClientConn
}

func NewVTCServiceClient(cfg *config.Config) (*VTCServiceClient, error) {
	address := fmt.Sprintf("%s:%s", cfg.Services.VTCService.Host, cfg.Services.VTCService.Port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to VTC service: %w", err)
	}

	client := pb.NewVTCServiceClient(conn)

	return &VTCServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (v *VTCServiceClient) GetVTCPrices(ctx context.Context, req *pb.PriceRequest) (*pb.PriceResponse, error) {
	return v.client.GetVTCPrices(ctx, req)
}

func (v *VTCServiceClient) Close() error {
	if v.conn != nil {
		return v.conn.Close()
	}
	return nil
}

type GRPCClients struct {
	RoutingService *RoutingServiceClient
	VTCService     *VTCServiceClient
}

func NewGRPCClients(cfg *config.Config) (*GRPCClients, error) {
	routingClient, err := NewRoutingServiceClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to routing service: %v", err)
	}

	vtcClient, err := NewVTCServiceClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to VTC service: %v", err)
	}

	return &GRPCClients{
		RoutingService: routingClient,
		VTCService:     vtcClient,
	}, nil
}

func (g *GRPCClients) Close() {
	if g.RoutingService != nil {
		g.RoutingService.Close()
	}
	if g.VTCService != nil {
		g.VTCService.Close()
	}
}
