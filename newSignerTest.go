package main

import (
	"fmt"
	//"os"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	//MQTT "github.com/eclipse/paho.mqtt.golang"

	//"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	jwt "github.com/golang-jwt/jwt"
)

func main() {

	//Set the content for jwt-token
	claims := jwt.MapClaims{
		"subject": "test",
		"name":    "name",
		"another": "thing",
	}

	rwc, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		fmt.Errorf("can't open TPM at  %v", err)
	}
	defer rwc.Close()

	//there's a EK key at 0x81010002
	var handleEk = tpmutil.Handle(0x81010002)
	fmt.Println("handle: ", handleEk)
	//there's a AK key at 0x81010003
	var handleAk = tpmutil.Handle(0x81010003)
	fmt.Println("handle: ", handleAk)

	//testing with pem reading
	pubKey, err := ioutil.ReadFile("keys/ak2.pem")
	fmt.Println("pub ak key : \n", pubKey)
	ss := string(pubKey)
	fmt.Println("key string : ", ss)
	//secret := []byte(pubKey)
	srs := []byte(ss)
	fmt.Println("byte string : ", srs)

	//read public key directly from handle
	kPublicKey, _, _, err := tpm2.ReadPublic(rwc, handleAk)
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

	//set up the jwt token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//and sign with the pem formatted public key
	signedToken, err := jwtToken.SignedString(rakPubPEM)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("signed token: ", signedToken)

	//confirm public key manually
	pKeyFromAtt := []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA88rw9mMriuKHvJ/OE2Bu\noMgrTQ7YyvZi8BOVD2k9cVWaCmYZ/I2nSveMUJuBFyWLMeHgvd97DOpbcxmMtQIj\nDzwjQyueKHuupw4fqXhZ5e2ZDg9ul4aw+yqjBFibZKn5WdD1+zdQpyicWPHe86Z8\n0B0/xs5apHuHtc6IYaHiT/CDs4RkJ2Y3iZPrdnKWGXjHIGUpTYquBQvAQmr8VUvZ\nnZUPAXTAflnziA+31tHUlKICcJXsU6DjacJohI/DbDMKX0zA1UxJwLzD2iXkbZlu\n81cjWBWbZFjZuaT1xEpcj4+gszE8s5iTqh/3jZOCiLFWJzv0V8ikIiP37ASennPB\nawIDAQAB\n-----END PUBLIC KEY-----\n")

	parsedToken, err := jwt.Parse(signedToken, func(parsedToken *jwt.Token) (interface{}, error) {
		return pKeyFromAtt, nil
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

}
