package rcore

// ----------------------------------------------------------------------
// Initil bolus as defined by manufacturer
func InitialBolus(drug string, wkg int, wcoef Double) Volume {
	//
	_wkg := Double(wkg)

	//
	switch drug {
	case DrugRocuronium:
		// 0.6 mg per [kg] of patient's weight
		return RocWSOL(Weight{0.6 * _wkg * wcoef, Mg}).In(ML)
	case DrugCisatracurium:
		// TODO
		return Volume{0, ML}
	}

	// default value if drug is set incorrectly
	return Volume{0, ML}
}
