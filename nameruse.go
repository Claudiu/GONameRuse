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
	GenerateCallBacks []Callback

	Callbacks struct {
		CheckCom cb
		CheckComAndNet cb
		CheckNone cb
		CheckLevenshtein cb
	}

	CheckHost bool
	CheckWhois bool
	Verbouse bool
}

type WhoIsServer struct {
	Addr string
	Calls int
}

type Callback struct {
	Params interface{}
	Function cb
}

func InitNameRuse() *NameRuse {
	nr := new(NameRuse)
	
	nr.CheckHost = true
	nr.CheckWhois = true
	nr.Verbouse = false
	nr.MaxLoopFails = 20

	nr.AddWhoIs("com.whois-servers.net:43")
	nr.AddWhoIs("whois.crsnic.net:43")

	nr.initCallbacks()

	nr.AddValidator(CheckNone, nil)

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

func (nr *NameRuse) AddValidator(fun cb, params interface{}) (*NameRuse) {
	that := Callback{Function: fun, Params: params}
	nr.GenerateCallBacks = append(nr.GenerateCallBacks, that)
	return nr
}

func (nr *NameRuse) AddLevensthein(source string, miniumum int) (*NameRuse) {
	nr.AddValidator(nr.Callbacks.CheckLevenshtein, LevensteinParam {Like: source, Minimum: miniumum});
	return nr
}

func (nr *NameRuse) Validate(word string) bool {
	for _, val := range nr.GenerateCallBacks {
		if(val.Function(nr, word, val.Params)) {
			continue
		} else {
			return false
		}				
	}

	return true	
}

func (nr *NameRuse) GenerateName(format string) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result := []rune(format)

	for index, letter := range result {
		if(letter == 'V') {
			result[index] = vowels[rand.Intn(len(vowels))]
		} else if (letter == 'C') {
			result[index] = consonants[rand.Intn(len(consonants))]
		} else if (letter == '*') {
			any := rand.Intn(10)
			
			if(any < 5) {
				result[index] = vowels[rand.Intn(len(vowels))]
			} else {
				result[index] = consonants[rand.Intn(len(consonants))]
			}
		}
	}

	if(nr.Verbouse) {
		log.Println(string(result) + " was generated. ")
	}

	return string(result)
}

func (nr *NameRuse) GenerateN(format string, n int) (*NameRuse) {
	nr.LoopFails = 0
	CheckedNames := []string{}
	for i := 0; i < n && (nr.LoopFails != nr.MaxLoopFails); i=i {
		result := nr.GenerateName(format)
		if(nr.IsRepeating(result, CheckedNames)) {
			if(nr.Verbouse) {
				log.Printf("Element is being repeated. [%d/%d]\n", nr.LoopFails, nr.MaxLoopFails)
			}

			nr.LoopFails = nr.LoopFails + 1
			continue
		} else {
			CheckedNames = append(CheckedNames, result) 
			if(nr.Validate(string(result))) {
				nr.Names = append(nr.Names, string(result))
				i = i + 1				
			}
		}
	}

	return nr
}

func (nr *NameRuse) GenerateNLike(format string, n int) (*NameRuse) {
	return nr.GenerateN(nr.Like(format), n)
}