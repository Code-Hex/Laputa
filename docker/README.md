# Raspbian container for ansible

This container is created based on the [Code-Hex/docker-rpi-raspbian](https://github.com/Code-Hex/docker-rpi-raspbian).  
It is very heavy and its size is over 70 MB. So, When using this container please use it as a test of ansible files.

# How to setup
## First
You need to make a rsa key.

    ./rsa_gen.sh

It is saved under the name "docker_rsa" in the `$HOME/.ssh`.  
  
## Second
You need to `git submodule update --init` in the project and

    cd docker/pi && docker build -t raspbian .

If you have not yet run them commands.
  
## Third
You can create a docker image with the following command.

    docker build -t virtual-pi .

# Usage
You can run the container in the background by following command

    docker run --rm -itdp 2222:22 virtual-pi:latest

Please stop the container by executing this command.

    docker stop <CONTAINER ID>