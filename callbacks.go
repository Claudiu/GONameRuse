package nameruse 

type cb func (a *NameRuse, b string, params interface{}) bool

type LevensteinParam struct {
	Like string
	Minimum int
}

var CheckComAndNet cb = func (nameruse *NameRuse, a string, params interface{}) bool {
	if (nameruse.IsDomainTaken(a + ".com") || nameruse.IsDomainTaken(a + ".net")) {
		return false
	} else {
		return true
	}
}

var CheckCom cb = func (nameruse *NameRuse, a string, params interface{}) bool {
	if (nameruse.IsDomainTaken(a + ".com")) {
		return false
	} else {
		return true
	}
}

var CheckNone cb = func (nameruse *NameRuse, a string, params interface{}) bool {
	return true
}

var CheckLevenshtein cb = func (nameruse *NameRuse, a string, params interface{}) bool {
	if(Levenshtein(params.(LevensteinParam).Like, a) < params.(LevensteinParam).Minimum) {
		//fmt.Println(nameruse.Levenshtein("reddit", a))
		return true
	} else {
		return false
	}
}

func (nr *NameRuse) initCallbacks() {
	nr.Callbacks.CheckCom = CheckCom
	nr.Callbacks.CheckComAndNet = CheckComAndNet
	nr.Callbacks.CheckNone = CheckNone
	nr.Callbacks.CheckLevenshtein = CheckLevenshtein
}