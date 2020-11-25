package main

import (
	"fmt"
	"os"
	"time"

	"github.com/layer5io/learn-layer5/smi-conformance/grpc"
)

func main() {

	service := &grpc.Service{
		Name:      "smi-conformance",
		Port:      "10011",
		Version:   "v1.0.0",
		StartedAt: time.Now(),
	}

	// Server Initialization
	fmt.Println("Conformance tool Started")
	err := grpc.Start(service)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
