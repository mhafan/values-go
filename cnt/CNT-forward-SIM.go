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
func FWSimMake(insim *rcore.SIMS) FWSim {
	//
	out := FWSim{}

	//
	out.insim = insim
	out.sbolusMax = 20
	out.forwardTime = 60 * 5

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
		fmt.Println("FWSIM ", _t0, " ", fwc.forwardTime, " bolus=", sbolus, " ", result)

		//
		if fwc.pred(result) {
			//
			return sbolus
		}
	}

	//
	return 0
}
