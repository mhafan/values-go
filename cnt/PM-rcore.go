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
// Current current Simulation state
var rsims *rcore.SIMS = nil

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
		"rtime", rsims.Time, "bolus ", _c.Bolus, " ", rsims.Bolus)

	// --------------------------------------------------------------------
	// reach Mtime in 1s simulation steps
	for rsims.Time <= _c.Mtime {
		// h = 1s, continuous simulation step
		rsims.RocSimStep()

		/*
			log.Println("HH ", rsims.Time, " ", rsims.Rocs.yROC, " ", rsims.Rocs.TOF0,
				rsims.Rocs.effect) */

		//
		_c.TOF = rsims.Rocs.TOF0
		_c.PTC = 0

		// +1s
		rsims = rsims.Next1S()
	}

	//
	rcore.CurrentExp.Save([]string{"TOF", "PTC"}, false)

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
// main function (called from main() when arg is -s)
func pmRCoreMAIN() {
	//
	fmt.Println("PM starting")

	//
	ent := rcore.Entity{
		rcore.CallPatMod,
		false,
		pmRCoreCycle, pmRCoreStart, pmRCoreEnd,
		func() {}}

	//
	rcore.EntityCore(&ent)
}
