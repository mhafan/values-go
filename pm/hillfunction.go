package main

//
import (
	"math"
	"rcore"
)

// ----------------------------------------------------------------------
// Hill function config
type Hill struct {
	//
	emax rcore.Double

	//
	ec50  rcore.Double
	gamma rcore.Double
}

// ----------------------------------------------------------------------
//
type Bounds struct {
	//
	bmin rcore.Double
	bmax rcore.Double
}

// ----------------------------------------------------------------------
// computed effect
func (h Hill) value(inp rcore.Double) rcore.Double {
	//
	ip := math.Pow(inp, h.gamma)
	ep := math.Pow(h.ec50, h.gamma)

	//
	out := h.emax * (ip / (ep + ip))

	//
	return math.Min(h.emax, out)
}

//
func (b Bounds) bound(v rcore.Double) rcore.Double {
	//
	return math.Min(math.Max(v, b.bmin), b.bmax)
}
