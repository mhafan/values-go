package rcore

// ----------------------------------------------------------------------
// Initil bolus as defined by manufacturer
func InitialBolus(drug string, wkg int, wcoef Double, kratio Double) Volume {
	//
	_wkg := Double(wkg)

	//
	switch drug {
	case DrugRocuronium:
		// 0.6 mg per [kg] of patient's weight
		return RocWSOL(Weight{kratio * _wkg * wcoef, Mg}).In(ML)
	case DrugCisatracurium:
		// TODO
		return Volume{0, ML}
	}

	// default value if drug is set incorrectly
	return Volume{0, ML}
}

// ----------------------------------------------------------------------
// Initil bolus as defined by manufacturer
func InitialBolusExprec(expr *Exprec) Volume {
	//
	return InitialBolus(expr.Drug, expr.Weight, expr.Wcoef, expr.IbolusMg)
}
