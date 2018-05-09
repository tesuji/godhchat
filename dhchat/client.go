package dhchat

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lzutao/godiffhellchat/aescrypt"
	"github.com/lzutao/godiffhellchat/dhkx"
)

func donothing(t interface{}) {
	return
}

func keyExchangeClient(r *bufio.Reader, w *bufio.Writer) *big.Int {
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

	a, _ := dhkx.NewDHKey(0)
	ga := a.PublicKey()
	log.Println("ga:", ga)

	if _, err := w.WriteString(ga.String() + "\n"); err != nil {
		return perror(err)
	}
	w.Flush()

	key, _ := a.SharedSecretKey(gb)
	log.Println("key:", key)
	return key
}

// ChatClientStart starts a connection to chat server
func ChatClientStart(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	var (
		r = bufio.NewReader(conn)
		w = bufio.NewWriter(conn)
	)

	key := keyExchangeClient(r, w)
	donothing(key)
}

func communicate(conn net.Conn) {
	var (
		r = bufio.NewReader(conn)
		w = bufio.NewWriter(conn)
	)

	key := keyExchangeClient(r, w)
	hashsum := sha256.Sum256(key.Bytes())
	cip := aescrypt.NewAESCipher(hashsum)

	chCon := make(chan string)
	chSer := make(chan string)
	go clientReadConsole(chCon)
	go clientReadServer(r, chSer)
IOLoop:
	for {
		select {
		case text := <-chCon:
			log.Print("Sent:", text)
			// fmt.Fprint(conn, text)
			cip.SendSocket(conn, text)
			fmt.Print("<<< ")
		case msg, ok := <-chSer:
			if !ok {
				break IOLoop
			} else {
				plain, _ := cip.Decrypt([]byte(msg))
				fmt.Print(">>> ", plain)
				fmt.Print("<<< ")
			}
		case <-time.After(30 * time.Second):
			// Do something when there is nothing read from stdin
			fmt.Println("You're too slow.")
			return
		}
	}
}

func clientReadConsole(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			close(ch)
			return
		}
		ch <- text
	}
}

func clientReadServer(reader *bufio.Reader, ch chan string) {
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			close(ch)
			return
		}
		ch <- text
	}
}
