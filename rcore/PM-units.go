package rcore

//
import (
	"math"
)

// ----------------------------------------------------------------------
// type alias for float (perhaps skip and keep saying float64)
type Double = float64

// ----------------------------------------------------------------------
// 1 unit = 1000
// -1 unit = 1/1000
func conv(inval Double, inunit int, outunit int) Double {
	//
	for inunit != outunit {
		//
		if inunit < outunit {
			inval /= 1000.0
			inunit++
		} else {
			inval *= 1000.0
			inunit--
		}
	}

	//
	return inval
}

// ----------------------------------------------------------------------
// Weight = base-number + 10^3*unit
type Weight struct {
	//
	Value Double
	//
	Unit int
}

// ----------------------------------------------------------------------
// [g] == 0, gram
const (
	// kilo-gram
	Kg = 1
	//
	G = 0
	//
	Mg = -1
	Ug = -2
	Ng = -3
)

// ----------------------------------------------------------------------
//
type Volume struct {
	//
	Value Double
	//
	Unit int
}

// ----------------------------------------------------------------------
//
const (
	L  = 0
	ML = -1
	UL = -2
	NL = -3
)

// ----------------------------------------------------------------------
// transformation Weight->outunit
func (w Weight) In(outunit int) Weight {
	//
	return Weight{conv(w.Value, w.Unit, outunit), outunit}
}

// ----------------------------------------------------------------------
//
func (v Volume) In(outunit int) Volume {
	//
	return Volume{conv(v.Value, v.Unit, outunit), outunit}
}

// ----------------------------------------------------------------------
//
func Volume_0() Volume { return Volume{0, ML} }
func Weight_0() Weight { return Weight{0, ML} }

// ----------------------------------------------------------------------
// Hill function config
type Hill struct {
	//
	Emax Double

	//
	EC50  Double
	Gamma Double
}

// ----------------------------------------------------------------------
//
type Hill4 [4]Hill

// ----------------------------------------------------------------------
//
type Bounds struct {
	//
	Bmin Double
	Bmax Double
}

//
func (h Hill) TOFOutput100(valEffect Double) int {
	//
	return int(Bounds{0, 100}.Bound(100.0 - valEffect))
}

// ----------------------------------------------------------------------
//
func (b Bounds) In(v Double) bool {
	//
	return v >= b.Bmin && v <= b.Bmax
}

// ----------------------------------------------------------------------
// computed effect
func (h Hill) Value(inp Double) Double {
	//
	ip := math.Pow(inp, h.Gamma)
	ep := math.Pow(h.EC50, h.Gamma)

	//
	out := h.Emax * (ip / (ep + ip))

	//
	return math.Min(h.Emax, out)
}

// ----------------------------------------------------------------------
// computed effect
func (h Hill) IValue(inp Double) int {
	//
	return int(math.Floor(h.Value(inp) + 0.5))
}

// ----------------------------------------------------------------------
// Cinp needed to achieve 99% effect
// C^g >= E50^g * 1000
// C ~= E50 * 1000^{-g}
func (h Hill) Cinp100() Double {
	//
	return h.EC50 * math.Pow(100, 1.0/h.Gamma)
}

// ----------------------------------------------------------------------
//
func (b Bounds) Bound(v Double) Double {
	//
	return math.Min(math.Max(v, b.Bmin), b.Bmax)
}

// ----------------------------------------------------------------------
//
func (b Bounds) IBound(v int) int {
	//
	return int(math.Min(math.Max(Double(v), b.Bmin), b.Bmax))
}
