package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// Data struct to hold incoming JSON data
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Struct to hold the Address, ServiceAddress, and ServicePort fields
type ServiceInfo struct {
	Address        string
	ServiceAddress string
	ServicePort    int
}

type catalogKey struct {
	Hostname [256]byte
}

type catalogValue struct {
	ServiceIP uint32
}

func main() {
	// Channel to communicate between goroutines
	dataChannel := make(chan Data)

	// Remove resource limits for kernels <5.11.
	// if err := rlimit.RemoveMemlock(); err != nil {
	// 	log.Fatal("Removing memlock:", err)
	// }

	spec, err := ebpf.LoadCollectionSpec("/app/kernel/xdp_main")
	if err != nil {
		panic(err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		panic(fmt.Sprintf("Failed to create new collection: %v\n", err))
	}
	defer coll.Close()

	prog_egress := coll.Programs["tc_egress_helper"]
	prog_ingress := coll.Programs["tc_ingress_helper"]
	service_catalog := coll.Maps["service_catalog"]
	if prog_egress == nil || prog_ingress == nil {
		panic("No program named 'tc_egress_helper' or 'tc_ingress_helper' found in collection")
	}
	if service_catalog == nil {
		panic("No map named 'service_catalog' found in collection")
	}

	iface := "lo"

	iface_idx, err := net.InterfaceByName(iface)
	if err != nil {
		panic(fmt.Sprintf("Failed to get interface %s: %v\n", iface, err))
	}

	lnk, _ := link.AttachTCX(link.TCXOptions{
		Interface: iface_idx.Index,
		Program:   prog_egress,
		Attach:    ebpf.AttachTCXEgress,
	})

	defer lnk.Close()

	lnk2, _ := link.AttachTCX(link.TCXOptions{
		Interface: iface_idx.Index,
		Program:   prog_ingress,
		Attach:    ebpf.AttachTCXIngress,
	})

	defer lnk2.Close()

	fmt.Println("Successfully loaded and attached BPF program.")

	// Start the server in a goroutine
	go startServer(dataChannel)

	// sleep until infra sets up
	time.Sleep(1 * time.Minute)

	// get service catalog and store in ebpf map
	serviceMap, err := getServiceCatalog()

	// DEBUG := Print the service map
	for serviceName, info := range serviceMap {
		fmt.Printf("Service: %s, Address: %s, ServiceAddress: %s, ServicePort: %d\n", serviceName, info.Address, info.ServiceAddress, info.ServicePort)

		// step-1 : convert IP to desired form
		ip := net.ParseIP(info.Address).To4()
		if ip == nil {
			fmt.Println("Invalid IP address")
			return
		}
		serviceIP := binary.BigEndian.Uint32(ip) // Convert IP to uint32 in network byte order

		// step-2 : convert serviceName to required form
		var key catalogKey
		copy(key.Hostname[:], serviceName)

		// Step 3: Prepare value (service IP)
		value := catalogValue{
			ServiceIP: serviceIP,
		}

		if err := service_catalog.Put(&key, &value); err != nil {
			fmt.Errorf("Error in storing service %s to the map :: %s", serviceName, err)
		} else {
			fmt.Printf("serviceip :: %d --  key :: %s \n", serviceIP, key.Hostname)
		}
	}

	// Create key and value holders
	var key catalogKey
	var value catalogValue

	iterator := service_catalog.Iterate()
	for iterator.Next(&key, &value) {
		// Convert the service name (key) to a Go string
		serviceName := string(key.Hostname[:])
		// Null terminate the service name if there are extra zeros
		serviceName = serviceName[:len(serviceName)-len(serviceName)+int(len(serviceName[:]))]

		// Print the service name and IP
		fmt.Printf("Service: %s, IP Address: %d\n", serviceName, value.ServiceIP)
	}

	// Check if there was an error during iteration
	if err := iterator.Err(); err != nil {
		log.Fatalf("Failed to iterate over map: %v", err)
	}

	// Main goroutine handles data from the channel
	for {

	}
}
