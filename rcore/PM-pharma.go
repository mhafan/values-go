package rcore

// ----------------------------------------------------------------------
// Abstract class for PK/PD models of drugs
type Drug struct {
	//
}

// ----------------------------------------------------------------------
// Abstract class for PK/PD models of drugs - an Interface
type DrugDef interface {
	// ...
	InitialBolus(wkg int, wcoef Double, kratio Double) Volume
	InitialBolusExprec(expr *Exprec) Volume

	// transformations Weight <-> Volume
	SolutionUnits(w Weight) Volume
	WeightUnits(v Volume) Weight

	// Default Volume of Distribution
	// Default Hill function coeffs
	DefVd(absValue Double, unitValue Double, weight Weight) Volume
	DefHillCoefs() Hill
	DefHill4Coefs() Hill4
	DefIBolusMgPerKg() Double

	// PK/PD model 1s simulationn step
	SimStep(ss *SIMS)

	//
	Effect(cinp Double, vd Double, hill Hill) LinScale
}

// ----------------------------------------------------------------------
//
func MakeDrugDef(drugName string) DrugDef {
	//
	switch drugName {
	case DrugRocuronium:
		return Rocuronium{}
	default:
		return Rocuronium{}
	}
}
