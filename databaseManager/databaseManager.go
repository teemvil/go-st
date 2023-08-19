package main

import (

	"fmt"
	"os"
	"time"
	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	jwt "github.com/golang-jwt/jwt"
)


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

func main(){

	//MQTT client info
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")
	
	var mes ManagementMessage 

	//list of validated devices
	var devices = []string

	//list of active sensors
	type Sensor struct {
		name, hostdevice, channel string
	}
	var sensors = []Sensor

	client := MQTT.NewClient(opts)
	
	if mqttToken := client.Connect(); mqttToken.Wait() && mqttToken.Error() != nil{
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	if mqttToken := client.Subscribe("management", 0, func(client MQTT.Client, msg MQTT.Message){
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		err := json.Unmarshal(msg.Payload(), &mes)
		if err != nil{
			fmt.Println(err)
		}

		//add new sensors to the list
		if (string(mes.Event)=="sensor-startup"){
			sensor = Sensor(mes.Name, mes.HostDevice, mes.SensorChannel)
			sensors = append(sensors, sensor)
			fmt.Println("sensor "+mes.Name+" added to the list")
			//Check if hostdevice is deemed safe
			n := 0
			for n < len(devices) {
				fmt.Println(devices[n])
				if (mes.Name = devices[n]){
					//TODO: start recording sensor data from sensor's channel
					// RAther: send message to databaseSaver.go to start saving
					savemessage = ManagementMessage(mes.Name, "null", "Start saving from channel "+mes.SensorChannel, "save", time.Now(), "null", mes.HostDevice, mes.SensorChannel)
					jsonmes, err := json.Marshal(savemessage)

					sToken := client.Publish("management", 0, false, jsonmes)
					sToken.Wait()

					n=len(devices)-1 //kludge for ending the loop
				}
				n = n+1
			}
		}

		//add new validated device to the list
		if (string(mes.Event)="validation-ok"){
			devices = append(devices, mes.Name)
			fmt.Println("device "+mes.Name+" added to the list")
			//check if there are sensors on the device
			n := 0
			for n < len(sensors) {
				fmt.Println(sensors[n])
				if (mes.Name = sensors[n].hostdevice){
					//TODO: start recording sensor data from sensor's channel
					// RAther: send message to databaseSaver.go to start saving
					savemessage = ManagementMessage(sensors[n].name, "null", "Start saving from channel "+sensors[n].channel, "save", time.Now(), "null", sensors[n].hostdevice, sensors[n].channel)
					jsonmes, err := json.Marshal(savemessage)

					sToken := client.Publish("management", 0, false, jsonmes)
					sToken.Wait()
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

