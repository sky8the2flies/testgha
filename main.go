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
		ReadTimeout: time.Second,
	}

	// Open serial port
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open port: %v", err)
	}
	defer port.Close()

	// Buffer for reading data
	readBuffer := make([]byte, 128)
	messageBuffer := []byte{}

	fmt.Println("Reading from serial port...")
	for {
		n, err := port.Read(readBuffer)
		if err != nil {
			log.Printf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			// Append new data to the message buffer
			messageBuffer = append(messageBuffer, readBuffer[:n]...)

			// Check for start and end delimiters (192) to identify a complete message
			startIndex := -1
			endIndex := -1

			for i, b := range messageBuffer {
				if b == 192 {
					if startIndex == -1 {
						startIndex = i
					} else {
						endIndex = i
						break
					}
				}
			}

			// If we found a complete message (from start to end delimiter)
			if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
				// Extract the complete message
				payload := messageBuffer[startIndex : endIndex+1]

				fmt.Println("Raw Data:", payload)

				parseZigbeeMessage(payload)

				// Remove the processed message from the buffer
				messageBuffer = messageBuffer[endIndex+1:]
			}
		}
	}
}

// Parsing function for Zigbee message structure
func parseZigbeeMessage(data []byte) {
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
