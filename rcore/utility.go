package rcore

import "math"

// --------------------------------------------------------------------
//
type Scheduling struct {
	//
	LastExecutedEvent int
	PlannedInterval   int

	//
	Running bool
}

// --------------------------------------------------------------------
//
func MakeScheduling(anInterval int) Scheduling {
	//
	out := Scheduling{PlannedInterval: anInterval}

	//
	out.Running = false

	//
	return out
}

// --------------------------------------------------------------------
//
func (s Scheduling) PlanNext(nowTime int) Scheduling {
	//
	out := s

	//
	out.Running = true
	out.LastExecutedEvent = nowTime

	//
	return out
}

// --------------------------------------------------------------------
//
func (s Scheduling) IsTimeToFire(nowTime int) bool {
	//
	return (nowTime >= s.LastExecutedEvent+s.PlannedInterval) && s.Running
}

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
