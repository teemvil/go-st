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
	Name          string    `json: "name"`
	Itemid        string    `json: "itemid"`
	Messsage      string    `json: "message"`
	Event         string    `json: "event"`
	Time          time.Time `json: "time"`
	Jwt           string    `json: "jwt"`
	HostDevice    string    `json: "hostDevice"`
	SensorChannel string    `json: "sensorChannel"`
}

func main() {

	//MQTT client info
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

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

			attestationMes := ManagementMessage{mes.Name, "null", "Validation in progress " + mes.Name, "attest", time.Now(), "null", "null", "null"}
			jsonmes, err := json.Marshal(attestationMes)
			if err != nil {
				fmt.Println(err)
			}

			mqttValidToken := client.Publish("management", 0, false, jsonmes)
			mqttValidToken.Wait()

		}

		if mes.Event == "attest-ok" {

			//secret := []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA88rw9mMriuKHvJ/OE2Bu\noMgrTQ7YyvZi8BOVD2k9cVWaCmYZ/I2nSveMUJuBFyWLMeHgvd97DOpbcxmMtQIj\nDzwjQyueKHuupw4fqXhZ5e2ZDg9ul4aw+yqjBFibZKn5WdD1+zdQpyicWPHe86Z8\n0B0/xs5apHuHtc6IYaHiT/CDs4RkJ2Y3iZPrdnKWGXjHIGUpTYquBQvAQmr8VUvZ\nnZUPAXTAflnziA+31tHUlKICcJXsU6DjacJohI/DbDMKX0zA1UxJwLzD2iXkbZlu\n81cjWBWbZFjZuaT1xEpcj4+gszE8s5iTqh/3jZOCiLFWJzv0V8ikIiP37ASennPB\nawIDAQAB\n-----END PUBLIC KEY-----\n")
			secret := mes.Messsage
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

				validationMes = ManagementMessage{mes.Name, "null", "Validation ok for " + mes.Name, "validation-ok", time.Now(), "null", "null", "null"}
				jsonmes, err := json.Marshal(validationMes)
				if err != nil {
					fmt.Println(err)
				}

				mqttValidToken := client.Publish("management", 0, false, jsonmes)
				mqttValidToken.Wait()
				/*
					fmt.Println("Subject: ", claims["subject"])
					fmt.Println("Name: ", claims["name"])
					fmt.Println("Another: ", claims["another"])*/
			} else {
				fmt.Println("JWT Token is not valid")
			}
		}

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second) //why is this loop here?
	}

}
