package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func getIPCSocket() net.Conn {
	socket, err := net.Dial("unix", os.Getenv("GREETD_SOCK"))
	if err != nil {
		log.Fatalln(err)
	}
	return socket
}

func main() {
	cmd := flag.String("cmd", "/usr/bin/sh", "Command to run after login")
	flag.Parse()

	file, err := os.OpenFile("/var/log/neogreet.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		os.Truncate(file.Name(), 0)
		log.SetOutput(file)
	}

	config := NewConfig("/etc/neogreet/neogreet.yaml")
	greeter := NewGreeter(getIPCSocket(), *cmd)
	greeter.OnResetConnection = getIPCSocket
	ui := NewUI(config, greeter)
	ui.Draw()
}
