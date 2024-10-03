package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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
