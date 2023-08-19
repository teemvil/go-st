package main

import (
	//"crypto/x509"

	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	//"github.com/google/go-tpm-tools/client"
	//"github.com/google/go-tpm/tpm2"
	//"github.com/google/go-tpm/tpmutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//jwt "github.com/golang-jwt/jwt"
)

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

type SensorFile struct {
	Name        string `json: "name"`
	HostDevice  string `json: "hostdevice"`
	MQTTchannel string `json: "mqttchannel"`
}

func main() {

	// set up the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

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

	var mes = ManagementMessage{info.Name, "null", "Hello, Temperature sensor here", "sensor-startup", time.Now(), "null", info.HostDevice, info.MQTTchannel}
	jsonmes, err := json.Marshal(mes)

	mqttToken := client.Publish("Management", 0, false, jsonmes)
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
