# Go Key Mime Pi

Inspired by the Python project [key-mime-py](https://github.com/mtlynch/key-mime-pi)

I decided to create a hard fork because I don't need the advanced features of the replacement project, [TinyPilot](https://github.com/mtlynch/tinypilot), and because I want to be able to submit blocks of text using an API.

## Overview

Use your Raspberry Pi as a remote-controlled keyboard that accepts keystrokes either through a web browser or through an API.

## Compatibility

* Raspberry Pi 4
* Raspberry Pi Zero W (maybe, haven't tested this yet)

## Pre-requisites

* Raspberry Pi OS Stretch or later
* Go 1.17+

## Quick Start

To begin, enable USB gadget support on the Pi by running the following commands:

```bash
sudo ./scripts/enable-usb-hid
sudo reboot
```

When the Pi reboots, log back into it and build and run Go Key Mime Pi with the following make command:

```bash
make build run
```

The primary GUI for Go Key Mime Pi can be found at `http://${PI_ADDR}:8000/`

Swagger docs can be found at `http://${PI_ADDR}:8000/swagger/index.html`