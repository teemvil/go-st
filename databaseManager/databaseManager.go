package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	var mes ManagementMessage

	//list of validated devices
	var devices []string

	//list of active sensors
	type Sensor struct {
		Name, Hostdevice, Channel string
	}
	var sensors []Sensor

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.Wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("management", 0, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		err := json.Unmarshal(msg.Payload(), &mes)
		if err != nil {
			fmt.Println(err)
		}

		//add new sensors to the list
		if string(mes.Event) == "sensor-startup" {
			var sensor = Sensor{mes.SensorName, mes.SensorHostDevice, mes.SensorChannel}
			sensors = append(sensors, sensor)
			fmt.Println("sensor " + mes.SensorName + " added to the list")
			//Check if hostdevice is deemed safe
			n := 0
			for n < len(devices) {
				fmt.Println(devices[n])
				if mes.SensorHostDevice == devices[n] {
					//send message to databaseSaver.go to start saving
					var savemessage = ManagementMessage{mes.DeviceName, "", "Start saving from channel " + mes.SensorChannel, "save", time.Now(), "", mes.SensorName, mes.SensorHostDevice, mes.SensorChannel, ""}
					jsonmes, err := json.Marshal(savemessage)
					if err != nil {
						fmt.Println(err)
					}
					sToken := client.Publish("management", 0, false, jsonmes)
					sToken.Wait()

					//n = len(devices) - 1 //kludge for ending the loop
				}
				n = n + 1
			}
		}

		//add new validated device to the list
		if string(mes.Event) == "validation-ok" {
			devices = append(devices, mes.DeviceName)
			fmt.Println("device " + mes.DeviceName + " added to the list")
			//check if there are sensors on the device
			n := 0
			for n < len(sensors) {
				fmt.Println(sensors[n])
				if mes.DeviceName == sensors[n].Hostdevice {
					// send message to databaseSaver.go to start saving
					var savemessage = ManagementMessage{sensors[n].Hostdevice, "", "Start saving from channel " + sensors[n].Channel, "save", time.Now(), "", sensors[n].Name, sensors[n].Hostdevice, sensors[n].Channel, ""}
					jsonmes, err := json.Marshal(savemessage)
					if err != nil {
						fmt.Println(err)
					}
					sToken := client.Publish("management", 0, false, jsonmes)
					sToken.Wait()
				}
				n = n + 1
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
