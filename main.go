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
			// Convert the buffer to a string and print it
			fmt.Printf("Received data: %v\n", buffer[:n])
		}
	}
}

// Sample parsing function for Zigbee message structure
func parseZigbeeMessage(data []byte) {
	// This is just a placeholder example for parsing
	// Actual parsing will depend on the structure of the data from RaspBee II

	if len(data) < 8 {
		fmt.Println("Invalid message length")
		return
	}

	// Assuming a hypothetical structure: header, message type, device ID, payload, footer
	header := data[0]
	messageType := data[1]
	deviceID := data[2]
	payload := data[3 : len(data)-1]
	footer := data[len(data)-1]

	fmt.Printf("Parsed Message:\n")
	fmt.Printf("  Header: %02X\n", header)
	fmt.Printf("  Message Type: %02X\n", messageType)
	fmt.Printf("  Device ID: %02X\n", deviceID)
	fmt.Printf("  Payload: % X\n", payload)
	fmt.Printf("  Footer: %02X\n", footer)
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
