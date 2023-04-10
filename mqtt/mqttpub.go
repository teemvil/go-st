package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	//"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	jwt "github.com/golang-jwt/jwt"
)

func getPubKey(handleaddr int) ([]byte, error) {
	//open TPM2 connection
	rwc, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		fmt.Errorf("can't open TPM at  %v", err)
	}
	defer rwc.Close()

	var handleint = 0
	if handleaddr == 0 {
		handleint = 0x81010003
	} else {
		handleint = handleaddr
	}
	var handle = tpmutil.Handle(handleint)
	fmt.Println("handle: ", handle)
	//there's a EK key at 0x81010002
	//var handleEk = tpmutil.Handle(0x81010002)
	//fmt.Println("handle: ", handleEk)
	//there's a AK key at 0x81010003
	var handleAk = tpmutil.Handle(0x81010003)
	fmt.Println("handle: ", handleAk)

	//read public key directly from handle
	kPublicKey, _, _, err := tpm2.ReadPublic(rwc, handle)
	if err != nil {
		fmt.Errorf("Reading handle failed: %s", err)
	}

	//make public key into into pem format
	ap, err := kPublicKey.Key()
	if err != nil {
		fmt.Errorf("reading Key() failed: %s", err)
	}
	akBytes, err := x509.MarshalPKIXPublicKey(ap)
	if err != nil {
		fmt.Errorf("Unable to convert ekpub: %v", err)
	}
	rakPubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: akBytes,
		},
	)
	fmt.Printf("     PublicKey: \n%v", string(rakPubPEM))
	return rakPubPEM, err
}

func main() {

	//Set the content for jwt-token
	claims := jwt.MapClaims{
		"subject": "test",
		"name":    "name",
		"another": "thing",
	}

	//give the key address memory position the default 0 gives key in 0x81010003
	var secret, err = getPubKey(0)
	if err != nil {
		fmt.Println("error: ", err)
	}

	//set up the jwt token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//and sign with the pem formatted public key
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("signed token: ", signedToken)

	// set up the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.0.24:1883")

	client := MQTT.NewClient(opts)

	if mqttToken := client.Connect(); mqttToken.wait() && mqttToken.Error() != nil {
		fmt.Println(mqttToken.Error())
		os.Exit(1)
	}

	//publish token over mqtt
	mqttToken := client.Publish("test-channel", 0, false, signedToken)
	mqttToken.Wait()

}
