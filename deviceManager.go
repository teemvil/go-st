package manager

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
	
	secret := []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA88rw9mMriuKHvJ/OE2Bu\noMgrTQ7YyvZi8BOVD2k9cVWaCmYZ/I2nSveMUJuBFyWLMeHgvd97DOpbcxmMtQIj\nDzwjQyueKHuupw4fqXhZ5e2ZDg9ul4aw+yqjBFibZKn5WdD1+zdQpyicWPHe86Z8\n0B0/xs5apHuHtc6IYaHiT/CDs4RkJ2Y3iZPrdnKWGXjHIGUpTYquBQvAQmr8VUvZ\nnZUPAXTAflnziA+31tHUlKICcJXsU6DjacJohI/DbDMKX0zA1UxJwLzD2iXkbZlu\n81cjWBWbZFjZuaT1xEpcj4+gszE8s5iTqh/3jZOCiLFWJzv0V8ikIiP37ASennPB\nawIDAQAB\n-----END PUBLIC KEY-----\n")

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

	if mqttToken := client.Subscribe("test-channel", 0, func(client MQTT.Client msg MQTT.Message){
		fmt.Println("Received message on topic: ", msg.Topic())
		fmt.Println("Received message: ", msg.Payload())

		err := json.Unmarshal(msg.Payload(), &mes)

		text := string(mes.jwt)
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
