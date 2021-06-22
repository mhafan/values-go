package rcore

// ----------------------------------------------------------------------
//
type Vec4 [4]int

// ----------------------------------------------------------------------
// Complex record about TOF measurement
// ----------------------------------------------------------------------
// LinScale measure:
// - TOFratio = 100...0
// - TOFcount = 3, 2, 1, 0
// - PTCcount = 15, 14, ..., 0
// ----------------------------------------------------------------------
// Effect = effect of the drug from 0 - 100%
// Amplitude = dropping from 100 to 0 with rising effect
// Amplitude = 100-Effect
type LinScale struct {
	// input concentration
	Cinp Double

	//
	TOF4Effect    Vec4
	TOF4Amplitude Vec4
	PTCCount      int
}

// ----------------------------------------------------------------------
//
func Norma100(v int) Double { return Double(v) / 100.0 }

//
func (l LinScale) TOFSimpleAmplitude() int {
	//
	return l.TOF4Amplitude[0]
}

// ----------------------------------------------------------------------
//
func (l LinScale) TOFcount() int {
	//
	out := 0

	//
	for i := 0; i < 4; i++ {
		//
		if l.TOF4Amplitude[i] > 0 {
			out += 1
		}
	}

	//
	return out
}

// ----------------------------------------------------------------------
//
func (l LinScale) TOFratio() Double {
	//
	if l.TOF4Amplitude[0] >= 10 {
		//
		return Norma100(l.TOF4Amplitude[3]) / Norma100(l.TOF4Amplitude[0])
	}

	//
	return 0
}

// ----------------------------------------------------------------------
// linscale value
func (ls LinScale) Value() int {
	//
	return 0
}

// ----------------------------------------------------------------------
// inp - input Cinp
func LinScaleModel4(inp Double, hill Hill4) LinScale {
	//
	out := LinScale{Cinp: inp}
	bnds := Bounds{0, 100}

	//
	for i := 0; i < 4; i++ {
		//
		out.TOF4Effect[i] = hill[i].IValue(out.Cinp)
		out.TOF4Amplitude[i] = bnds.IBound(100 - out.TOF4Effect[i])
	}

	//
	return out
}
