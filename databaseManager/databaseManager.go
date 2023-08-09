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
		SensorChannel    string `json: "channel"`
	}
	var mes ManagementMessage 

	//list of validated devices
	var devices = []string

	//list of active sensors
	var sensors = []string

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
			sensor = mes
			sensors[].add(sensor)
		}

		if (string(mes.Event)="validaton-ok"){
			devices[].add(mes.Name)
		}


		

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
	}

}