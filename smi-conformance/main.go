package main

import (
	"os"
	"time"

	"github.com/layer5io/learn-layer5/smi-conformance/grpc"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	service := &grpc.Service{
		Name:      "smi-conformance",
		Port:      "10011",
		Version:   "v1.0.0",
		StartedAt: time.Now(),
	}
	// Initialize Logger instance
	logger := log.Log.WithName(service.Name)
	logger.Info("Conformance tool Started")
	// Server Initialization
	err := grpc.Start(service)
	if err != nil {
		logger.Error(err, "Conformance tool crashed!!")
		os.Exit(1)
	}
}
