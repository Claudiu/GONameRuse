package nameruse

import (
	"log"
	"math/rand"
	"time"
	"bufio"
	"strings"
	"net"
	"fmt"
)

type cb func(string) bool

type NameRuse struct {
	Names []string
	LoopFails int
	Servers []WhoIsServer
	GenerateCallBack cb
	CheckHost bool
	CheckWhois bool
	Verbouse bool
}

type WhoIsServer struct {
	Addr string
	Calls int
}

func InitNameRuse() *NameRuse {
	nr := new(NameRuse)
	
	nr.CheckHost = true
	nr.CheckWhois = true
	nr.Verbouse = false

	nr.AddWhoIs("com.whois-servers.net:43")
	nr.AddWhoIs("whois.crsnic.net:43")

	nr.GenerateCallBack = func (a string) bool {
		if(nr.Verbouse) {
			log.Println("Default callback fired with " + a + ".")
		}

		return true;
	}

	return nr;
}

func (nr *NameRuse) AddWhoIs(addr string) {
	if(nr.Verbouse) {
		log.Println("Added WhoIs Server at " + addr + ". ")
	}
	
	nr.Servers = append(nr.Servers, WhoIsServer {Addr: addr, Calls: 0})
}


func (nr *NameRuse) GetWhoisServer() string {
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

func (nr *NameRuse) Hipsterize(word string) (string, bool) {
	//TODO: Add more TLDs
	tld := []string{"ru", "porn" , "ro", "io", "ie", "eo", "iu", "ae"}

	for _, domain := range tld {
		if(string([]rune(word)[len(word) - len(domain):]) == domain) {
			return string([]rune(word)[:len(word) - len(domain)]) + "." + domain, true
		}
	}

	if(nr.Verbouse) {
		log.Println("Failed to hipsterize domain with name " + word)
	}

	return string(word), false
}

func (nr *NameRuse) IsFreeByHost(word string) bool {
	val, err := net.LookupHost(word)
	if(err == nil && len(val) == 1) {
		if(nr.Verbouse) {
			log.Println(word + " is taken by Host with IP: "  + val[0])
		}
		return false;
	} else {
		if(nr.Verbouse) {
			log.Println(word + " is free by Host check. ")
		}

		return true;
	}	
}

func (nr *NameRuse) IsDomainTaken(word string) bool {
	if(nr.CheckHost == false && nr.CheckWhois == false) {
		panic("CheckWhois and CheckHost cannot be both false.")
	}

	if(!nr.IsFreeByHost(word) && nr.CheckHost) {
		return true;
	}

	if(!nr.IsFreeByWhois(word) && nr.CheckWhois) {
		return true;
	}

	return false;
}

func (nr *NameRuse) IsRepeating(word string) bool {
	for _, name := range nr.Names {
		if name == word {
			if(nr.Verbouse) {
				log.Println(word + " is already in list. ")
			}
			return true;
		}
	}

	return false
}

func (nr *NameRuse) GenerateName(format string) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result 	   := []rune(format)
	vowels 	   := []rune{'a', 'e', 'i', 'o', 'u'}
	consonants := []rune{'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'x', 'z'}

	for index, letter := range result {
		if(letter == 'V') {
			result[index] = vowels[rand.Intn(len(vowels))]
		} else if (letter == 'C') {
			result[index] = consonants[rand.Intn(len(consonants))]
		}
	}

	if(nr.Verbouse) {
		log.Println(string(result) + " was generated. ")
	}

	return string(result)
}

func (nr *NameRuse) GenerateN(format string, n int) []string {
	nr.LoopFails = 0
	for i := 0; i < n && nr.LoopFails != 3; i++ {
		result := nr.GenerateName(format)
		if(nr.IsRepeating(result)) {
			nr.LoopFails = nr.LoopFails + 1
			nr.GenerateName(format)
		} else {
			if(nr.GenerateCallBack(string(result))) {
				nr.Names = append(nr.Names, string(result))
			}
		}
	}

	return nr.Names
}