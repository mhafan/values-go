//
package rcore

// --------------------------------------------------------------------
// ...
import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// --------------------------------------------------------------------
//
const (
	// start of a new experiment
	CallStart = "START"
	// end of the current experiment
	CallEnd = "END"

	// Simulation loop calls
	// TCM -> CNT -> PUMP -> PATMOD -> SENSOR -> TCM
	CallCNT    = "CNT"
	CallPump   = "PUMP"
	CallPatMod = "PATMOD"
	CallSensor = "SENSOR"
	CallTCM    = "TCM"

	//
	DrugRocuronium    = "ROC"
	DrugCisatracurium = "CIS"

	//
	SexMale   = "male"
	SexFemale = "female"

	MasterChannel = "vm.master"
)

// --------------------------------------------------------------------
// to be mirrored at REDIS
type Exprec struct {
	// ----------------------------------------------------------------
	// e.g. vm.someTestCase.123
	// set by TCM
	Varname string

	// ----------------------------------------------------------------
	// Experiment definition: (set by TCM)
	// MUST: DrugRocuronium | DrugCisatracurium
	Drug string
	// MUST: [kg], to compute Vd
	Weight int
	// for future use
	Age int
	// for future use
	TargetTOF int
	TargetPTC int
	// new added: regime of infusion
	// { none, rboluses (default) }
	CNTStrategy string

	// ----------------------------------------------------------------
	// Decisions made by CNT (CNT can be disabled and values set from TCM)
	// Values can be slightly updated by PUMP (noise/fault injection)
	// [mL]
	Bolus int
	// [mL/hr]
	Infusion int

	// ----------------------------------------------------------------
	// Parameters essential for PatientModel (PM).
	// Might be set by TCM.
	// Typically left in default values.
	// Vd is either Weight*UnitVd or AbsoluteVd
	// Vd => concentration [ug/mL]
	UnitVd     int
	AbsoluteVd int
	EC50       Double

	// ----------------------------------------------------------------
	// Outputs from PM
	// newly added: Cinp
	// concentration in blood plasma (C0)
	// [ug/mL]
	Cinp Double
	// estimated TOF and PTC
	TOF int
	PTC int
	// newly added
	// cumulative consumption of drug by patient in [mL] of solution
	ConsumedTotal Double
	// estimated time till full recovery [s]
	RecoveryTime int

	// ----------------------------------------------------------------
	// Controlled by TCM.
	// Cycle += 1 after every simulation cycle
	// Mtime += step [s]
	// ----------------------------------------------------------------
	// PM works in [s], 1 second time granularity
	Mtime int
	Cycle int
}

// --------------------------------------------------------------------
//
var CurrentExp *Exprec

// --------------------------------------------------------------------
//
func NewExpID(testcase string) string {
	//
	n, _ := redis.Int(Global.handleOUT.Do("INCR", "vm.expcounter"))

	//
	return "vm." + testcase + "." + strconv.Itoa(n)
}

// --------------------------------------------------------------------
//
func MakeExpID(vn string) *Exprec {
	//
	return &Exprec{Varname: vn}
}

// --------------------------------------------------------------------
//
func (r *Exprec) channel() string {
	//
	return r.Varname
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

// --------------------------------------------------------------------
//
func (r *Exprec) Save_i(key string, i int, keys []string, all bool) {
	//
	if all == false {
		//
		if Contains(keys, key) == false {
			return
		}
	}

	//
	Global.handleOUT.Do("hset", r.Varname, key, i)
}

// --------------------------------------------------------------------
//
func (r *Exprec) Save_s(key string, i string, keys []string, all bool) {
	//
	if all == false {
		//
		if Contains(keys, key) == false {
			return
		}
	}

	//
	Global.handleOUT.Do("hset", r.Varname, key, i)
}

// --------------------------------------------------------------------
//
func (r *Exprec) Save_d(key string, i Double, keys []string, all bool) {
	//
	if all == false {
		//
		if Contains(keys, key) == false {
			return
		}
	}

	//
	Global.handleOUT.Do("hset", r.Varname, key, i)
}

// --------------------------------------------------------------------
//
func (r *Exprec) Save(keys []string, all bool) bool {
	//
	r.Save_s("drug", r.Drug, keys, all)
	r.Save_d("EC50", r.EC50, keys, all)

	//
	r.Save_i("bolus", r.Bolus, keys, all)
	r.Save_i("infusion", r.Infusion, keys, all)

	//
	r.Save_i("weight", r.Weight, keys, all)
	r.Save_i("age", r.Age, keys, all)

	r.Save_i("mtime", r.Mtime, keys, all)
	r.Save_i("cycle", r.Cycle, keys, all)

	r.Save_i("unitVd", r.UnitVd, keys, all)
	r.Save_i("absoluteVd", r.AbsoluteVd, keys, all)
	r.Save_i("targetTOF", r.TargetTOF, keys, all)
	r.Save_i("targetPTC", r.TargetPTC, keys, all)
	r.Save_i("TOF", r.TOF, keys, all)
	r.Save_i("PTC", r.PTC, keys, all)

	//
	return true
}

// --------------------------------------------------------------------
//
func (r *Exprec) Load(keys []string, all bool) bool {
	//
	var _s = Global.handleOUT

	f := func(k string) bool { return all == true || Contains(keys, k) }

	//
	if f("drug") {
		r.Drug, _ = redis.String(_s.Do("hget", r.Varname, "drug"))
	}

	if f("bolus") {
		r.Bolus, _ = redis.Int(_s.Do("hget", r.Varname, "bolus"))
	}

	if f("infusion") {
		r.Infusion, _ = redis.Int(_s.Do("hget", r.Varname, "infusion"))
	}

	//
	if f("weight") {
		r.Weight, _ = redis.Int(_s.Do("hget", r.Varname, "weight"))
	}

	if f("age") {
		r.Age, _ = redis.Int(_s.Do("hget", r.Varname, "age"))
	}

	if f("mtime") {
		r.Mtime, _ = redis.Int(_s.Do("hget", r.Varname, "mtime"))
	}

	if f("cycle") {
		r.Cycle, _ = redis.Int(_s.Do("hget", r.Varname, "cycle"))
	}

	if f("unitVd") {
		r.UnitVd, _ = redis.Int(_s.Do("hget", r.Varname, "unitVd"))
	}

	if f("absoluteVd") {
		r.AbsoluteVd, _ = redis.Int(_s.Do("hget", r.Varname, "absoluteVd"))
	}

	if f("targetTOF") {
		r.TargetTOF, _ = redis.Int(_s.Do("hget", r.Varname, "targetTOF"))
	}

	if f("targetPTC") {
		r.TargetPTC, _ = redis.Int(_s.Do("hget", r.Varname, "targetPTC"))
	}

	if f("EC50") {
		r.EC50, _ = redis.Float64(_s.Do("hget", r.Varname, "EC50"))
	}
	if f("TOF") {
		r.TOF, _ = redis.Int(_s.Do("hget", r.Varname, "TOF"))
	}
	if f("PTC") {
		r.PTC, _ = redis.Int(_s.Do("hget", r.Varname, "PTC"))
	}

	//
	return true
}

// --------------------------------------------------------------------
//
func (r *Exprec) Say(message string) {
	//
	Global.topublish <- rmsg{r.channel(), message}
}
