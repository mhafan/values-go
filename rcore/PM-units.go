package rcore

//
import (
	"math"
)

// ----------------------------------------------------------------------
// type alias for float (perhaps skip and keep saying float64)
type Double = float64

// ----------------------------------------------------------------------
// Complex record about TOF measurement
// ----------------------------------------------------------------------
// LinScale measure:
// - TOFratio = 100...0
// - TOFcount = 3, 2, 1, 0
// - PTCcount = 15, 14, ..., 0
// ----------------------------------------------------------------------
//
type LinScale struct {
	//
	Cinp Double

	//
	TOFRatio int
	TOFCount int
	PTCCount int

	//
	HillEffect Double
}

// ----------------------------------------------------------------------
//
func (ls LinScale) value() int {
	//
	return 0
}

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
//
func VdFor(drug string, absValue Double, unitValue Double, weight Weight) Volume {
	//
	if absValue > 0 {
		//
		return Volume{absValue, L}
	}

	//
	if unitValue > 0 {
		//
		return Volume{unitValue * weight.In(Kg).Value, ML}
	}

	//
	return Volume{38.0 * weight.In(Kg).Value, ML}
}

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

// ----------------------------------------------------------------------
//
func (b Bounds) bound(v Double) Double {
	//
	return math.Min(math.Max(v, b.bmin), b.bmax)
}

// ----------------------------------------------------------------------
// inp - input Cinp
func LinScaleModel(inp Double, hill Hill) LinScale {
	//
	out := LinScale{}

	//
	out.HillEffect = hill.value(inp)
	out.TOFRatio = int(_TOFbounds.bound(100.0 - out.HillEffect))

	//
	return out
}
