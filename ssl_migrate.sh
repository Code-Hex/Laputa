#!/usr/bin/env sh
openssl genrsa -out laputa.key 2048
openssl req -new -key laputa.key -sha256 -out laputa.csr
openssl x509 -in laputa.csr -days 3650 -req -signkey laputa.key -sha256 -out laputa.crt
