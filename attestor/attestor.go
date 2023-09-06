package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

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

func main() {

	//MQTT client info
	opts := MQTT.NewClientOptions()
	//opts.AddBroker("192.168.0.24:1883")
	opts.AddBroker("test.mosquitto.org:1883")

	//open database connection
	db, err := sql.Open("mysql", "user:password@/attestation")
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

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
		//attestation using itemid, to get secret from attestation database
		if mes.Event == "attestation-start" {
			fmt.Println(itemid)

			//get public key from database
			smtOut, err := db.Prepare("SELECT pubkey FROM devices WHERE itemid = ?")
			if err != nil {
				fmt.Println(err)
			}
			defer smtOut.Close()

			var res string
			err2 := smtOut.QueryRow(itemid).Scan(&res)
			//if unsucceful, send out attest fail message, else send attest success message
			if err2 != nil {
				fmt.Println("Error getting pubkey from database", err2)
				validationMes = ManagementMessage{mes.DeviceName, mes.Itemid, "Attestation failed for " + mes.DeviceName, "attest-fail", time.Now(), mes.Jwt, "", "", "", ""}
				jsonmes, err := json.Marshal(validationMes)
				if err != nil {
					fmt.Println(err)
				}

				mqttValidToken := client.Publish("management", 0, false, jsonmes)
				mqttValidToken.Wait()
			} else {
				fmt.Println(res)
				validationMes = ManagementMessage{mes.DeviceName, mes.Itemid, "Attestation ok for " + mes.DeviceName, "attest-ok", time.Now(), mes.Jwt, "", "", "", res}
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
