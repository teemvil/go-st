package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Channel for receiving messages
var messageChannel = make(chan string)

func main() {
	// Set up the MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

	// Create a new client
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Channel for receiving messages
	messageChannel := make(chan string)

	// Start goroutine to handle messages from the channel
	go func() {
		for {
			select {
			case msg := <-messageChannel:
				topic := msg
				if topic != "" {
					subscribeToTopic(client, topic)
				}
			}
		}
	}()

	// Subscribe to initial topic
	subscribeToTopic(client, "management")

	/* Simulate receiving messages
	messageChannel <- "new-topic-1"
	messageChannel <- "trusted"
	messageChannel <- "new-topic-2"*/

	// Wait for new messages to arrive
	select {}
}

// Helper function to subscribe to a topic
func subscribeToTopic(client MQTT.Client, topic string) {
	token := client.Subscribe(topic, 0, messageHandler)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic %s: %s\n", topic, token.Error())
	} else {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

type ManagementMessage struct {
	Name          string    `json: "name"`
	Itemid        string    `json: "itemid"`
	Messsage      string    `json: "message"`
	Event         string    `json: "event"`
	Time          time.Time `json: "time"`
	Jwt           string    `json: "jwt"`
	HostDevice    string    `json: "hostDevice"`
	SensorChannel string    `json: "sensorChannel"`
}

var mes ManagementMessage

// Message handler callback function
func messageHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message on topic: %s\n", msg.Topic())
	fmt.Printf("Message payload: %s\n", msg.Payload())

	err := json.Unmarshal(msg.Payload(), &mes)
	if err != nil {
		fmt.Println(err)
	}
	if mes.Event == "save" {
		messageChannel <- mes.SensorChannel

	}

	if msg.Topic != "management" {
		//TODO: save to database
	}

}
