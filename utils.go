package nameruse

import (
	"net"
	"log"
)

func isVowel(letter rune) bool {
	for _, val := range vowels {
		if(val == letter) {
			return true
		}
	}

	return false
}

func (nr *NameRuse) Clear() {
	nr.Names = nr.Names[0:0]
}

func (nr *NameRuse) Like(word string) string {
	result := make([]rune, len(word))

	for index, letter := range word {
		if(isVowel(letter)) {
			result[index] = 'V'
		} else {
			result[index] = 'C'
		}
	}

	return string(result)
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

