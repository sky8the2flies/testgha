package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"swift-hub-app/api"
	"swift-hub-app/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var msgCount = 0

// InitMQTTClient initializes and connects to the MQTT broker
func InitMQTTClient(cfg *config.Config, handler mqtt.MessageHandler) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID(cfg.ClientID)
	opts.SetDefaultPublishHandler(handler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	log.Print("Connected to MQTT broker")

	return client
}

// SubscribeToTemperature subscribes to the Zigbee2MQTT topic for temperature data
func SubscribeToTemperature(client mqtt.Client) {
	msgCount++
	topic := "zigbee2mqtt/+"
	token := client.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
		type Payload struct {
			Battery     float32 `json:"battery"`
			Humidity    float32 `json:"humidity"`
			LinkQuality float32 `json:"linkquality"`
			Temperature float32 `json:"temperature"`
			Voltage     float32 `json:"voltage"`
		}
		var p Payload
		err := json.Unmarshal(m.Payload(), &p)
		if err != nil {
			log.Fatalf("Failed to unmmarshal payload err: %+v", err)
			return
		}
		log.Printf("(%d) message recieved %s: %+v", msgCount, m.Topic(), p)
	})
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
}

// MessageHandler handles incoming MQTT messages
func MessageHandler(apiClient *api.Client) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

		// Parse temperature value from message payload
		var temp float64
		if _, err := fmt.Sscanf(string(msg.Payload()), "%f", &temp); err != nil {
			log.Printf("Error parsing temperature data: %v\n", err)
			return
		}

		// Extract device ID from topic
		parts := strings.Split(msg.Topic(), "/")
		deviceID := parts[len(parts)-2]

		// Send data to the API
		err := apiClient.SendTemperatureData(deviceID, temp)
		if err != nil {
			log.Printf("Error sending data to API: %v", err)
		}
	}
}
