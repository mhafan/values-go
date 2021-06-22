// ----------------------------------------------------------------------
// PatMod for R-system
// ----------------------------------------------------------------------
package main

//
import (
	"fmt"
	"log"
	"rcore"
)

// ----------------------------------------------------------------------
//
type PMEntity struct {
	//
	rsims *rcore.SIMS
}

// ----------------------------------------------------------------------
//
func (e *PMEntity) MyTurn() string {
	//
	return rcore.CallPatMod
}

// ----------------------------------------------------------------------
//
func (e *PMEntity) ResetState() {
	//
	e.rsims = nil
}

// NONE
func (e *PMEntity) DefaultFunction(msg rcore.Rmsg) {
	// NONE
}

// ----------------------------------------------------------------------
//
func pmTotalRecoveryPredicate(sims *rcore.SIMS) bool {
	//
	return sims.Effect.TOFSimpleAmplitude() < 95
}

// ----------------------------------------------------------------------
// For every cycle of distributed simulation.
func (e *PMEntity) CycleFunction() {
	// dont forget to pass the token
	defer rcore.CurrentExp.Say(rcore.CallSensor)

	//
	tol := []string{"mtime", "cycle", "bolus", "infusion", "SensorStatus"}

	// --------------------------------------------------------------------
	// load Redis record
	if rcore.EntityExpRecReload(tol) == false {
		//
		panic("REDIS record not found")
	}

	// ...
	var _c = rcore.CurrentExp

	//
	if _c.SensorStatus > 0 {
		//
		fmt.Println("PM: Sensor correction")
	}

	// --------------------------------------------------------------------
	// REDIS record -> SIMS (patmod simulation state)
	e.rsims.UpdateFrom(_c)

	//
	e.rsims = e.rsims.SimSteps(_c.Mtime)

	//
	trec := e.rsims.Clone().SimStepsWhile(pmTotalRecoveryPredicate)

	// --------------------------------------------------------------------
	//
	_c.TOF = e.rsims.Effect.TOFSimpleAmplitude()
	_c.PTC = 0
	_c.Cinp = e.rsims.Cinp()
	_c.ConsumedTotal += e.rsims.BolusConsumptionML
	_c.RecoveryTime = (trec.Time - e.rsims.Time)

	//
	log.Println("PMA; cycle: ", _c.Cycle, "mtime", _c.Mtime,
		"rtime", e.rsims.Time, "bolus ", e.rsims.Bolus.Value, " ", e.rsims.Bolus, "ConsT ", _c.ConsumedTotal)

	//
	rcore.CurrentExp.Save([]string{"TOF", "PTC", "Cinp", "ConsumedTotal", "RecoveryTime"}, false)
}

// ----------------------------------------------------------------------
// vm.X.Y -> START
// --- Initialization of new experiment
func (e *PMEntity) StartFunction() {
	//
	e.rsims = rcore.EmptySIMS()

	//
	e.rsims.SetupFrom(rcore.CurrentExp)
	e.rsims.UpdateFrom(rcore.CurrentExp)

	//
	log.Println("Starting new patient: ", e.rsims)
}

// ----------------------------------------------------------------------
// Deallocation of the current experiment
func (e *PMEntity) EndFunction() {
	//
	e.rsims = nil

	//
	log.Println("Ending the experiment")
}

// ----------------------------------------------------------------------
//
func pmEntityCFG() *PMEntity {
	//
	out := &PMEntity{}

	//
	return out
}
