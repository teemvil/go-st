package main

import (

	"fmt"
	"os"
	"time"
	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	jwt "github.com/golang-jwt/jwt"
)

func main(){

	//MQTT client info
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")
	
	type ManagementMessage struct {
		Name     string `json: "name"`
		Itemid   string `json: "itemid"`
		Messsage string `json: "message"`
		Event	 string `json: "event"`
		Time     time.Time  `json: "time"`
		Jwt      string `json: "jwt"`
		HostDevice string `json: "hostDevice"`
		SensorChannel    string `json: "sensorChannel"`
	}
	var mes ManagementMessage 

	//list of validated devices
	var devices = []string

	//list of active sensors
	type Sensor struct {
		name, hostdevice, channel string
	}
	var sensors = []Sensor

	client := MQTT.NewClient(opts)
	
	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil{
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("management", 0, func(client MQTT.Client msg MQTT.Message){
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		err := json.Unmarshal(msg.Payload(), &mes)

		//add new sensors to the list
		if (string(mes.Event)="sensor-startup"){
			sensor = Sensor(mes.Name, mes.HostDevice, mes.SensorChannel)
			sensors = append(sensors, sensor)
			fmt.Println("sensor "+mes.Name+" added to the list")
			n := 0
			for n < len(devices) {
				fmt.Println(devices[n])
				if (mes.Name = devices[n]){
					//TODO: start recording sensor data from sensor's channel

					n=len(devices)-1 //kludge for ending the loop
				}
				n = n+1
			}
		}

		//add new validated device to the list
		if (string(mes.Event)="validation-ok"){
			devices = append(devices, mes.Name)
			fmt.Println("device "+mes.Name+" added to the list")
			n := 0
			for n < len(sensors) {
				fmt.Println(sensors[n])
				if (mes.Name = sensors[n].hostdevice){
					//TODO: start recording sensor data from sensor's channel

					subscribeToTopic(client, sensors[n].channel) //??

					n=len(sensors)-1 //kludge for ending the loop
				}
				n = n+1
			}
		}
		

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
	}

}



func subscribeToTopic(client MQTT.Client, topic string) {
	token := client.Subscribe(topic, 0, messageHandler)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic %s: %s\n", topic, token.Error())
	} else {
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

func messageHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message on topic: %s\n", msg.Topic())
	fmt.Printf("Message payload: %s\n", msg.Payload())
}