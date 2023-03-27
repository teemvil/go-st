#!bin/bash
 
mkdir /tmp/myvtpm
modprobe tpm_vtpm_proxy
swtpm chardev --vtpm-proxy --tpmstate dir=/tmp/myvtpm --tpm2 --ctrl type=tcp,port=2322