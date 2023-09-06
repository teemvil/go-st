package keyget

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
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
