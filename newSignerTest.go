package main

import (
	"fmt"
	//"os"
	"io/ioutil"

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

	//there's a AK key at 0x81010003
	var handle = tpmutil.Handle(0x81010003)
	fmt.Println("handle: ", handle)

	pubKey, err := ioutil.ReadFile("keys/ak.pem")
	fmt.Println("pub ak key : ", pubKey)

	secret := []byte(pubKey)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("signed token: ", signedToken)

	parsedToken, err := jwt.Parse(signedToken, func(parsedToken *jwt.Token) (interface{}, error) {
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

}
