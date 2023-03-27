package main

import (
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"crypto/x509"

	"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"

	"encoding/base64"
	//"encoding/hex"
	"encoding/json"

	"crypto/sha256"

	"github.com/lestrrat-go/jwx/jwk"
)

var (
	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket)")
)

func main() {

	rwc, err := tpm2.OpenTPM(*tpmPath)
	if err != nil {
		fmt.Errorf("can't open TPM at  %v", err)
	}
	defer rwc.Close()

	nonce := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	var handle = tpmutil.Handle(0x81010003)

	var k = tpmutil.Handle(0x81010001)

	var keyFile = flag.String("keyFile", "key.bin", "TPM KeyFile")
	var publicKeyFile = flag.String("publicKeyFile", "key.pem", "PEM File to write the public key")

	fmt.Println("Handle is", handle, k)
	fmt.Println("Type:", reflect.TypeOf(handle).String())
	fmt.Println("Type:", reflect.TypeOf(k).String())

	// This is a test thing from the TPMcourse..
	// This bit only works if you've run the createkeys.sh in the keys folder,
	// which creates the handle at 0x81010003. Otherwise it'll give a nil
	// quote retuns 3 values,  att, sig, err
	att, sig, err := tpm2.Quote(
		rwc,
		handle,
		"",
		"",
		nonce,
		tpm2.PCRSelection{tpm2.AlgSHA256, []int{0}},
		tpm2.AlgNull)

	if err != nil {
		fmt.Errorf("Problem getting quote  %s", err)
		fmt.Println(err)
	}
	if sig != nil {
		fmt.Errorf("Sig is nil\n")
	}
	if att != nil {
		fmt.Errorf("Att is nil\n")
	}
	// att is of type []byte
	// sig is of type tpm2.Signature ??
	// err is ??
	fmt.Println("Err is *", err, "*")
	fmt.Println("length  ", len(att))
	fmt.Println("Att [% x] ", att)
	fmt.Println("Sig  ", sig)

	// Endorsementkey method creates handle at 0x81010001
	ekkey, err := client.EndorsementKeyRSA(rwc)
	fmt.Println("EK key (??) : ", ekkey)
	//fmt.Println("EK public key: ", ekkey.PublicKey())
	fmt.Println("EK handle: ", ekkey.Handle())
	//fmt.Println("EK public area: ", ekkey.PublicArea().Attributes)
	ee, err := ekkey.PublicArea().Key()
	en, err := ekkey.PublicArea().Name()
	//fmt.Println("EK public area key: ", ee)
	//fmt.Println("EK public area name: ", en)

	// Attstionkey creates handle at 0x81008F01
	akkey, err := client.AttestationKeyRSA(rwc)
	//fmt.Println("AK key (??) : ", akkey)
	//fmt.Println("AK public key : ", akkey.PublicKey())

	//keyParameters = client.AKTemplateRSA()
	/*
		keey, err := client.EndorsementKeyFromNvIndex(rwc, 0x81010003)
		if err != nil {
			fmt.Println("cattestatko key not: ", err)
		}
		fmt.Println("ek kssey ?: ", keey)
		fmt.Println("ek kssey ?: ", keey.PublicKey())
	*/

	// Newkey method does not persist the handle
	anotherkey, err := client.NewKey(rwc, tpm2.HandleEndorsement, client.AKTemplateRSA())
	if err != nil {
		fmt.Println("can't create SRK %q: %v", tpmPath, err)
	}

	kh := anotherkey.Handle()
	fmt.Println("key handle  ", kh)
	fmt.Println("======= ContextSave (k) ========")
	khBytes, err := tpm2.ContextSave(rwc, kh)
	if err != nil {
		fmt.Printf("ContextSave failed for ekh1: %v", err)
	}

	//fmt.Println("khBytes ? : ", khBytes)

	err = ioutil.WriteFile(*keyFile, khBytes, 0644)
	if err != nil {
		fmt.Printf("ContextSave failed for ekh2: %v", err)
	} else {
		fmt.Println("succesfully saved into key.bin  ")
	}
	defer tpm2.FlushContext(rwc, kh)

	// Creating a public pem file
	kPublicKey, _, _, err := tpm2.ReadPublic(rwc, kh)
	if err != nil {
		log.Fatalf("Error tpmEkPub.Key() failed: %s", err)
	}

	ap, err := kPublicKey.Key()
	if err != nil {
		log.Fatalf("reading Key() failed: %s", err)
	}
	akBytes, err := x509.MarshalPKIXPublicKey(ap)
	if err != nil {
		log.Fatalf("Unable to convert ekpub: %v", err)
	}

	rakPubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: akBytes,
		},
	)
	log.Printf("     PublicKey: \n%v", string(rakPubPEM))

	err = ioutil.WriteFile(*publicKeyFile, rakPubPEM, 0644)
	if err != nil {
		log.Fatalf("Could not write file %v", err)
	}
	log.Printf("Public Key written to: %s", *publicKeyFile)

	// this bit converts the key into JWK format for some reason??
	der, err := x509.MarshalPKIXPublicKey(ap)
	if err != nil {
		log.Fatalf("keycreate: error converting public key: %v", err)
	}
	hasher := sha256.New()
	hasher.Write(der)
	kid := base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))

	jkey, err := jwk.New(ap)
	if err != nil {
		log.Fatalf("failed to create symmetric key: %s\n", err)
	}
	jkey.Set(jwk.KeyIDKey, kid)

	buf, err := json.MarshalIndent(jkey, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal key into JSON: %s\n", err)
		return
	}
	fmt.Printf("JWK Format:\n%s\n", buf)

}
