package rcore

import "math"

// --------------------------------------------------------------------
//
type MStat struct {
	//
	N int

	Min Double
	Max Double
}

//
func (ms *MStat) Update(inval Double) {
	//
	if ms.N <= 0 {
		//
		ms.Max = inval
		ms.Min = inval
	} else {
		//
		ms.Max = math.Max(ms.Max, inval)
		ms.Min = math.Min(ms.Min, inval)
	}

	//
	ms.N++
}

//
func (ms *MStat) Reset() {
	//
	ms.N = 0
	ms.Max = 0
	ms.Min = 0
}

// --------------------------------------------------------------------
//
func Contains(keys []string, key string) bool {
	//
	for _, v := range keys {
		//
		if v == key {
			return true
		}
	}

	//
	return false
}
