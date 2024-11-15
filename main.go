// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// // Configure these based on your setup
// const (
// 	apiBaseURL = "http://localhost:8080/api" // Replace with your deCONZ API URL
// 	apiKey     = "YOUR_API_KEY_HERE"         // Replace with your deCONZ API Key
// )

// // Light represents a light in the deCONZ system
// type Light struct {
// 	Name  string `json:"name"`
// 	Type  string `json:"type"`
// 	State struct {
// 		On bool `json:"on"`
// 	} `json:"state"`
// }

// // FetchLights fetches a list of lights from deCONZ
// func FetchLights() (map[string]Light, error) {
// 	url := fmt.Sprintf("%s/%s/temperature", apiBaseURL, apiKey)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching lights: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %w", err)
// 	}

// 	lights := make(map[string]Light)
// 	if err := json.Unmarshal(body, &lights); err != nil {
// 		return nil, fmt.Errorf("error unmarshalling lights: %w", err)
// 	}

// 	return lights, nil
// }

// // ToggleLight toggles the state of a specific light
// func ToggleLight(lightID string, on bool) error {
// 	url := fmt.Sprintf("%s/%s/lights/%s/state", apiBaseURL, apiKey, lightID)
// 	data := map[string]bool{"on": on}
// 	payload, err := json.Marshal(data)
// 	if err != nil {
// 		return fmt.Errorf("error marshalling toggle data: %w", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return fmt.Errorf("error toggling light: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error toggling light: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		return fmt.Errorf("error toggling light, status: %s, response: %s", resp.Status, string(body))
// 	}

// 	return nil
// }

// func main() {
// 	// Fetch and display lights
// 	lights, err := FetchLights()
// 	if err != nil {
// 		log.Fatalf("Failed to fetch lights: %v", err)
// 	}

// 	for id, light := range lights {
// 		fmt.Printf("Light ID: %s, Name: %s, Type: %s, State: %t\n", id, light.Name, light.Type, light.State.On)
// 	}

// 	// Example: Toggle the first light found
// 	for id, light := range lights {
// 		fmt.Printf("Toggling light %s (%s)...\n", id, light.Name)
// 		newState := !light.State.On
// 		if err := ToggleLight(id, newState); err != nil {
// 			log.Fatalf("Failed to toggle light %s: %v", id, err)
// 		}
// 		fmt.Printf("Light %s is now set to %t\n", id, newState)
// 		break
// 	}
// }

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
