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
	
	type Message struct {
		name string
		itemid string
		messsage string
		time int64
		jwt string
	}
	var mes Message 

	client := MQTT.NewClient(opts)
	
	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil{
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("sensor-channel", 0, func(client MQTT.Client msg MQTT.Message){
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		err := json.Unmarshal(msg.Payload(), &mes)

		text := string(mes.jwt)
		

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
	}

}