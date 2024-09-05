package main

import (
	"bytes"
	"encoding/binary"
)

// Go representation of the struct catalog_key
type CatalogKey struct {
	Hostname [256]byte // Define a fixed-size array for hostname
}

// Go representation of the struct catalog_value
type CatalogValue struct {
	ServiceIP uint32 // Resolved IP provided by consul
}

// Function to convert string hostname to CatalogKey
func NewCatalogKey(hostname string) CatalogKey {
	var key CatalogKey
	copy(key.Hostname[:], hostname)
	return key
}

// Function to convert byte array to CatalogValue (for map lookup)
func BytesToCatalogValue(data []byte) (CatalogValue, error) {
	var value CatalogValue
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &value)
	return value, err
}

// Function to convert CatalogValue to byte array (for map update)
func CatalogValueToBytes(value CatalogValue) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	return buf.Bytes(), err
}
