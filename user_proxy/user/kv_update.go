package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func update_kv_store(key string, value string) error {

	// make url
	url := fmt.Sprintf("http://127.0.0.1:8500/v1/kv/%s", key)

	fmt.Println(url, key)
	fmt.Println(value)

	// Create a new PUT request with the value as the request body
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ------------------ CUSTOMIZE : ADD HELPER FUNCTIONS TO UPDATE THE KV-STORE AS PER REQUIREMENT (below is one example) ------------------

// KV-Store naming convection
// metric/<node-name>/<metric-type>
// Eg : metric/node1/cpu, metric/node1/memory

func compute_metrics(node_name string) {
	for {
		percentages, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Printf("Error retrieving CPU utilization: %s\n", err)
			continue
		}

		// percentages will hold the CPU utilization percentage for each CPU
		// since the second parameter is false, it aggregates over all CPUs
		// fmt.Printf("CPU Utilization: %.2f%%\n", percentages[0])

		// Get virtual memory stats
		vMemStat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error retrieving memory stats:", err)
			return
		}

		// Calculate memory utilization
		memUtilization := 100.0 * float64(vMemStat.Used) / float64(vMemStat.Total)
		cpuUtilization := percentages[0]

		fmt.Printf("Memory Utilization: %.2f%%\n", memUtilization)

		update_kv_store(fmt.Sprintf("metric/%s/cpu", node_name), fmt.Sprintf("%f", cpuUtilization))
		update_kv_store(fmt.Sprintf("metric/%s/memory", node_name), fmt.Sprintf("%f", memUtilization))

		// Sleep for a bit before the next reading (optional, adjust as needed)
		time.Sleep(time.Second * 10)
	}
}
