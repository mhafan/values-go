// ----------------------------------------------------------------------
// Math part of CNT
// ----------------------------------------------------------------------
// Decision making system
package main

import (
	"fmt"
	"rcore"
)

// ----------------------------------------------------------------------
// decision made about the infusion. This is an output from decision
// procedure.
type Decision struct {
	// valid at mtime
	Mtime int

	// [mL] of solution
	InfusionML int
	BolusML    int
}

// ----------------------------------------------------------------------
// Decision context valid during total START->...->END experiment
type DecContext struct {
	// ------------------------------------------------------------------
	// the initial bolus is special
	// it may be a bigger dosage in order to achieve intubation state
	InitialBolusGiven bool
	InitialBolusMTime int

	// ------------------------------------------------------------------
	// context of the previous decision
	LastDecision *Decision
	LastNonzero  *Decision

	// ------------------------------------------------------------------
	// CNT controls TOF/PTC sensor
	LastScheduledTOFMeasurement int
	LastScheduledPTCMeasurement int
}

// ----------------------------------------------------------------------
//
func MakeDecContext() *DecContext {
	//
	out := DecContext{}

	//
	out.InitialBolusGiven = false
	out.LastDecision = nil
	out.LastNonzero = nil

	//
	out.LastScheduledPTCMeasurement = 0
	out.LastScheduledTOFMeasurement = 0

	//
	return &out
}

// ----------------------------------------------------------------------
//
func makeDecision(expID *rcore.Exprec) Decision {
	//
	return Decision{Mtime: expID.Mtime}
}

// ----------------------------------------------------------------------
// assuming
func (context *DecContext) lastAnyBolusAt() int {
	// skip
	if context.InitialBolusGiven == false {
		//
		return -1
	}

	// time of last bolus is either initial bolus
	lastBol := context.InitialBolusMTime

	//
	if context.LastNonzero != nil {
		//
		lastBol = context.LastNonzero.Mtime
	}

	//
	return lastBol
}

// ----------------------------------------------------------------------
// Basic algorithm
// - it schedules a bolus to be administered at (lastBolus+expID.repeStep)
// - it skips if inital bolus hasnt been administered yet
func (context *DecContext) decisionBasic(expID *rcore.Exprec, insim *rcore.SIMS) Decision {
	//
	out := makeDecision(expID)

	//
	lastBolusAt := context.lastAnyBolusAt()

	//
	if lastBolusAt >= 0 {
		//
		if expID.Mtime >= lastBolusAt+expID.RepeStep {
			//
			out.BolusML = expID.RepeBolus
		}
	}

	//
	return out
}

// ----------------------------------------------------------------------
// CNT Algorithm: Forward Simulation
func (context *DecContext) decisionFWSim(expID *rcore.Exprec, insim *rcore.SIMS) Decision {
	//
	out := Decision{Mtime: expID.Mtime}

	// ------------------------------------------------------------------
	// internal state & config
	fwcfg := FWSimMake(insim, expID)

	// ------------------------------------------------------------------
	// default predicate
	// true => quit the search
	// false => continue search
	fwcfg.pred = func(sim *rcore.SIMS) bool {
		//
		return sim.CinpAboveLow(expID)
	}

	// run the bolus search
	suggestBolus := forwardSimulationBolus(fwcfg)

	//
	out.BolusML = suggestBolus

	//
	return out
}

// ----------------------------------------------------------------------
// Decision branching point:
// input:
// -- current exprec
// -- last SIM state from PATMOD, assert(SIM.Time <= expID.Mtime)
func (context *DecContext) decision(expID *rcore.Exprec, insim *rcore.SIMS) Decision {
	// ------------------------------------------------------------------
	//
	out := Decision{Mtime: expID.Mtime}

	// ------------------------------------------------------------------
	// no initial bolus given so far
	if context.InitialBolusGiven == false {
		//
		ibolus := rcore.InitialBolusExprec(expID)

		//
		if ibolus.Value > 0 {
			//
			fmt.Println("IBOLUS ", ibolus.Value)

			//
			out.BolusML = int(ibolus.Value)
			context.InitialBolusGiven = true
			context.InitialBolusMTime = expID.Mtime

			//
			return out
		}
	}

	// ------------------------------------------------------------------
	//
	switch expID.CNTStrategy {
	//
	case rcore.CNTStratBasic:
		//
		out = context.decisionBasic(expID, insim)
		//
	case rcore.CNTStratFWSim:
		//
		out = context.decisionFWSim(expID, insim)
	}

	//
	return out
}
