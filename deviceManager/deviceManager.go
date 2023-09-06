package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	jwt "github.com/golang-jwt/jwt"
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

func main() {

	//MQTT client info
	opts := MQTT.NewClientOptions()
	//opts.AddBroker("192.168.0.24:1883")
	opts.AddBroker("test.mosquitto.org:1883")

	var mes ManagementMessage

	var validationMes ManagementMessage

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("management", 0, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		jsonerr := json.Unmarshal(msg.Payload(), &mes)
		if jsonerr != nil {
			fmt.Println("JSON Error: ", jsonerr)
		}

		//attestation
		itemid := mes.Itemid
		//TODO: attestation using itemid, to get secret from attestation database
		if mes.Event == "device-startup" {
			//secret = attest(itemid)
			fmt.Println(itemid)

			attestationMes := ManagementMessage{mes.DeviceName, mes.Itemid, "Validation in progress " + mes.DeviceName, "attest", time.Now(), mes.Jwt, "", "", "", ""}
			jsonmes, err := json.Marshal(attestationMes)
			if err != nil {
				fmt.Println(err)
			}

			mqttValidToken := client.Publish("management", 0, false, jsonmes)
			mqttValidToken.Wait()

		}

		if mes.Event == "attest-ok" {

			secret := []byte(mes.Misc)

			//if attestation succesful, we should be able to open the jwt
			messageJwt := string(mes.Jwt)
			parsedToken, err := jwt.Parse(messageJwt, func(parsedToken *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil {
				fmt.Println("Error parsing jwt: ", err)
			}
			if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
				fmt.Println("JWT is valid")

				validationMes = ManagementMessage{mes.DeviceName, "", "Validation ok for " + mes.DeviceName, "validation-ok", time.Now(), mes.Jwt, "", "", "", ""}
				jsonmes, err := json.Marshal(validationMes)
				if err != nil {
					fmt.Println(err)
				}

				mqttValidToken := client.Publish("management", 0, false, jsonmes)
				mqttValidToken.Wait()

				fmt.Println("Subject: ", claims["subject"])
				fmt.Println("Name: ", claims["name"])
				fmt.Println("Another: ", claims["another"])
			} else {
				fmt.Println("JWT Token is not valid")
				validationMes = ManagementMessage{mes.DeviceName, "", " JWT validation failed for " + mes.DeviceName, "validation-fail", time.Now(), mes.Jwt, "", "", "", ""}
				jsonmes, err := json.Marshal(validationMes)
				if err != nil {
					fmt.Println(err)
				}

				mqttValidToken := client.Publish("management", 0, false, jsonmes)
				mqttValidToken.Wait()
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
