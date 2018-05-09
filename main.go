package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lzutao/godhchat/dhchat"
)

func init() {
	// disable all log
	log.SetOutput(ioutil.Discard)
}

func usage() {
	fmt.Printf("usage: %s [-h] [--port PORT] [--listen | --host HOST]\n\n", os.Args[0])
	fmt.Printf("Diffie-Hellman Chat - encrypted chat with Diffie-Hellman key exchange\n\n")
	fmt.Println("optional arguments:")
	flag.PrintDefaults()
	fmt.Println("\n[+] Written by 15520599")
}

func main() {
	var isListen bool
	var portNumber uint
	var hostIP string

	flag.UintVar(&portNumber, "port", 6001, "port number for connection")
	flag.BoolVar(&isListen, "listen", false, "listen to a port (Chat Server)")
	flag.StringVar(&hostIP, "host", "", "IP address to connect to")

	flag.Usage = usage
	flag.Parse()

	if portNumber == 0 {
		portNumber = 6001
	}

	if isListen {
		dhchat.ChatServerStart(int(portNumber))
	} else if hostIP != "" {
		dhchat.ChatClientStart(hostIP, int(portNumber))
	} else {
		fmt.Println("Choose --listen or --host to continue.")
		flag.PrintDefaults()
	}

}
