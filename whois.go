package nameruse

import (
	"bufio"
	"strings"
	"net"
	"log"
	"fmt"
	"os"
)

func (nr *NameRuse) AddWhoIs(addr string) {
	//TODO: Check if host is alive
	if(nr.Verbouse) {
		log.Println("Added WhoIs Server at " + addr + ". ")
	}
	
	nr.Servers = append(nr.Servers, WhoIsServer {Addr: addr, Calls: 0})
}

func (nr *NameRuse) AddWhoIsFile(path string) {
	 file, err := os.Open(path)

	 if (err != nil) {
	 	panic(err)
	 }

	 defer file.Close()
	 scanner := bufio.NewScanner(file)

	 for scanner.Scan() {
	 	nr.AddWhoIs(scanner.Text())
	 }
}

func (nr *NameRuse) GetWhoisServer() string {
	//TODO: Implement Chaining rather than iterating over again
	smallestCall := nr.Servers[0].Calls
	smallestIndex := 0

	for index, server := range nr.Servers {
		if(server.Calls < smallestCall) {
			smallestCall = server.Calls
			smallestIndex = index
		}
	}

	nr.Servers[smallestIndex].Calls = nr.Servers[smallestIndex].Calls + 1
	return nr.Servers[smallestIndex].Addr
}

func (nr *NameRuse) IsFreeByWhois(domain string) bool {
	//FIXME: Timeout
	whois, err := net.Dial("tcp", nr.GetWhoisServer())

	if err != nil {
		fmt.Println(err)
	}
	defer whois.Close()

	whois.Write([]byte(domain + "\r\n"))

	whoisbuf := bufio.NewReader(whois)

	for {
		str, err := whoisbuf.ReadString('\n')
		
		if len(str)>0 {
			if(strings.Contains(str, "Expiration Date")) {
				if(nr.Verbouse) {
					log.Println(domain + " is taken by WhoIs check. ")
					log.Println(str[3:len(str) - 1])
				}
				return false;
			}
		}

		if err!= nil {
			break
		}
	}

	if(nr.Verbouse) {
		log.Println(domain + " is free by WhoIs check. ")
	}
	return true
} 