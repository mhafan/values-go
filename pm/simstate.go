package main

//
import "math"

// ----------------------------------------------------------------------
// 3 compartment model
type COMP_X  [3+1] double

//
func bound(v double, vmin double, vmax double) double {
  //
  return math.Min(math.Max(v, vmin), vmax)
}

// ----------------------------------------------------------------------
//
type Hill struct {
  //
  emax  double
  ec50  double
  gamma double
}

// ----------------------------------------------------------------------
//
func (h Hill) value(inp double) double {
  //
  ip := math.Pow(inp, h.gamma)
  ep := math.Pow(h.ec50, h.gamma)

  //
  out := h.emax * (ip / (ep + ip))

  //
  return math.Min(h.emax, out)
}

// ----------------------------------------------------------------------
//
type ROCS struct {
  // --------------------------------------------------------------------
  // model, internal
  yROC      COMP_X
  rocHill   Hill

  //
  effect    double
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
  bolus     volume
  infusion  volume
}


// ----------------------------------------------------------------------
// time(0) zero simstate constructor
func emptySIMS(pat *Patient) SIMS {
  //
  return SIMS {
    //
    pat, 0,
    //
    ROCS { COMP_X{}, rocDefHill(), 0, 0 },
    // bolus & infusion
    volume_0(), volume_0() }
}


// ----------------------------------------------------------------------
// next state by shifting time
func (from SIMS) nextState(at int) SIMS {
  // ... copy ...
  ns := from

  // shift time
  ns.time = at
  // reset inputs
  ns.bolus = volume_0()
  ns.infusion = volume_0()

  //
  return ns
}

// ----------------------------------------------------------------------
//
func (from SIMS) next1S() SIMS {
  //
  return from.nextState(from.time + 1)
}
