#!/usr/bin/env sh

mkdir -p ./ssl
cd ./ssl
openssl genrsa -out laputa.key 2048
openssl req -new -key laputa.key -sha256 -out laputa.csr
openssl x509 -in laputa.csr -days 36500 -req -signkey laputa.key -sha256 -out laputa.crt
