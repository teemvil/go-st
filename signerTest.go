package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	tpmjwt "github.com/salrashid123/golang-jwt-tpm"
)

var ()

func main() {

	ctx := context.Background()

	var keyctx interface{}
	/*claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Issuer:    "test",
	}*/

	// Set the content of the jwt-token
	claims := jwt.MapClaims{
		"ExpiresAt": time.Now().Add(time.Minute * 1).Unix(),
		"Issuer":    "test",
		"Another":   "thing",
	}

	//Set the signin method
	tpmjwt.SigningMethodTPMRS256.Override()
	token := jwt.NewWithClaims(tpmjwt.SigningMethodTPMRS256, claims)

	config := &tpmjwt.TPMConfig{
		TPMDevice:     "/dev/tpm0",
		KeyHandleFile: "key.bin",
		KeyTemplate:   tpmjwt.AttestationKeyParametersRSA256,
		//KeyTemplate: tpmjwt.UnrestrictedKeyParametersRSA256,
	}

	keyctx, err := tpmjwt.NewTPMContext(ctx, config)
	if err != nil {
		log.Fatalf("Unable to initialize tpmJWT: %v", err)
	}

	fmt.Println("key context: ", keyctx)

	token.Header["kid"] = config.GetKeyID()
	rrr := config.GetKeyID()
	fmt.Println("key id: ", rrr)
	tokenString, err := token.SignedString(keyctx)
	if err != nil {
		log.Fatalf("Error signing %v", err)
	}
	fmt.Printf("TOKEN: %s\n", tokenString)

	// verify with TPM based publicKey
	keyFunc, err := tpmjwt.TPMVerfiyKeyfunc(ctx, config)
	if err != nil {
		log.Fatalf("could not get keyFunc: %v", err)
	}

	vtoken, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Fatalf("Error verifying token %v", err)
	}
	if vtoken.Valid {
		log.Println("     verified with TPM PublicKey")
	}

	// verify with provided RSAPublic key
	pubKey := config.GetPublicKey()
	fmt.Println("pub key: ", pubKey)

	//yetAnotherkey, err := ioutil.ReadFile("key.pem")
	//fmt.Println("pub key ?: ", yetAnotherkey)

	//pp := "ZDcxYjczM2NmMTdiOTQ5MWRjNzdiZDIyMTcyYTU0YzU2ZGI5NWVmMWU3NzE0NWRhNDgwZjhmZWMyYjUwZTI4MQ=="

	v, err := jwt.Parse(vtoken.Raw, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if v.Valid {
		log.Println("     verified with exported PubicKey")
	}

	/*parsedToken, err := jwt.Parse(text, func(parsedToken *jwt.Token)(interface{}, error){
		return secret, nil
	})
	*/
	if err != nil {
		fmt.Println("Error: ", err)
	}
	//Print the content
	if claims, ok := v.Claims.(jwt.MapClaims); ok && v.Valid {
		fmt.Println("Decoded JWT claims:")
		fmt.Println("Expires: ", claims["ExpiresAt"])
		fmt.Println("Name: ", claims["Issuer"])
		fmt.Println("Another: ", claims["Another"])
	} else {
		fmt.Println("JWT Token is not valid")
	}

}
