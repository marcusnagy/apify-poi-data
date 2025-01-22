package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"apify-poi-data/internal/services"
	"apify-poi-data/pkg/apify"
	maps_v1 "apify-poi-data/proto/apify/maps/v1"
	poi_v1 "apify-poi-data/proto/apify/poi/v1"
	tripsadvisor_v1 "apify-poi-data/proto/apify/tripsadvisor/v1"
)

func setupGRPCServer(port int, tlsConfig *tls.Config) (*grpc.Server, net.Listener, error) {
	grpc.EnableTracing = true
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))

	server := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	// Register services
	maps_v1.RegisterMapsServiceServer(server, &services.MapsService{
		ApifyClient: apify.NewClient(
			cfg.Apify.Key,
			cfg.Apify.ActorExtractorID,
			cfg.Apify.ActorScraperID,
		),
		Database: db,
	})
	tripsadvisor_v1.RegisterTripadvisorServiceServer(server, &services.TripadvisorService{})
	poi_v1.RegisterPoiServiceServer(
		server,
		&services.PoiService{
			Database: db,
		},
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}

	reflection.Register(server)

	return server, listener, nil
}

func ShutdownGRPCServer(server *grpc.Server) {
	server.GracefulStop()
}

func SetupHTTPMux(ctx context.Context, grpcPort, httpPort int, tlsConfig *tls.Config) (*http.Server, error) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))}

	err := maps_v1.RegisterMapsServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		return nil, err
	}

	err = tripsadvisor_v1.RegisterTripadvisorServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		return nil, err
	}

	err = poi_v1.RegisterPoiServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		return nil, err
	}

	// Register gRPC gateway handlers
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: mux,
		// TLSConfig: tlsConfig,
	}

	return server, nil
}

func loadTLSCredentials() (*tls.Config, error) {
	certFile := os.Getenv("TLS_CERT_FILE")
	if certFile == "" {
		certFile = "certs/server.crt" // default value
	}

	keyFile := os.Getenv("TLS_KEY_FILE")
	if keyFile == "" {
		keyFile = "certs/server.key" // default value
	}

	caFile := os.Getenv("TLS_CA_FILE")
	if caFile == "" {
		caFile = "certs/rootCA.pem" // default value
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Load your root CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	return config, nil
}
