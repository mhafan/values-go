// ----------------------------------------------------------------------
// PatMod for R-system
// ----------------------------------------------------------------------
package main

//
import (
	"log"
	"rcore"
)

// ----------------------------------------------------------------------
// Current current Simulation state
var rsims *rcore.SIMS = nil

// ----------------------------------------------------------------------
//
func pmTotalRecoveryPredicate(sims *rcore.SIMS) bool {
	//
	return sims.TOF0 < 95
}

// ----------------------------------------------------------------------
// For every cycle of distributed simulation.
func pmRCoreCycle() {
	//
	tol := []string{"mtime", "cycle", "bolus", "infusion"}

	// --------------------------------------------------------------------
	// load Redis record
	if rcore.EntityExpRecReload(tol) == false {
		//
		panic("REDIS record not found")
	}

	// ...
	var _c = rcore.CurrentExp

	// --------------------------------------------------------------------
	// REDIS record -> SIMS (patmod simulation state)
	rsims.UpdateFrom(_c)

	//
	log.Println("PMA; cycle: ", _c.Cycle, "mtime", _c.Mtime,
		"rtime", rsims.Time, "bolus ", rsims.Bolus.Value, " ", rsims.Bolus)

	//
	rsims = rsims.SimSteps(_c.Mtime)

	//
	trec := rsims.Clone().SimStepsWhile(pmTotalRecoveryPredicate)

	//
	log.Println("TREC ", trec.Time)

	// --------------------------------------------------------------------
	//
	_c.TOF = rsims.TOF0
	_c.PTC = 0
	_c.Cinp = rsims.YROC[1]
	_c.ConsumedTotal += rsims.BolusConsumptionML
	_c.RecoveryTime = (trec.Time - rsims.Time)

	//
	rcore.CurrentExp.Save([]string{"TOF", "PTC", "Cinp", "ConsumedTotal", "RecoveryTime"}, false)

	//
	_c.Say(rcore.CallSensor)
}

// ----------------------------------------------------------------------
// vm.X.Y -> START
// --- Initialization of new experiment
func pmRCoreStart() {
	//
	rsims = rcore.EmptySIMS()

	//
	rsims.SetupFrom(rcore.CurrentExp)
	rsims.UpdateFrom(rcore.CurrentExp)

	//
	log.Println("Starting new patient: ", rsims)
}

// ----------------------------------------------------------------------
// Deallocation of the current experiment
func pmRCoreEnd() {
	//
	rsims = nil

	//
	log.Println("Ending the experiment")
}

// ----------------------------------------------------------------------
//
func pmEntityCFG() *rcore.Entity {
	//
	ent := rcore.MakeNewEntity()

	//
	ent.MyTurn = rcore.CallPatMod
	ent.What = pmRCoreCycle
	ent.WhatStart = pmRCoreStart
	ent.WhatEnd = pmRCoreEnd

	//
	return ent
}
