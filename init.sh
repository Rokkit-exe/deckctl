#!/bin/bash

if ! command -v systemctl &> /dev/null; then
    echo "systemctl could not be found, please install systemd to use deckctl as a service."
    exit
fi

if [ ! -f "./deckctl.service" ]; then
    echo "deckctl.service file not found in the current directory. Please ensure it is present."
    exit
fi

if [ ! -f "./config.yaml" ]; then
    echo "config.yaml file not found in the current directory. Please ensure it is present."
    exit
fi

sudo systemctl stop deckctl.service

go build -o deckctl main.go

sudo cp ./deckctl.service /etc/systemd/system/deckctl.service

sudo cp ./deckctl /usr/local/bin/deckctl

mkdir -p "$HOME/.config/deckctl"
cp ./config.yaml "$HOME/.config/deckctl/config.yaml"

if [ -f "/etc/systemd/system/deckctl.service" ]; then
    sudo systemctl daemon-reload
    sudo systemctl enable deckctl.service
    sudo systemctl start deckctl.service
else
    echo "Failed to install deckctl service."
fi

