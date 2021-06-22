package rcore

//
type Rocuronium Drug

// ----------------------------------------------------------------------
//
func (d Rocuronium) InitialBolus(wkg int, wcoef Double, kratio Double) Volume {
	// 0.6 mg per [kg] of patient's weight
	return d.SolutionUnits(Weight{kratio * Double(wkg) * wcoef, Mg}.In(ML))
}

// ----------------------------------------------------------------------
//
func (d Rocuronium) InitialBolusExprec(expr *Exprec) Volume {
	//
	return d.InitialBolus(expr.Weight, expr.Wcoef, expr.IbolusMg)
}

// ----------------------------------------------------------------------
//
func (d Rocuronium) SolutionUnits(w Weight) Volume {
	//
	return Volume{w.In(Mg).Value / 10.0, ML}
}

// ----------------------------------------------------------------------
//
func (d Rocuronium) WeightUnits(v Volume) Weight {
	//
	return Weight{v.In(ML).Value * 10.0, Mg}
}

// ----------------------------------------------------------------------
//
func (d Rocuronium) DefHillCoefs() Hill {
	//
	return Hill{100.0, 0.823, 4.79}
}

func (d Rocuronium) DefHill4Coefs() Hill4 {
	//
	p := d.DefHillCoefs()

	//
	fait := func(a, b Double, h Hill) Hill {
		//
		return Hill{h.Emax, h.EC50 / a, h.Gamma / b}
	}

	//
	return Hill4{p, fait(1.1, 1.1, p), fait(1.3, 1.25, p), fait(1.5, 1.33, p)}
}

//
func (d Rocuronium) DefIBolusMgPerKg() Double {
	//
	return 0.6
}

// ----------------------------------------------------------------------
//
func (d Rocuronium) DefVd(absValue Double, unitValue Double, weight Weight) Volume {
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
//
func (d Rocuronium) Effect(cinp Double, vd Double, hill Hill) LinScale {
	//
	out := LinScaleModel4(cinp, d.DefHill4Coefs())

	// ------------------------------------------------------------------
	// Simple output from the hill function
	//out.HillEffect = hill.Value(out.Cinp)
	//out.TOFSimple = hill.TOFOutput100(out.HillEffect)

	// ------------------------------------------------------------------
	// 1) Let us assume TOFcount = 4 (TOF[4]>0)

	//
	return out
}

// ----------------------------------------------------------------------
//
const (
	//
	rocK12 = 0.259 / 60.0
	rocK21 = 0.163 / 60.0
	rocK13 = 0.060 / 60.0
	rocK31 = 0.012 / 60.0
	rocK10 = 0.119 / 60.0
)

// ----------------------------------------------------------------------
//
func __rocInputs(y COMP_X, infConc Double) COMP_X {
	//
	var out COMP_X

	//
	out[1] = y[2]*rocK21 + y[3]*rocK31 - y[1]*(rocK10+rocK12+rocK13)
	out[2] = y[1]*rocK12 - y[2]*rocK21
	out[3] = y[1]*rocK13 - y[3]*rocK31

	//
	out[1] += infConc

	//
	return out
}

// ----------------------------------------------------------------------
//
func __rocSimStep1H(yin COMP_X, infConc Double) COMP_X {
	//
	f := __rocInputs(yin, infConc)

	//
	var out COMP_X

	//
	for i := 1; i < 4; i++ {
		//
		out[i] = yin[i] + f[i]
	}

	//
	return out
}

// ----------------------------------------------------------------------
// PK/PD Model for Rocuronium
// ----------------------------------------------------------------------
func (d Rocuronium) SimStep(ss *SIMS) {
	// Volume of distribution
	Vd := ss.VdCentral

	// eventual infusion input
	ic := 0.0

	// [ml/hr] => effective weight / hr => per s =>
	if ss.Infusion.Value > 0.0 {
		// [ug/hr]
		inWeightHour := ss.Drug.WeightUnits(ss.Infusion).In(Ug)

		// [ug/s]
		inWeightS := inWeightHour.Value / 3600.0

		//
		ic = inWeightS / Vd.In(ML).Value
	}

	//
	if ss.Time > 0 {
		//
		ss.YROC = __rocSimStep1H(ss.YROC, ic)
	}

	//
	if ss.Bolus.Value > 0.0 {
		//
		ss.YROC[1] += ss.Drug.WeightUnits(ss.Bolus).In(Ug).Value / Vd.In(ML).Value

		//
		ss.BolusConsumptionML += ss.Bolus.Value

		//
		//fmt.Println("BOLSTAT ", ss.BolusConsumptionML, ss.Bolus.Value)
	}

	//
	//ss.Effect = d.Effect(ss.Cinp(), Vd.In(ML).Value, ss.HillCoefs)
	ss.Effect = LinScaleModel4(ss.Cinp(), d.DefHill4Coefs())
}
