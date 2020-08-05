package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kumarabd/gokit/logger"
	"github.com/layer5io/learn-layer5/smi-conformance/grpc"
)

func main() {

	service := &grpc.Service{
		Name:      "smi-conformance",
		Port:      "10008",
		Version:   "v1.0.0",
		StartedAt: time.Now(),
	}

	// Initialize Logger instance
	log, err := logger.New(service.Name)
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	// Server Initialization
	log.Info("Conformance tool Started")
	err = grpc.Start(service)
	if err != nil {
		log.Err("Conformance tool crashed!!", err.Error())
		os.Exit(1)
	}
}
