package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/gin-gonic/gin"
)

// Data struct to hold incoming JSON data
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// Channel to communicate between goroutines
	dataChannel := make(chan Data)

	// Remove resource limits for kernels <5.11.
	// if err := rlimit.RemoveMemlock(); err != nil {
	// 	log.Fatal("Removing memlock:", err)
	// }

	spec, err := ebpf.LoadCollectionSpec("../kernel/xdp_main")
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

	for {
		// Handle map configuration and handling incoming information
	}

	// Main goroutine handles data from the channel
	// for {
	// 	select {
	// 	case data := <-dataChannel:
	// 		// Perform operations with the data
	// 		log.Printf("Handling data: %s = %s", data.Key, data.Value)

	// 		// Add your data handling logic here
	// 	case <-signal.Notify(make(chan os.Signal, 1), os.Interrupt):
	// 		log.Print("Received signal, exiting..")
	// 		return
	// 	}
	// }
}

// Function to start the Gin server
func startServer(dataChannel chan Data) {
	r := gin.Default()

	r.POST("/data", func(c *gin.Context) {
		var data Data
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send the data to the channel
		dataChannel <- data

		c.JSON(http.StatusOK, gin.H{"status": "data received"})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
