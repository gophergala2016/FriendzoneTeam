#!/bin/bash
sudo apt-get -y install gccgo-4.9
wget https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz -O /path/
sudo tar -C /usr/local -xzf go1.5.3.linux-amd64.tar.gz
echo 'export GOPATH=$HOME/work' >> $HOME/.profile
echo 'export GOROOT=/usr/local/go' >> $HOME/.profile
echo 'export export PATH=$PATH:$GOROOT/bin' >> $HOME/.profile
source .profile
