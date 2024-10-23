package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	// get cli arguments
	node_name := os.Args[1]
	fmt.Println("Node name : %s", node_name)

	// Load ebpf program
	service_catalog, err := ebpf_loader()
	if err != nil {
		panic("Error in loading ebpf program")
	}

	// Channel to communicate between goroutines
	dataChannel := make(chan Data)
	go startServer(dataChannel)

	// sleep until infra sets up
	time.Sleep(1 * time.Minute)

	// initialize map and other required structure
	initialize_maps(service_catalog)

	// go routine for regular metric collection and updation in KV Store
	go compute_metrics(node_name)

	// Main goroutine handles data from the channel
	for {

	}
}
