package main

import (

	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	jwt "github.com/golang-jwt/jwt"
)

func main(){

	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")
	
	secret := []byte("secret")

	client := MQTT.NewClient(opts)
	
	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil{
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("test-channel", 0, func(client MQTT.Client msg MQTT.Message){
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		text := string(msg.Payload())
		parsedToken, err := jwt.Parse(text, func(parsedToken *jwt.Token)(interface{}, error){
			return secret, nil
		})

		if err != nil {
			fmt.Println("Error: ", err)
		}
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			fmt.Println("Decoded JWT claims:")
			fmt.Println("Subject: ", claims["subject"])
			fmt.Println("Name: ", claims["name"])
			fmt.Println("Another: ", claims["another"])
		} else {
			fmt.Println("JWT Token is not valid")
		}

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
	}

}
