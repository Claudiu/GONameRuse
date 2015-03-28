package nameruse 

type cb func (a *NameRuse, b string) bool

var CheckComAndNet cb = func (nameruse *NameRuse, a string) bool {
	if (nameruse.IsDomainTaken(a + ".com") || nameruse.IsDomainTaken(a + ".net")) {
		return false
	} else {
		return true
	}
}

var CheckCom cb = func (nameruse *NameRuse, a string) bool {
	if (nameruse.IsDomainTaken(a + ".com")) {
		return false
	} else {
		return true
	}
}

var CheckNone cb = func (nameruse *NameRuse, a string) bool {
	return true
}

func (nr *NameRuse) initCallbacks() {
	nr.Callbacks.CheckCom = CheckCom
	nr.Callbacks.CheckComAndNet = CheckComAndNet
	nr.Callbacks.CheckNone = CheckNone
}