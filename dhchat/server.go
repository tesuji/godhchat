package dhchat

import (
	"bufio"
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

func perror(err error) *big.Int {
	log.Println(err)
	return nil
}

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

func keyExchangeServer(r *bufio.Reader, w *bufio.Writer) *big.Int {
	a, _ := dhkx.NewDHKey(0)
	ga := a.PublicKey()
	log.Println("ga:", ga)

	if _, err := w.WriteString(ga.String() + "\n"); err != nil {
		return perror(err)
	}
	w.Flush()

	src, err := r.ReadString('\n')
	if err != nil {
		return perror(err)
	}

	src = src[:len(src)-1]
	gb, ok := new(big.Int).SetString(src, 10)
	if !ok {
		log.Println("Not good")
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

	var (
		r = bufio.NewReader(conn)
		w = bufio.NewWriter(conn)
	)

	key := keyExchangeServer(r, w)
	donothing(key)
}
