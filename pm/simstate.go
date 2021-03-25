package main

import "rcore"

// ----------------------------------------------------------------------
// 3 compartment model
type COMP_X  [3+1] rcore.Double

//
var _TOFbounds = Bounds{0, 100}

// ----------------------------------------------------------------------
//
type ROCS struct {
  // --------------------------------------------------------------------
  // model, internal
  yROC      COMP_X
  rocHill   Hill

  //
  effect    rcore.Double
  TOF0      int
}

// ----------------------------------------------------------------------
// Simulation state:
// ----------------------------------------------------------------------
type SIMS struct {
  // --------------------------------------------------------------------
  //
  patient *Patient

  // --------------------------------------------------------------------
  // simulation internal data
  time int

  //
  rocs      ROCS

  // --------------------------------------------------------------------
  // inputs
  bolus     rcore.Volume
  infusion  rcore.Volume
}


// ----------------------------------------------------------------------
// time(0) zero simstate constructor
func emptySIMS(pat *Patient) *SIMS {
  //
  return &SIMS {
    //
    pat, 0,
    //
    ROCS { COMP_X{}, rocDefHill(), 0, 0 },
    // bolus & infusion
    rcore.Volume_0(), rcore.Volume_0() }
}


// ----------------------------------------------------------------------
// next state by shifting time
func (from *SIMS) nextState(at int) *SIMS {
  // ... copy ...
  ns := *from

  // shift time
  ns.time = at
  // reset inputs
  ns.bolus = rcore.Volume_0()
  ns.infusion = rcore.Volume_0()

  //
  return &ns
}

// ----------------------------------------------------------------------
//
func (from *SIMS) next1S() *SIMS {
  //
  return from.nextState(from.time + 1)
}
