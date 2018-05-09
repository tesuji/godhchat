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

	"github.com/lzutao/godhchat/aescrypt"
	"github.com/lzutao/godhchat/dhkx"
)

// keyExchangeClient exchanges Diffie-Hellman key with keyExchangeServer
func keyExchangeClient(conn net.Conn) *big.Int {
	var src string
	_, err := fmt.Fscanln(conn, &src)
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

	a, _ := dhkx.NewDHKey(0)
	ga := a.PublicKey()
	log.Println("ga:", ga)

	_, err = fmt.Fprintln(conn, ga.String())
	if err != nil {
		return perror(err)
	}

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

	key := keyExchangeClient(conn)
	communicate(conn, key)
}

// communicate starts communications between client and server
// after recv shared secret key
func communicate(conn net.Conn, key *big.Int) {
	hashsum := sha256.Sum256(key.Bytes())
	cip := aescrypt.NewAESCipher(hashsum)

	chCon := make(chan string)
	chRmt := make(chan string)
	go readConsole(chCon)
	go readRemote(conn, chRmt)
IOLoop:
	for {
		select {
		case text := <-chCon:
			log.Print("Sent:", text)
			_, err := cip.SendSocket(conn, text)

			if err != nil {
				perror(err)
			}

			fmt.Print("<<< ")
		case msg, ok := <-chRmt:
			if !ok {
				break IOLoop
			} else {
				s, err := cip.Decrypt(msg)
				if err != nil {
					perror(err)
				}
				fmt.Print(">>> ", s)
				fmt.Print("<<< ")
			}
		case <-time.After(30 * time.Second):
			fmt.Println("You're too slow.")
			return
		}
	}
}

// readConsole reads lines from keyboard and return them
func readConsole(ch chan string) {
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

// readRemote reads lines other remote party and return them
func readRemote(conn net.Conn, ch chan string) {
	for {
		var text string
		_, err := fmt.Fscanln(conn, &text)
		if err != nil {
			close(ch)
			return
		}
		ch <- text
	}
}
