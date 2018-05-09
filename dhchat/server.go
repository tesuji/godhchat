package dhchat

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"

	// pb "github.com/lzutao/godiffhellchat/simple"

	"github.com/lzutao/godiffhellchat/dhkx"
)

const (
	// BUFSIZE is default buffersize
	BUFSIZE = 4096
)

// ChatServerStart implements method to communication to the chat client.
func ChatServerStart(port int) {
	addr := strings.Join([]string{"localhost", strconv.Itoa(port)}, ":")
	// listen on all interfaces
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Socket listen port %d failed: %s", port, err.Error())
	}

	defer func() {
		listener.Close()
		log.Println("Listener closed")
	}()

	log.Printf("Listenning on port: %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleRequest(conn)
	}

}

// p := &pb.SimpleMessage {
// 	Mode: pb.SimpleMessage_KEY
// 	Data: ga.Bytes()
// }

// data, err := proto.Marshal(p)
// if err != nil {
// 	log.Fatalln("unmarshaling error:", err)
// }

func containsByte(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func perror(err error) *big.Int {
	log.Println(err)
	return nil
}

func keyExchangeServer(conn net.Conn) *big.Int {
	a, _ := dhkx.NewDHKey(0)
	ga := a.PublicKey()
	log.Println("ga:", ga)

	_, err := fmt.Fprintln(conn, ga.String())
	if err != nil {
		return perror(err)
	}

	var src string
	_, err = fmt.Fscanln(conn, &src)
	if err != nil {
		return perror(err)
	}

	src = strings.TrimSuffix(src, "\n")

	gb, ok := new(big.Int).SetString(src, 10)
	if !ok {
		log.Println("Cannot convert string to big.Int")
		return nil
	}
	log.Println("gb:", gb)

	key, _ := a.SharedSecretKey(gb)
	log.Println("key:", key)
	return key
}

// handleRequest handles incoming requests.
func handleRequest(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Println("Closing connection from", conn.RemoteAddr())
	}()
	log.Println("Connection from", conn.RemoteAddr())

	key := keyExchangeServer(conn)
	communicate(conn, key)
}
