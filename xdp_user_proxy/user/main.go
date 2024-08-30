package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/gin-gonic/gin"
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

	prog := coll.Programs["xdp_prog_simple"]
	if prog == nil {
		panic("No program named 'xdp_prog_simple' found in collection")
	}

	iface := "lo"

	iface_idx, err := net.InterfaceByName(iface)
	if err != nil {
		panic(fmt.Sprintf("Failed to get interface %s: %v\n", iface, err))
	}
	opts := link.XDPOptions{
		Program:   prog,
		Interface: iface_idx.Index,
		// Flags is one of XDPAttachFlags (optional).
	}

	lnk, err := link.AttachXDP(opts)
	if err != nil {
		panic(err)
	}
	defer lnk.Close()

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
	}

	for {
		// Handle map configuration and handling incoming information
	}

	// Main goroutine handles data from the channel
}

func getServiceCatalog() (map[string]ServiceInfo, error) {
	// Make the initial GET request to fetch the service catalog
	res, err := http.Get("http://127.0.0.1:8500/v1/catalog/services")
	if err != nil {
		log.Fatalf("unable to fetch service catalog: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", res.StatusCode)
	}

	// Read and unmarshal the JSON response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("unable to read response body: %v", err)
	}

	var services map[string]interface{}
	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Fatalf("unable to unmarshal JSON response: %v", err)
	}

	// Create a map to hold the service information
	serviceMap := make(map[string]ServiceInfo)

	// Iterate over the keys and make individual requests
	for key := range services {
		// Make a GET request for each service
		res, err := http.Get("http://127.0.0.1:8500/v1/catalog/service/" + key)
		if err != nil {
			log.Printf("error fetching service %s: %v", key, err)
			continue
		}
		defer res.Body.Close()

		// Read the response body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("error reading service response for %s: %v", key, err)
			continue
		}

		// Unmarshal the response body into a slice of maps
		var tempdata []map[string]interface{}
		err = json.Unmarshal(body, &tempdata)
		if err != nil {
			log.Printf("error unmarshaling JSON for service %s: %v", key, err)
			continue
		}

		// Extract and store Address, ServiceAddress, and ServicePort
		if len(tempdata) > 0 {
			serviceInfo := ServiceInfo{}

			// Extract Address
			if address, ok := tempdata[0]["Address"].(string); ok {
				serviceInfo.Address = address
			}

			// Extract ServiceAddress
			if serviceAddress, ok := tempdata[0]["ServiceAddress"].(string); ok {
				serviceInfo.ServiceAddress = serviceAddress
			}

			// Extract ServicePort
			if servicePort, ok := tempdata[0]["ServicePort"].(float64); ok {
				serviceInfo.ServicePort = int(servicePort)
			}

			// Store the service information in the map
			serviceMap[key] = serviceInfo
		}
	}

	return serviceMap, nil
}

// Function to start the Gin server
func startServer(dataChannel chan Data) {
	r := gin.Default()

	r.GET("/services", func(ctx *gin.Context) {
		fmt.Println("--------- GOT UPDATE MESSAGE ---------")
		ctx.JSON(http.StatusOK, gin.H{"status": "data received"})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
