package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ManagementMessage struct {
	DeviceName       string    `json: "name"`
	Itemid           string    `json: "itemid"`
	Message          string    `json: "message"`
	Event            string    `json: "event"`
	Time             time.Time `json: "time"`
	Jwt              string    `json: "jwt"`
	SensorName       string    `json: "sensorName"`
	SensorHostDevice string    `json: "sensorHostDevice"`
	SensorChannel    string    `json: "sensorChannel"`
	Misc             string    `json: "misc"`
}

type SensorFile struct {
	Name        string `json: "name"`
	HostDevice  string `json: "hostdevice"`
	MQTTchannel string `json: "mqttchannel"`
}

func main() {

	// set up the mqtt client
	opts := MQTT.NewClientOptions()
	//opts.AddBroker("192.168.0.24:1883")
	opts.AddBroker("test.mosquitto.org:1883")

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	//publish token over mqtt

	//get config info from file
	var info SensorFile
	conf, err := os.ReadFile("sensor_config.json")
	if err != nil {
		fmt.Println("error: ", err)
	}

	err2 := json.Unmarshal(conf, &info)
	if err2 != nil {
		fmt.Println(err2)
	}

	var mes = ManagementMessage{info.HostDevice, "", "Hello, " + info.Name + " is online", "sensor-startup", time.Now(), "", info.Name, info.HostDevice, info.MQTTchannel, ""}
	jsonmes, err := json.Marshal(mes)

	mqttToken := client.Publish("management", 0, false, jsonmes)
	mqttToken.Wait()

	//start broadcasting in a loop
	for {
		sensorValue := rand.Intn(100)         //random number for now
		var value = strconv.Itoa(sensorValue) //convert to string
		secondMqttToken := client.Publish(mes.SensorChannel, 0, false, value)
		secondMqttToken.Wait()
		fmt.Printf("Published message: %s\n", value)
		time.Sleep(1 * time.Second)
	}

}
