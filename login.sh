#!/bin/bash
export PATH=$PATH:/home/user/bin
export GOCACHE=$(pwd)/cache
export GOROOT=/usr/lib/go #gosetup
export GOPATH=/home/$USER/go #gosetup
export DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$(id -u)/bus
current_dir=$(pwd)
cd "$current_dir"/ && go build main.go && ./main -dir /home/"$USER"/Documents/Personal/wallpaper/wallpapers/
