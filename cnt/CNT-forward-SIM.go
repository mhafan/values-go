package main

import (
	"fmt"
	"rcore"
)

// ----------------------------------------------------------------------
//
func SimPredicateCinp(sim *rcore.SIMS) bool {
	//
	return sim.Cinp() >= 2
}

// ----------------------------------------------------------------------
// Forward Simulation CONFIG
type FWSim struct {
	//
	insim       *rcore.SIMS
	forwardTime int
	pred        rcore.SimPredicate

	//
	sbolusMax int
}

// ----------------------------------------------------------------------
//
func FWSimMake(insim *rcore.SIMS, inexp *rcore.Exprec) FWSim {
	//
	out := FWSim{}

	//
	out.insim = insim
	out.sbolusMax = 20
	out.forwardTime = inexp.FwRange

	//
	return out
}

// ----------------------------------------------------------------------
// Forward Simulator
// - increasing bolus at t0
func forwardSimulationBolus(fwc FWSim) int {
	// entering time in the simulation context
	_t0 := fwc.insim.Time

	// increasing bolus and simulating forward
	for sbolus := 0; sbolus < fwc.sbolusMax; sbolus++ {
		// ... copy the current state
		insimClone := fwc.insim.Clone()

		// prepare for forward simulation
		insimClone.CinpStat.Reset()
		// ... with that given bolus
		insimClone.Bolus = rcore.Volume{rcore.Double(sbolus), rcore.ML}

		// do simulation steps
		result := insimClone.SimSteps(_t0 + fwc.forwardTime)

		//
		fmt.Println("FWSIM ", rcore.CurrentExp.TargetCinpLow, " bolus=", sbolus, "range", fwc.forwardTime, " ", result)

		//
		if fwc.pred(result) {
			//
			return sbolus
		}
	}

	//
	return 0
}
