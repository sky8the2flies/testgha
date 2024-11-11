package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	// Configure serial connection settings
	config := &serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        38400,
		Size:        8,
		StopBits:    serial.Stop1,
		Parity:      serial.ParityNone,
		ReadTimeout: time.Second * 5,
	}

	// Open serial port
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open port: %v", err)
	}
	defer port.Close()

	// Buffer for reading data
	buffer := make([]byte, 128)

	fmt.Println("Reading from serial port...")
	for {
		n, err := port.Read(buffer)
		if err != nil {
			log.Printf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			payload := buffer[:n]
			fmt.Println("Raw Data:", buffer[:n])

			crc0, crc1 := CalculateChecksum(payload)
			fmt.Println("Checksum: ", crc0, crc1)
		}
	}
}

// CalculateChecksum computes the 16-bit checksum for the given payload
func CalculateChecksum(payload []byte) (uint8, uint8) {
	var crc uint16 = 0

	// Sum each byte in the payload
	for _, b := range payload {
		crc += uint16(b)
	}

	// Two's complement and split into two bytes
	crc = ^crc + 1
	crc0 := uint8(crc & 0xFF)        // Lower byte
	crc1 := uint8((crc >> 8) & 0xFF) // Upper byte

	return crc0, crc1
}

// Parsing function for Zigbee message structure
func parseZigbeeMessage(data []byte) {
	if len(data) < 3 || data[0] != 192 || data[len(data)-1] != 192 {
		fmt.Println("Invalid message format")
		return
	}

	// Extract fields based on observed pattern
	startByte := data[0]
	endByte := data[len(data)-1]

	messageType := data[1]           // 2nd byte, possible message type or command
	deviceID := data[2]              // 3rd byte, could be device identifier
	payload := data[3 : len(data)-1] // Remaining data except start and end bytes

	// Interpret payload as fields or a single integer (depends on message format)
	fmt.Printf("Parsed Message:\n")
	fmt.Printf("  Start Byte: %02X\n", startByte)
	fmt.Printf("  Message Type: %02X\n", messageType)
	fmt.Printf("  Device ID: %02X\n", deviceID)
	fmt.Printf("  Payload: %v\n", payload)
	fmt.Printf("  End Byte: %02X\n\n", endByte)
}

// package main

// import (
// 	"log"
// 	"os"
// 	"os/signal"
// 	"swift-hub-app/api"
// 	"swift-hub-app/config"
// 	"swift-hub-app/mqtt"
// 	"syscall"
// )

// func main() {
// 	log.Print("Starting")
// 	// Load configuration
// 	cfg := config.LoadConfig()

// 	// Initialize API client
// 	apiClient := api.NewClient(cfg)

// 	// Initialize MQTT client
// 	mqttClient := mqtt.InitMQTTClient(cfg, mqtt.MessageHandler(apiClient))

// 	// Subscribe to temperature topics
// 	mqtt.SubscribeToTemperature(mqttClient)

// 	// Set up signal handling for graceful shutdown
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	<-sigs // Wait for termination signal

// 	log.Println("Shutting down...")
// 	mqttClient.Disconnect(250)
// }
