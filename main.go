package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	PORT  = "5000"
	HOST  = ""
	CTYPE = "udp4"
)

type Client struct {
	Ip       string
	Port     string
	Addr     *net.UDPAddr
	Num      int
	Nickname string
	Count    uint
}

func handleRequest(lisen *net.UDPConn, cbuf []byte, addr *net.UDPAddr, n int, clientsChan map[string]*Client, addrAlone string) {
	mutex.Lock()
	defer mutex.Unlock()
	writeChar := []byte("$ ")
	red := color.New(color.FgRed).SprintFunc()
	bleu := color.New(color.FgBlue).SprintFunc()

	if clientsChan[addrAlone].Count == 1 {
		lisen.WriteToUDP([]byte("What's your nickname: "), addr)
	} else if clientsChan[addrAlone].Count == 2 {
		msg := fmt.Sprintf("\n[%s] is in the game ;)\n$ ", bleu(clientsChan[addrAlone].Nickname))

		for _, v := range clientsChan {
			/*Other chatters*/
			if v.Ip!=addrAlone {
				if len(v.Nickname)>0 {
					lisen.WriteToUDP([]byte(msg), v.Addr)
				}
			}else{
				lisen.WriteToUDP(writeChar, v.Addr)
			}
		}
	} else if clientsChan[addrAlone].Count > 2 {
		msg := fmt.Sprintf("\n[%s] %s \n$ ", red(clientsChan[addrAlone].Nickname), cleanUp(string(cbuf[0:n])))
		for _, v := range clientsChan {
			/*Other chatters*/
			if v.Ip!=addrAlone {
				if len(v.Nickname)>0 {
					lisen.WriteToUDP([]byte(msg), v.Addr)
				}
			}else{
				lisen.WriteToUDP(writeChar, v.Addr)
			}
		}
	}
}

func makeClients(clients map[string]*Client,clientsChan chan map[string]*Client,addrSplit []string,addr *net.UDPAddr,n int, cbuf []byte){
	mutex.Lock()
	if _,ok := clients[addrSplit[0]];!ok {
		clients[addrSplit[0]] = &Client{Ip: addrSplit[0], Port: addrSplit[1], Addr: addr, Nickname: "", Num: n, Count: 0}
	}
	clients[addrSplit[0]].Count += 1
	if clients[addrSplit[0]].Count == 2 {
		clients[addrSplit[0]].Nickname = cleanUp(string(cbuf[0:n]))
	}
	clientsChan <- clients
	mutex.Unlock()
}

var mutex sync.Mutex

func main() {
	fmt.Println("UDP server")
	udpAdd, err := net.ResolveUDPAddr(CTYPE, fmt.Sprintf("%s:%s", HOST, PORT))
	errorHand(err)

	lisen, err := net.ListenUDP(CTYPE, udpAdd)
	errorHand(err)
	defer lisen.Close()

	clients := make(map[string]*Client, 0)
	clientsChan := make(chan map[string]*Client, 0)
	cbuf := make([]byte, 1024)
	for {
		n, addr, err := lisen.ReadFromUDP(cbuf)
		errorHand(err)
		addrSplit := strings.Split(addr.String(), ":")
		go makeClients(clients,clientsChan,addrSplit,addr,n, cbuf)
		go handleRequest(lisen, cbuf, addr, n, <-clientsChan, addrSplit[0])
	}
}

func errorHand(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cleanUp(s string) string {
	return strings.Replace(s, "\n", "", -1)
}