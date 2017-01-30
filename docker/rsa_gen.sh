#!/usr/bin/env sh

ssh-keygen -t rsa -N "" -f $HOME/.ssh/docker_rsa
cp $HOME/.ssh/docker_rsa.pub authorized_keys

echo "Created $HOME/.ssh/docker_rsa $HOME/.ssh/docker_rsa.pub"
echo "Copied $HOME/.ssh/docker_rsa.pub as ./authorized_keys"