package main

//
import "math"

// ----------------------------------------------------------------------
// Hill function config
type Hill struct {
  //
  emax  double
  //
  ec50  double
  gamma double
}

// ----------------------------------------------------------------------
//
type Bounds struct {
  bmin  double
  bmax  double
}

// ----------------------------------------------------------------------
// computed effect
func (h Hill) value(inp double) double {
  //
  ip := math.Pow(inp, h.gamma)
  ep := math.Pow(h.ec50, h.gamma)

  //
  out := h.emax * (ip / (ep + ip))

  //
  return math.Min(h.emax, out)
}

//
func (b Bounds) bound(v double) double {
  //
  return math.Min(math.Max(v, b.bmin), b.bmax)
}
