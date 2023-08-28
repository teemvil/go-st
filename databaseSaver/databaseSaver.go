package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Channel for receiving messages
var messageChannel = make(chan string)

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

// influx data
var token = "D5pA_hs2oPhmCip_EpynBa1GcThw_M29ivprrznLsTM4tE0nPyAnTlb7Zy2ZSSH84nQqVv6YhClxwDzqd8PDzQ=="
var url = "http://localhost:8086"
var client_in = influxdb2.NewClient(url, token)
var org = "metropolia"
var bucket = "testbucket"
var writeAPI = client_in.WriteAPIBlocking(org, bucket)

func main() {
	// Set up the MQTT client options
	opts := MQTT.NewClientOptions()
	//opts.AddBroker("192.168.0.24:1883")
	opts.AddBroker("test.mosquitto.org:1883")

	// Create a new client
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

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

var mes ManagementMessage

// Message handler callback function
func messageHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message on topic: %s\n", msg.Topic())
	fmt.Printf("Message payload: %s\n", msg.Payload())

	topic := string(msg.Topic())

	if topic == "management" {
		err := json.Unmarshal(msg.Payload(), &mes)
		if err != nil {
			fmt.Println(err)
		}
		if mes.Event == "save-data" {
			messageChannel <- mes.SensorChannel

		}
	} else {

		//save to database
		value := string(msg.Payload())
		fmt.Println("value: " + value)
		var valueInt, err2 = strconv.Atoi(value)
		if err2 != nil {
			fmt.Println(err2)
		}
		tags := map[string]string{
			"channel": topic,
		}
		fields := map[string]interface{}{
			"field1": valueInt,
		}
		point := write.NewPoint(topic, tags, fields, time.Now())
		//time.Sleep(1 * time.Second) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}

}
