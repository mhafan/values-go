// ----------------------------------------------------------------------
// Controller - model of the control algorithm
// (will be replaced by real TofCuff from RGB)
// ----------------------------------------------------------------------
// CNT currently implements simple infusion technique:
// 1) time=0 => initial bolus
// 2) regular time staps (flBolusInterval) => repetitive (smaller) boluses
// ------ (flBolusAmount)
// ----------------------------------------------------------------------

// ----------------------------------------------------------------------
// TODOs:
// - error/warning system between NMTSimulator and NMTSimulator-TCM
// - CNT - scheduling TOF/PTC measurement pattern
// - CNT - generating TOF/PTC correction for PM

package main

//
import (
	"flag"
	"log"
	"rcore"
)

// ----------------------------------------------------------------------
// Data structure for CNT Entity and its internal state
type CNTEntity struct {
	// CNT module provides default behavior for PUMP & SENSOR
	defaultBehavior bool

	// decision context
	decContext *DecContext

	// link to Patient Module Entity
	PM   *PMEntity
	CUFF *CUFFEntity
}

// ----------------------------------------------------------------------
// direct dosing of NMT blockator:
// time = 0: initial bolus
// then in "flBolusInterval" intervals, bolus "flBolusAmount"
var flBolusInterval = flag.Int("b", 200, "Interval between boluses [s]")
var flBolusAmount = flag.Int("B", 5, "Bolus volume [mL of solution]")

// ----------------------------------------------------------------------
// bypass PUMP & SENSOR
var flDefaultBehavior = flag.Bool("X", false, "Default PUMP/Cuff")

// ----------------------------------------------------------------------
// clear and reset the internal state
func (e *CNTEntity) ResetState() {
	//
	e.decContext = nil
	// call reset to Patient Model
	e.PM.ResetState()
}

// ----------------------------------------------------------------------
//
func (e *CNTEntity) MyTurn() string {
	//
	return rcore.CallCNT
}

// ----------------------------------------------------------------------
// with every START msg, do reset internals
func (e *CNTEntity) StartFunction() {
	// remove the previous one
	e.ResetState()

	// start a new context descriptor
	e.decContext = MakeDecContext()

	// initiate some REDIS record variables
	rcore.CurrentExp.StartupFromCNT()

	//
	log.Println("CNT Starting experiment: ", rcore.CurrentExp.Channel())

	//
	e.PM.StartFunction()

	//
	if e.defaultBehavior {
		//
		e.CUFF.StartFunction()
	}
}

// ----------------------------------------------------------------------
// do on END message (closing the experiment)
func (e *CNTEntity) EndFunction() {
	//
	rcore.CurrentExp.TerminateFromCNT()

	//
	log.Println("CNT Terminating experiment: ", rcore.CurrentExp.Channel())

	// nothing special
	e.ResetState()

	//
	e.PM.EndFunction()

	//
	if e.defaultBehavior {
		//
		e.CUFF.EndFunction()
	}
}

// ----------------------------------------------------------------------
// Default behavior for CNT-Entity
// 1) calling PM Entity on its commands
// 2) calling PUMP&Sensor entities if -X
func (e *CNTEntity) DefaultFunction(msg rcore.Rmsg) {
	// depending on the message in the channel
	switch msg.Message {
	// calling PM
	case e.PM.MyTurn():
		//
		e.PM.CycleFunction()

		//
	case rcore.CallPump:
		//
		if e.defaultBehavior {
			//
			rcore.CurrentExp.Say(rcore.CallPatMod)
		}

		//
	case rcore.CallSensor:
		//
		if e.defaultBehavior {
			//
			e.CUFF.CycleFunction()
		}
	}
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
func (e *CNTEntity) CycleFunction() {
	// dont forget to pass the token
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

	// --------------------------------------------------------------------
	// first of all, clear both outputs
	rcore.CurrentExp.Bolus = 0
	rcore.CurrentExp.Infusion = 0
	rcore.CurrentExp.SensorCommand = 0

	// load additional attributes for control algorithms
	cycleReloadCNTStrategyArgs()

	// --------------------------------------------------------------------
	// call decision making procedure where it should branch further
	dec := e.decContext.decision(rcore.CurrentExp, e.PM.rsims)

	// --------------------------------------------------------------------
	// if there is an nonzero output, write it to data rec
	if dec.BolusML > 0 || dec.InfusionML > 0 {
		//
		rcore.CurrentExp.Bolus = dec.BolusML
		rcore.CurrentExp.Infusion = dec.InfusionML

		//
		e.decContext.LastNonzero = &dec
	}

	//
	e.decContext.LastDecision = &dec

	// update the central data record
	rcore.CurrentExp.Save([]string{"bolus", "infusion", "SensorCommand"}, false)
}

// ----------------------------------------------------------------------
// well, main routine
func main() {
	// read and analyse input prompt args
	flag.Parse()

	// system data record describing an entity within the HiL
	// CNT, in this case
	ent := &CNTEntity{}
	ent.PM = pmEntityCFG()
	ent.CUFF = MakeCUFFEntity()

	// bypass PUMP & SENSOR if enabled
	// arg -X
	ent.defaultBehavior = *flDefaultBehavior

	// start listening CallCNT message ("CNT")
	// runLoop
	rcore.EntityCore(ent)
}
