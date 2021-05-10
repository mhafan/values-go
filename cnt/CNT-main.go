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
	"rcore"
)

// ----------------------------------------------------------------------
// direct dosing of NMT blockator:
// time = 0: initial bolus
// then in "flBolusInterval" intervals, bolus "flBolusAmount"
var flBolusInterval = flag.Int("b", 200, "Interval between boluses [s]")
var flBolusAmount = flag.Int("B", 5, "Bolus volume [mL of solution]")

// bypass PUMP & SENSOR
var flDefaultBehavior = flag.Bool("X", false, "Default PUMP/Cuff")

// ----------------------------------------------------------------------
//
var _decContext *DecContext = nil

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

	//
	_decContext = nil
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
	_decContext = MakeDecContext()

	//
	fmt.Println("CNT Start")
}

// ----------------------------------------------------------------------
// do on END message (closing the experiment)
func endWithExperiment() {
	// nothing special
	resetState()
}

// ----------------------------------------------------------------------
// reload CNTStrategy attributes depending on the current strategy
// (saving comp time when accessing REDIS)
func cycleReloadCNTStrategyArgs() {
	//
	rcore.EntityExpRecReload([]string{"targetCinpLow", "targetCinpHi"})

	//
	switch rcore.CurrentExp.CNTStrategy {
	//
	case rcore.CNTStratNone:
		//
	case rcore.CNTStratBasic:
		//
		rcore.EntityExpRecReload([]string{"repeStep", "repeBolus"})
		//
	case rcore.CNTStratFWSim:
		//
		rcore.EntityExpRecReload([]string{"fwRange"})
	}
}

// ----------------------------------------------------------------------
// Regulation cycle =>
// 1) direct mode
// 2) feedback mode
// ----------------------------------------------------------------------
// it needs to refresh data records:
// mtime/cycle
// CNTstrategy
// CNTstrategy args
func cycle() {
	// pass the token
	defer rcore.CurrentExp.Say(rcore.CallPump)

	// --------------------------------------------------------------------
	// ... first update
	if rcore.EntityExpRecReload([]string{"CNTStrategy"}) == false {
		//
		return
	}

	// strategy = none, do NOTHING
	if rcore.CurrentExp.CNTStrategy == rcore.CNTStratNone {
		//
		return
	}

	// --------------------------------------------------------------------
	// otherwise, CNT should make some decision
	if rcore.EntityExpRecReload([]string{"cycle", "mtime"}) == false {
		//
		return
	}

	// first of all, clear both outputs
	rcore.CurrentExp.Bolus = 0
	rcore.CurrentExp.Infusion = 0

	// load additional attributes for control algorithms
	cycleReloadCNTStrategyArgs()

	// --------------------------------------------------------------------
	// call decision making procedure where it should branch further
	dec := _decContext.decision(rcore.CurrentExp, rsims)

	// --------------------------------------------------------------------
	// if there is an nonzero output, write it to data rec
	if dec.BolusML > 0 || dec.InfusionML > 0 {
		//
		rcore.CurrentExp.Bolus = dec.BolusML
		rcore.CurrentExp.Infusion = dec.InfusionML

		//
		_decContext.LastNonzero = &dec
	}

	//
	_decContext.LastDecision = &dec

	// update the central data record
	rcore.CurrentExp.Save([]string{"bolus", "infusion"}, false)
}

// ----------------------------------------------------------------------
// Bypassing PUMP and SENSOR
// (if promp arg -X entered)
func PumpCuffDefaultBehavior(msg rcore.Rmsg) {
	//
	switch msg.Message {
	case rcore.CallPump:
		//
		rcore.CurrentExp.Say(rcore.CallPatMod)

	case rcore.CallSensor:
		//
		rcore.CurrentExp.Say(rcore.CallTCM)
	}
}

// ----------------------------------------------------------------------
// well, main routine
func main() {
	// read and analyse input prompt args
	flag.Parse()

	// system data record describing an entity within the HiL
	// CNT, in this case
	ent := rcore.MakeNewEntity()

	// the entity gets activated on msg:
	ent.MyTurn = rcore.CallCNT

	// link 3 procedures implementing
	// START command
	// CNT command (for each cycle)
	// END command
	ent.WhatStart = startupWithExperiment
	ent.What = cycle
	ent.WhatEnd = endWithExperiment

	// secondary entities (PATMOD)
	ent.Slave = pmEntityCFG()

	// bypass PUMP & SENSOR if enabled
	if *flDefaultBehavior {
		//
		ent.WhatDefault = PumpCuffDefaultBehavior
	}

	// start listening CallCNT message ("CNT")
	// runLoop
	rcore.EntityCore(ent)
}
