package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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
		if mes.Event == "attestation-start" {
			//secret = attest(itemid)
			fmt.Println(itemid)

			//TODO: get public key from somewhere
			secret := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA88rw9mMriuKHvJ/OE2Bu\noMgrTQ7YyvZi8BOVD2k9cVWaCmYZ/I2nSveMUJuBFyWLMeHgvd97DOpbcxmMtQIj\nDzwjQyueKHuupw4fqXhZ5e2ZDg9ul4aw+yqjBFibZKn5WdD1+zdQpyicWPHe86Z8\n0B0/xs5apHuHtc6IYaHiT/CDs4RkJ2Y3iZPrdnKWGXjHIGUpTYquBQvAQmr8VUvZ\nnZUPAXTAflnziA+31tHUlKICcJXsU6DjacJohI/DbDMKX0zA1UxJwLzD2iXkbZlu\n81cjWBWbZFjZuaT1xEpcj4+gszE8s5iTqh/3jZOCiLFWJzv0V8ikIiP37ASennPB\nawIDAQAB\n-----END PUBLIC KEY-----\n"

			validationMes = ManagementMessage{mes.Name, "null", secret, "attest-ok", time.Now(), mes.Jwt, "null", "null"}
			jsonmes, err := json.Marshal(validationMes)
			if err != nil {
				fmt.Println(err)
			}

			mqttValidToken := client.Publish("management", 0, false, jsonmes)
			mqttValidToken.Wait()

		}

	}); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second)
	}

}
