// ----------------------------------------------------------------------
// Controller - model of the control algorithm
// (will be replaced by real TofCuff from RGB)
// ----------------------------------------------------------------------
// CNT currently implements simple infusion technique:
// 1) time=0 => initial bolus
// 2) regular time staps (flBolusInterval) => repetitive (smaller) boluses
// ------ (flBolusAmount)
// ----------------------------------------------------------------------

package main

//
import (
	"flag"
	"fmt"
	"log"
	"rcore"
)

// ----------------------------------------------------------------------
// direct dosing of NMT blockator:
// time = 0: initial bolus
// then in "flBolusInterval" intervals, bolus "flBolusAmount"
var flBolusInterval = flag.Int("b", 200, "Interval between boluses [s]")
var flBolusAmount = flag.Int("B", 5, "Bolus volume [mL of solution]")

// ----------------------------------------------------------------------
// state of the experiment:
// time of the last bolus
var _lastTimeBolus = 0

// time of the next scheduled bolus
var _scheduledBolusAt = -1

// was any initial bolus given
var _initialBolusGiven = false

// ----------------------------------------------------------------------
//
func resetState() {
	//
	_lastTimeBolus = 0
	_scheduledBolusAt = -1
	_initialBolusGiven = false
}

// ----------------------------------------------------------------------
// Initil bolus as defined by manufacturer
func initialBolus(drug string, wkg int) rcore.Volume {
	//
	switch drug {
	case rcore.DrugRocuronium:
		// 0.6 mg per [kg] of patient's weight
		return rcore.RocWSOL(rcore.Weight{0.6 * 100.0, rcore.Mg}).In(rcore.ML)
	case rcore.DrugCisatracurium:
		// TODO
		return rcore.Volume{0, rcore.ML}
	}

	// default value if drug is set incorrectly
	return rcore.Volume{0, rcore.ML}
}

// ----------------------------------------------------------------------
// with every START msg, do reset internals
func startupWithExperiment() {
	//
	resetState()

	// schedule next bolus
	if *flBolusInterval > 0 {
		//
		_scheduledBolusAt = *flBolusInterval
	}

	//
	fmt.Println("CNT Start")
}

// ----------------------------------------------------------------------
// do on END message (closing the experiment)
func endWithExperiment() {
	// nothing special
}

// ----------------------------------------------------------------------
// Direct MODE:
// time == 0 => INITIAL bolus
// intervals
func regulationInDirectMode(_r *rcore.Exprec) bool {
	// by defualt, set both zero
	_r.Bolus = 0
	_r.Infusion = 0

	// --------------------------------------------------------------------
	// time == 0 || Cycle == 0
	// --------------------------------------------------------------------
	if _r.Mtime <= 0 || _r.Cycle == 0 {
		//
		if _initialBolusGiven == true {
			//
			log.Println("Will not set the initial bolus multiple times!")

			// error
			return false
		}

		// initial bolus in recommended volume 0.6mg/kg
		_r.Bolus = int(initialBolus(_r.Drug, _r.Weight).Value)
		_initialBolusGiven = true

		//
		log.Println("CNT:initial bolus [mL]: ", _r.Bolus)

		//
		return true
	}

	// --------------------------------------------------------------------
	// repetitive bolus, if enabled:
	// the time has reached scheduled moment
	if _r.Mtime >= _scheduledBolusAt && _scheduledBolusAt > 0 {
		// now
		_lastTimeBolus = _r.Mtime
		// schedule the next moment
		_scheduledBolusAt = _r.Mtime + (*flBolusInterval)
		// set the bolus
		_r.Bolus = *flBolusAmount

		//
		log.Println("CNT:repetitive bolus [mL]:", _r.Bolus, " time=", _r.Mtime)

		//
		return true
	}

	//
	return true
}

// ----------------------------------------------------------------------
// Regulation cycle =>
// 1) direct mode
// 2) feedback mode
func cycle() {
	// --------------------------------------------------------------------
	// update the exp variable, rcore.CurrentExp
	if rcore.EntityExpRecReload([]string{"cycle", "mtime"}) == false {
		//
		return
	}

	// --------------------------------------------------------------------
	// step of regulation in simple/direct regime
	if regulationInDirectMode(rcore.CurrentExp) == true {
		// update bolus/infusion
		rcore.CurrentExp.Save([]string{"bolus", "infusion"}, false)
	}

	//
	rcore.CurrentExp.Say(rcore.CallPump)
}

// ----------------------------------------------------------------------
//
func main() {
	//
	flag.Parse()

	//
	ent := rcore.MakeNewEntity()

	//
	ent.MyTurn = rcore.CallCNT
	ent.IsMaster = true
	ent.What = cycle
	ent.WhatStart = startupWithExperiment
	ent.WhatEnd = endWithExperiment

	//
	ent.Slave = pmEntityCFG()

	// start listening CallCNT message ("CNT")
	rcore.EntityCore(ent)
}
