#!/usr/bin/env bash
set -euo pipefail

echo "Installing deckctl..."

# Validate requirements
command -v systemctl >/dev/null || { echo "systemd required"; exit 1; }
command -v go >/dev/null || { echo "Go required"; exit 1; }

for file in deckctl.service config.yaml 50-deckctl.rules; do
    [[ -f "$file" ]] || { echo "$file missing"; exit 1; }
done

echo "Requesting sudo..."
sudo -v

echo "Stopping old service..."
sudo systemctl stop deckctl.service 2>/dev/null || true
sudo systemctl disable deckctl.service 2>/dev/null || true

echo "Building binary..."
go mod tidy
go build -o deckctl main.go

echo "Installing binary..."
sudo install -m 755 ./deckctl /usr/local/bin/deckctl

echo "Installing service..."
sudo install -m 644 ./deckctl.service /etc/systemd/system/deckctl.service

sudo sed -i "s/^User=.*/User=$USER/" /etc/systemd/system/deckctl.service
sudo sed -i "s/^Group=.*/Group=$USER/" /etc/systemd/system/deckctl.service
sudo sed -i "s|^WorkingDirectory=.*|WorkingDirectory=$HOME|" /etc/systemd/system/deckctl.service

echo "Installing udev rules..."
sudo install -m 644 ./50-deckctl.rules /etc/udev/rules.d/

if getent group uucp >/dev/null; then
    sudo usermod -aG uucp "$USER"
    echo "You must log out and back in for group changes to apply."
fi

sudo udevadm control --reload-rules
sudo udevadm trigger

echo "Installing config..."
install -d "$HOME/.config/deckctl/scripts"
install -m 644 ./config.yaml "$HOME/.config/deckctl/config.yaml"

echo "Enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable deckctl.service
sudo systemctl start deckctl.service

echo "Installation complete."
