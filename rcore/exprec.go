//
package rcore

// --------------------------------------------------------------------
// ...
import (
	"fmt"
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
	CNTStratBasic = "basic"
	CNTStratFWSim = "fwsim"
	CNTStratNone  = "none"

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
	TargetCinpLow Double
	TargetCinpHi  Double
	// new added: regime of infusion
	// { none, rboluses (default) }
	CNTStrategy string
	RepeStep    int
	RepeBolus   int
	FwRange     int
	IbolusMg    Double

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
	// ----------------------------------------------------------------
	// these should ramain constant/intact during the experiment
	UnitVd     int
	AbsoluteVd int
	EC50       Double

	// ----------------------------------------------------------------
	// Outputs from PM
	// newly added: Cinp
	// concentration in blood plasma (C0)
	// [ug/mL]
	Cinp Double

	// ----------------------------------------------------------------
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

	// ----------------------------------------------------------------
	// Outputs from TOF/PTC sensor
	// 4x values of TOF, PTC value
	// Status:
	// == 0, no update
	// == 1 (bit 0), TOF updated
	// == 2 (bit 1), PTC updated
	SensorTOF0    int
	SensorTOF1    int
	SensorTOF2    int
	SensorTOF3    int
	SensorPTC     int
	SensorStatus  int
	SensorCommand int

	// ----------------------------------------------------------------
	// To be added
	// - weight coefficient for the initial bolus
	// <0, xx>
	Wcoef Double
}

// --------------------------------------------------------------------
//
var CurrentExp *Exprec

// --------------------------------------------------------------------
// New experiment ID in the form "vm." + testcase + "." + num
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
func CsvHeader() string {
	//
	return "Cycle,Mtime,Bolus,Infusion,Cinp,TOF,PTC,Consumed,RecoveryTime\n"
}

// --------------------------------------------------------------------
//
func (r *Exprec) CsvExport() string {
	//
	out := fmt.Sprintf("%d,%d,%d,%d,%.2f,%d,%d,%.2f,%d", r.Cycle, r.Mtime, r.Bolus, r.Infusion, r.Cinp, r.TOF, r.PTC, r.ConsumedTotal, r.RecoveryTime)

	//
	return out
}

// --------------------------------------------------------------------
//
func (r *Exprec) channel() string {
	//
	return r.Varname
}

// --------------------------------------------------------------------
//
func (r *Exprec) ischannel(nm string) bool {
	//
	return r.channel() == nm
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
	r.Save_s("CNTStrategy", r.CNTStrategy, keys, all)

	r.Save_d("EC50", r.EC50, keys, all)
	r.Save_d("Cinp", r.Cinp, keys, all)
	r.Save_d("ConsumedTotal", r.ConsumedTotal, keys, all)
	r.Save_d("targetCinpHi", r.TargetCinpHi, keys, all)
	r.Save_d("targetCinpLow", r.TargetCinpLow, keys, all)
	r.Save_d("wcoef", r.Wcoef, keys, all)
	r.Save_d("ibolusMg", r.IbolusMg, keys, all)

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
	r.Save_i("TOF", r.TOF, keys, all)
	r.Save_i("PTC", r.PTC, keys, all)
	r.Save_i("RecoveryTime", r.RecoveryTime, keys, all)

	//
	r.Save_i("repeStep", r.RepeStep, keys, all)
	r.Save_i("repeBolus", r.RepeBolus, keys, all)
	r.Save_i("fwRange", r.FwRange, keys, all)

	//
	r.Save_i("SensorTOF0", r.SensorTOF0, keys, all)
	r.Save_i("SensorTOF1", r.SensorTOF1, keys, all)
	r.Save_i("SensorTOF2", r.SensorTOF2, keys, all)
	r.Save_i("SensorTOF3", r.SensorTOF3, keys, all)
	r.Save_i("SensorPTC", r.SensorPTC, keys, all)
	r.Save_i("SensorStatus", r.SensorStatus, keys, all)
	r.Save_i("SensorCommand", r.SensorCommand, keys, all)

	//
	return true
}

// --------------------------------------------------------------------
//
func (r *Exprec) Load_i(key string, i *int, keys []string, all bool) {
	//
	if all || Contains(keys, key) {
		//
		*i, _ = redis.Int(Global.handleOUT.Do("hget", r.Varname, key))
	}
}

// --------------------------------------------------------------------
//
func (r *Exprec) Load_d(key string, i *Double, keys []string, all bool) {
	//
	if all || Contains(keys, key) {
		//
		*i, _ = redis.Float64(Global.handleOUT.Do("hget", r.Varname, key))
	}
}

// --------------------------------------------------------------------
//
func (r *Exprec) Load_s(key string, i *string, keys []string, all bool) {
	//
	if all || Contains(keys, key) {
		//
		*i, _ = redis.String(Global.handleOUT.Do("hget", r.Varname, key))
	}
}

// --------------------------------------------------------------------
//
func (r *Exprec) LoadAll() {
	//
	r.Load([]string{}, true)
}

// --------------------------------------------------------------------
//
func (r *Exprec) Load(keys []string, all bool) bool {
	//
	r.Load_s("drug", &r.Drug, keys, all)
	r.Load_s("CNTStrategy", &r.CNTStrategy, keys, all)

	//
	r.Load_i("bolus", &r.Bolus, keys, all)
	r.Load_i("infusion", &r.Infusion, keys, all)
	r.Load_i("weight", &r.Weight, keys, all)
	r.Load_i("age", &r.Age, keys, all)
	r.Load_i("mtime", &r.Mtime, keys, all)
	r.Load_i("cycle", &r.Cycle, keys, all)
	r.Load_i("unitVd", &r.UnitVd, keys, all)
	r.Load_i("absoluteVd", &r.AbsoluteVd, keys, all)
	r.Load_i("TOF", &r.TOF, keys, all)
	r.Load_i("PTC", &r.PTC, keys, all)
	r.Load_i("RecoveryTime", &r.RecoveryTime, keys, all)

	r.Load_i("repeStep", &r.RepeStep, keys, all)
	r.Load_i("repeBolus", &r.RepeBolus, keys, all)
	r.Load_i("fwRange", &r.FwRange, keys, all)

	//
	r.Load_d("EC50", &r.EC50, keys, all)
	r.Load_d("Cinp", &r.Cinp, keys, all)
	r.Load_d("ConsumedTotal", &r.ConsumedTotal, keys, all)
	r.Load_d("targetCinpHi", &r.TargetCinpHi, keys, all)
	r.Load_d("targetCinpLow", &r.TargetCinpLow, keys, all)
	r.Load_d("wcoef", &r.Wcoef, keys, all)
	r.Load_d("ibolusMg", &r.IbolusMg, keys, all)

	//
	r.Load_i("SensorTOF0", &r.SensorTOF0, keys, all)
	r.Load_i("SensorTOF1", &r.SensorTOF1, keys, all)
	r.Load_i("SensorTOF2", &r.SensorTOF2, keys, all)
	r.Load_i("SensorTOF3", &r.SensorTOF3, keys, all)
	r.Load_i("SensorPTC", &r.SensorPTC, keys, all)
	r.Load_i("SensorStatus", &r.SensorStatus, keys, all)
	r.Load_i("SensorCommand", &r.SensorCommand, keys, all)

	//
	return true
}

// --------------------------------------------------------------------
//
func (r *Exprec) Say(message string) {
	//
	Global.topublish <- Rmsg{r.channel(), message}
}
