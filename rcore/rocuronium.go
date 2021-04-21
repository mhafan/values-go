package rcore

//
func RocDefHill() Hill {
	//
	return Hill{100.0, 0.823, 4.79}
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
func RocInputs(y COMP_X, infConc Double) COMP_X {
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
func RocSimStep1H(yin COMP_X, infConc Double) COMP_X {
	//
	f := RocInputs(yin, infConc)

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
func (ss *SIMS) RocSimStep() {
	// Volume of distribution
	Vd := ss.Patient.Vc_roc()

	// eventual infusion input
	ic := 0.0

	// [ml/hr] => effective weight / hr => per s =>
	if ss.Infusion.Value > 0.0 {
		// [ug/hr]
		inWeightHour := RocSOLW(ss.Infusion).In(Ug)

		// [ug/s]
		inWeightS := inWeightHour.Value / 3600.0

		//
		ic = inWeightS / Vd.In(ML).Value
	}

	//
	if ss.Time > 0 {
		//
		ss.Rocs.yROC = RocSimStep1H(ss.Rocs.yROC, ic)
	}

	//
	if ss.Bolus.Value > 0.0 {
		//
		ss.Rocs.yROC[1] += RocSOLW(ss.Bolus).In(Ug).Value / Vd.In(ML).Value
	}

	//
	ss.Rocs.effect = ss.Rocs.rocHill.value(ss.Rocs.yROC[1])
	ss.Rocs.TOF0 = int(_TOFbounds.bound(100.0 - ss.Rocs.effect))
}