package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

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

func ebpf_loader() ([]*ebpf.Map, error) {
	// Remove resource limits for kernels <5.11.
	// if err := rlimit.RemoveMemlock(); err != nil {
	// 	log.Fatal("Removing memlock:", err)
	// }

	spec, err := ebpf.LoadCollectionSpec("/users/anakin/ebpf_service_mesh/user_proxy/kernel/xdp_main")
	if err != nil {
		panic(err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		panic(fmt.Sprintf("Failed to create new collection: %v\n", err))
	}
	defer coll.Close()

	prog := coll.Programs["tc_egress"]
	service_catalog := coll.Maps["service_catalog"]
	if prog == nil {
		panic("No program named 'tc_egress' found in collection")
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
		Program:   prog,
		Attach:    ebpf.AttachTCXEgress,
	})

	defer lnk.Close()

	fmt.Println("Successfully loaded and attached BPF program.")

	// return the reference to maps and other required information
	return []*ebpf.Map{service_catalog}, nil
}

func initialize_maps(maps []*ebpf.Map) error {

	// service catalog
	service_catalog := maps[0]

	// get service catalog and store in ebpf map
	serviceMap, err := getServiceCatalog()

	if err != nil {
		fmt.Errorf("Error in fetching service catalog")
	}

	// DEBUG := Print the service map
	for serviceName, info := range serviceMap {
		fmt.Printf("Service: %s, Address: %s, ServiceAddress: %s, ServicePort: %d\n", serviceName, info.Address, info.ServiceAddress, info.ServicePort)

		// step-1 : convert IP to desired form
		ip := net.ParseIP(info.Address).To4()
		if ip == nil {
			fmt.Println("Invalid IP address")
			continue
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
		serviceName = serviceName[:int(len(serviceName[:]))]

		// Print the service name and IP
		fmt.Printf("Service: %s, IP Address: %d\n", serviceName, value.ServiceIP)
	}

	// Check if there was an error during iteration
	if err := iterator.Err(); err != nil {
		log.Printf("Failed to iterate over map: %v", err)
	}

	return nil
}
