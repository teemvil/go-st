package main

import (
	//"crypto/x509"

	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	//"github.com/google/go-tpm-tools/client"
	//"github.com/google/go-tpm/tpm2"
	//"github.com/google/go-tpm/tpmutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//jwt "github.com/golang-jwt/jwt"
)

func main() {

	// set up the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	//publish token over mqtt
	type ManagementMessage struct {
		Name          string    `json: "name"`
		Itemid        string    `json: "itemid"`
		Messsage      string    `json: "message"`
		Event         string    `json: "event"`
		Time          time.Time `json: "time"`
		Jwt           string    `json: "jwt"`
		HostDevice    string    `json: "hostDevice"`
		SensorChannel string    `json: "channel"`
	}

	//get config info from file
	type SensorFile struct {
		Name        string `json: "name"`
		HostDevice  string `json: "host-device"`
		MQTTchannel string `json: "mqtt-channel"`
	}

	var info SensorFile
	conf, err := os.ReadFile("device_config.json")
	if err != nil {
		fmt.Println("error: ", err)
	}

	err2 := json.Unmarshal(conf, &info)
	if err2 != nil {
		fmt.Println(err2)
	}

	var mes = ManagementMessage{info.Name, "null", "Hello, Temperature sensor here", "sensor-startup", time.Now(), "null", info.HostDevice, info.MQTTchannel}
	jsonmes, err := json.Marshal(mes)

	mqttToken := client.Publish("Management", 0, false, jsonmes)
	mqttToken.Wait()

	//start broadcasting in a loop
	for {
		sensorValue := rand.Intn(100) //random number for now
		secondMqttToken := client.Publish(mes.Channel, 0, false, sensorValue)
		secondMqttToken.Wait()
		fmt.Printf("Published message: %s\n", sensorValue)
		time.Sleep(1 * time.Second)
	}

}
