package main

import (
	"log"
	"os"
	"os/signal"
	"swift-hub-app/api"
	"swift-hub-app/config"
	"swift-hub-app/mqtt"
	"syscall"
)

func main() {
	log.Print("Starting")
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize API client
	apiClient := api.NewClient(cfg)

	// Initialize MQTT client
	mqttClient := mqtt.InitMQTTClient(cfg, mqtt.MessageHandler(apiClient))

	// Subscribe to temperature topics
	mqtt.SubscribeToTemperature(mqttClient)

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // Wait for termination signal

	log.Println("Shutting down...")
	mqttClient.Disconnect(250)
}
