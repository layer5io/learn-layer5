package grpc

import (
	"errors"
	"fmt"
	"net"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"

	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
)

// Service object holds all the information about the server parameters.
type Service struct {
	Name      string    `json:"name"`
	Port      string    `json:"port"`
	Version   string    `json:"version"`
	StartedAt time.Time `json:"startedat,string"`
}

// panicHandler is the handler function to handle panic errors
func panicHandler(r interface{}) error {
	fmt.Println("500 panic Error")
	return errors.New(fmt.Sprintf("Panic error for: %+v", r))
}

// StartServer starts grpc server
func Start(s *Service) error {

	address := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	middlewares := middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandler(panicHandler),
		),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares),
	)

	//Register Proto
	conformance.RegisterConformanceTestingServer(server, s)

	// Start serving requests
	if err = server.Serve(listener); err != nil {
		return err
	}
	return nil

}
