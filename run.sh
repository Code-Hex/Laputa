#!/usr/bin/env sh

export LAPUTA_CERTFILE="./ssl/laputa.crt"
export LAPUTA_KEYFILE="./ssl/laputa.key"
export LAPUTA_AKATSUKI="https://127.0.0.1:3000/"
export LAPUTA_MODE="develop"
export LAPUTA_FLOOR="F321"
export LAPUTA_PORT="4000"
export LAPUTA_DEBUG=true
./laputa
