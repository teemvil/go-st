# go-st
This project is an attempt at creating a trustworthy scalable IoT-framework. 

The architecture is largely based on the collaborative innovation project done with Metropolia and Nokia Bell Labs, which you can find here: github.com/teemvil/iot

The framework features have ben redone with Go. There is an added JWT-based security feature, that uses TPM2 to create a key that is the used to sign the JWT-token. After a succefull remote-attestation the JWT is unlocked with the TPM2-created key, and only after that is the device deemed trusted.
