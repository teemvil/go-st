package sensor

import (
	//"crypto/x509"

	"encoding/json"
	"fmt"
	"os"

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
		name       string
		hostDevice string
		messsage   string
		time       int64
		channel    string
	}
	var mes = Message{"Temperature sensor", "pi014", "Hello, Temperature sensor here", 32323, "temp1"}
	jsonmes, err := json.Marshal(mes)

	mqttToken := client.Publish("sensor-channel", 0, false, jsonmes)
	mqttToken.Wait()

}
