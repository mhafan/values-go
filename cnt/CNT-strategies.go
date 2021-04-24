// ----------------------------------------------------------------------
// Math part of CNT
// ----------------------------------------------------------------------
// Decision making system
package main

import (
	"rcore"
)

// ----------------------------------------------------------------------
//
type Decision struct {
	//
	Mtime int

	//
	InfusionML int
	BolusML    int
}

// ----------------------------------------------------------------------
//
func MakeDecContext() *DecContext {
	//
	out := DecContext{}

	//
	out.InitialBolusGiven = false
	out.LastDecision = nil

	//
	return &out
}

// ----------------------------------------------------------------------
//
type DecContext struct {
	//
	InitialBolusGiven bool

	//
	LastDecision *Decision
}

// ----------------------------------------------------------------------
//
func (context *DecContext) decision(expID *rcore.Exprec, insim *rcore.SIMS) Decision {
	// ------------------------------------------------------------------
	//
	out := Decision{}

	// ------------------------------------------------------------------
	//
	if context.InitialBolusGiven == false {
		//
		ibolus := rcore.InitialBolus(expID.Drug, expID.Weight, 1.0)

		//
		if ibolus.Value > 0 {
			//
			out.BolusML = int(ibolus.Value)
			context.InitialBolusGiven = true

			//
			return out
		}
	}

	///
	fwcfg := FWSimMake(insim)

	//
	fwcfg.pred = SimPredicateCinp

	suggestBolus := forwardSimulationBolus(fwcfg)

	//
	out.BolusML = suggestBolus

	//
	return out
}
