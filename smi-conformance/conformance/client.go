package conformance

import (
	context "context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// ConformanceClient represents a gRPC adapter client
type ConformanceClient struct {
	CClient ConformanceTestingClient
	conn    *grpc.ClientConn
}

// CreateClient creates a ConformanceClient for the given params
func CreateClient(ctx context.Context, conformanceLocationURL string) (*ConformanceClient, error) {
	var opts []grpc.DialOption
	// creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to create TLS credentials %v", err)
	// 	}
	// 	opts = append(opts, grpc.WithTransportCredentials(creds))
	// } else {
	opts = append(opts, grpc.WithInsecure())
	// }
	conn, err := grpc.Dial(conformanceLocationURL, opts...)
	if err != nil {
		logrus.Errorf("fail to dial: %v", err)
	}

	cClient := NewConformanceTestingClient(conn)

	return &ConformanceClient{
		conn:    conn,
		CClient: cClient,
	}, nil
}

// Close closes the ConformanceClient
func (c *ConformanceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
