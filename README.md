# deckctl
Minimal daemon and cli to receive and send data to the `Rokkit_deck`.
Control what your deck looks like and how it behaves with a simple yaml configuration file.
Execute commands on button presses or slider changes, and read deck events to trigger actions on your computer.

See [Rokkit_deck](https://github.com/Rokkit-exe/rokkit_deck) for more information about the deck itself.

### Support
Deckctl has only been tested on `Arch Linux`, but should work on any Linux distribution with `systemd`.
If you encounter any issues, please open an issue on the GitHub repository.

### Features
- Send UI configuration to the deck (button text/color, slider color/label, etc.)
- Read deck events (button presses, slider changes)
- Execute commands on button presses or slider changes

### Future features
- Support for more complex Event types (long press, double press, toggle, etc.)
- Plugins for common applications (Spotify, Discord, Obs-Studio, etc.)
- Flashing Icons/Images to the deck

### Requirements
- Go
- `Rokkit_deck` hardware with firmware

### Installation
```bash
# clone repo
git clone https://github.com/Rokkit-exe/deckctl.git
cd deckctl

# 1. Install dependencies and build the binary
# 2. Place service, config and binary in the right place
# 3. Enable and start the service
./init.sh
```

### Usage
```bash
# daemon modes
deckctl daemon stop
# or 
systemctl stop deckctl


deckctl daemon start
# or
systemctl start deckctl


deacktl daemon restart
# or
systemctl restart deckctl

# flash configuration to the deck
deckctl flash -f "~/.config/deckctl/config.yaml"
```

### configuration

The configuration live in `~/.config/deckctl/config.yaml` and is read by the daemon on startup and when the `flash` command is executed.
I recommend placing custom scripts in `~/.config/deckctl/scripts` and referencing them in the config file under the `actions` section.

See [example config](./config.yaml)

### About

This project is in early development and is not yet feature complete. 
The goal is to have a simple and flexible way to control the deck and execute commands on button presses or slider changes. 
If you have any suggestions or want to contribute, feel free to open an issue or a pull request on the GitHub repository.

### History

I had bought an Elgato Stream Deck but found it to be expensive and did not have stable support for Linux. 
I wanted to have a similar device that I could use on my Linux machine, and have native support for it.
So i created my own deck using an ESP32, WaveShare touch display and rotary encoders to recreate the Stream Deck plus.

### License
This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.


