package main

import (
	"os"
	"sync"
	"time"
)

func main() {

	// get cli arguments
	node_name := os.Args[1]

	// wait group for go routines
	var wg sync.WaitGroup

	// Load ebpf program
	// ---- CUSTOMIZE ----
	// implement the ebpf_loader function as per your application requirements
	// Return Values : array of ebpf map required, Error
	bpf_maps, err := ebpf_loader()
	if err != nil {
		panic("Error in loading ebpf program")
	}

	//  OPTIONAL: Channel to communicate between goroutines
	dataChannel := make(chan Data)
	wg.Add(1)

	// ---- CUSTOMIZE ----
	// Add new endpoints to the http server as per user requirements
	go startServer(dataChannel)

	// sleep until infra sets up
	time.Sleep(1 * time.Minute)

	// ---- CUSTOMIZE ----
	// implement your logic to insitialize your ebpf maps present in array returned
	// by the ebpf_loader() function.
	initialize_maps(bpf_maps)

	// ---- ADDITIONS ----
	// Add custom go routine to add extra required functionalities
	// Below is the example that computes and updates the KV Store

	// go routine for regular metric collection and updation in KV Store
	wg.Add(1)
	go compute_metrics(node_name)

	// wait for all go routines to terminate
	wg.Wait()
}
