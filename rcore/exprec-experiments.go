//
package rcore

// --------------------------------------------------------------------
// ...
import (
	"fmt"
	"os"
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

	//
	MasterChannel = "vm.master"

	//
	CuffCommandTOF = 1
	CuffCommandPTC = 2
)

// --------------------------------------------------------------------
// New experiment ID in the form "vm." + testcase + "." + num
func NewExpID(testcase string) string {
	//
	n, err := redis.Int(Global.handleOUT.Do("INCR", "vm.expcounter"))

	//
	if err != nil {
		//
		os.Exit(1)
	}

	//
	return "vm." + testcase + "." + strconv.Itoa(n)
}

// --------------------------------------------------------------------
//
func CsvHeader() string {
	//
	return "Cycle,Mtime,Bolus,Infusion,Cinp,TOF,PTC,Consumed,RecoveryTime,SensorCommand,SensorStatus\n"
}

// --------------------------------------------------------------------
//
func (r *Exprec) CsvExport() string {
	//
	out := fmt.Sprintf("%d,%d,%.2f,%.2f,%.2f,%d,%d,%.2f,%d,%d,%d", r.Cycle, r.Mtime, r.Bolus, r.Infusion, r.Cinp, r.TOF, r.PTC, r.ConsumedTotal, r.RecoveryTime, r.SensorCommand, r.SensorStatus)

	//
	return out
}
