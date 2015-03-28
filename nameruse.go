package nameruse

import (
	"log"
	"math/rand"
	"time"
)

var vowels 	   = []rune{'a', 'e', 'i', 'o', 'u'}
var consonants = []rune{'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'x', 'z'}


type NameRuse struct {
	Names []string
	LoopFails int
	MaxLoopFails int
	Servers []WhoIsServer
	GenerateCallBack cb

	Callbacks struct {
		CheckCom cb
		CheckComAndNet cb
		CheckNone cb
	}

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
	nr.MaxLoopFails = 5

	nr.AddWhoIs("com.whois-servers.net:43")
	nr.AddWhoIs("whois.crsnic.net:43")

	nr.initCallbacks()

	nr.GenerateCallBack = CheckComAndNet

	return nr;
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


func (nr *NameRuse) GenerateName(format string) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result 	   := []rune(format)

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
	for i := 0; i < n && nr.LoopFails != nr.MaxLoopFails; i++ {
		result := nr.GenerateName(format)
		if(nr.IsRepeating(result)) {
			nr.LoopFails = nr.LoopFails + 1
			nr.GenerateName(format)
		} else {
			if(nr.GenerateCallBack(nr, 	string(result))) {
				nr.Names = append(nr.Names, string(result))
			}
		}
	}

	return nr.Names
}

func (nr *NameRuse) GenerateNLike(format string, n int) []string {
	return nr.GenerateN(nr.Like(format), n)
}