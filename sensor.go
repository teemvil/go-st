package sensor

import (
	//"crypto/x509"

	"encoding/json"
	"fmt"
	"os"
	"time"

	//"github.com/google/go-tpm-tools/client"
	//"github.com/google/go-tpm/tpm2"
	//"github.com/google/go-tpm/tpmutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//jwt "github.com/golang-jwt/jwt"
)

func main() {

	// set up the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	//publish token over mqtt
	type Message struct {
		Name       string `json: "name"`
		HostDevice string `json: "hostDevice"`
		Messsage   string `json: "message"`
		Time       int64  `json: "time"`
		Channel    string `json: "channel"`
	}
	var mes = Message{"Temperature sensor", "pi014", "Hello, Temperature sensor here", 32323, "temp1"}
	jsonmes, err := json.Marshal(mes)

	mqttToken := client.Publish("hello-channel", 0, false, jsonmes)
	mqttToken.Wait()

	//start broadcasting in a loop
	for {
		sensorValue := 0 //TODO: get value from sensor
		secondMqttToken := client.Publish(mes.Channel, 0, false, sensorValue)
		secondMqttTtoken.Wait()
		fmt.Printf("Published message: %s\n", text)
		time.Sleep(1 * time.Second)
	}

}
