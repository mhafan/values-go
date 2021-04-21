package rcore

//
import (
	"math"
)

// ----------------------------------------------------------------------
// Hill function config
type Hill struct {
	//
	emax Double

	//
	ec50  Double
	gamma Double
}

// ----------------------------------------------------------------------
//
type Bounds struct {
	//
	bmin Double
	bmax Double
}

// ----------------------------------------------------------------------
// computed effect
func (h Hill) value(inp Double) Double {
	//
	ip := math.Pow(inp, h.gamma)
	ep := math.Pow(h.ec50, h.gamma)

	//
	out := h.emax * (ip / (ep + ip))

	//
	return math.Min(h.emax, out)
}

//
func (b Bounds) bound(v Double) Double {
	//
	return math.Min(math.Max(v, b.bmin), b.bmax)
}
