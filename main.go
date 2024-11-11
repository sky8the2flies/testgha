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
			// Print received data
			fmt.Printf("Received data: %s\n", buffer[:n])
		}
	}
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
