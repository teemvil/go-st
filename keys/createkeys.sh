#!/bin/bash
 
echo "creating (owner hierarchy) primary key"
tpm2_createprimary -C o -g sha256 -G rsa -c o.ctx
 
echo "creating (owner hierarchy) RSA key)"
tpm2_create -C o.ctx -u k1.pub -r k1.priv
tpm2_load -C o.ctx -u k1.pub -r k1.priv -c k1.ctx
tpm2_evictcontrol -C o -c k1.ctx 0x81010004
tpm2_readpublic -Q -c k1.ctx -f "pem" -o k1.pem
 
echo "creating ek key"
tpm2_createek -c 0x81010002 -G rsa -u ek.pub
tpm2_readpublic -c 0x81010002 -o ek.pem -f pem
 
echo "creating ak key"
tpm2_createak -C 0x81010002 -c ak.ctx -G rsa -g sha256 -s rsassa -u ak.pub -f pem -n ak.name
tpm2_evictcontrol -c ak.ctx 0x81010003
tpm2_readpublic -c 0x81010003 -o ak.pem -f pem