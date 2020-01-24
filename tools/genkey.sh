openssl genrsa 4096 > ../logger/keys/prikey.pem
openssl rsa -pubout < ../logger/keys/prikey.pem > ../logger/keys/pubkey.key
