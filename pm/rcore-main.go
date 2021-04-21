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
// Current patient and current Simulation state
var rpatient *rcore.Patient = nil
var rsims *rcore.SIMS = nil

// ----------------------------------------------------------------------
// For every cycle of distributed simulation.
func rserverCycle() {
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
func rserverStart() {
	//
	rpatient = rcore.NewPatient()

	//
	rpatient.SetDefaults()
	rpatient.WeightKG = rcore.Double(rcore.CurrentExp.Weight)

	//
	rsims = rcore.EmptySIMS(rpatient)

	//
	log.Println("Starting new patient: ", rsims)
}

// ----------------------------------------------------------------------
// Deallocation of the current experiment
func rserverEnd() {
	//
	rpatient = nil
	rsims = nil

	//
	log.Println("Ending the experiment")
}

// ----------------------------------------------------------------------
// main function (called from main() when arg is -s)
func rserverMain() {
	//
	rcore.EntityCore(rcore.CallPatMod, rserverCycle, rserverStart, rserverEnd)
}
