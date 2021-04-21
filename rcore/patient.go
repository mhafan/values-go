package rcore

// ----------------------------------------------------------------------
//
type PatientROC struct {
	// [mL/kg]
	Vd_kg Double
	// [ug/mL]
	ec50 Double
}

//
func PatientROC_def() PatientROC {
	//
	return PatientROC{Vd_kg: 38, ec50: 0.823}
}

// ----------------------------------------------------------------------
// Patient descriptor
// ----------------------------------------------------------------------
type Patient struct {
	// [kg]
	WeightKG Double

	//
	Age int
	Sex string

	//
	RocCFG PatientROC
}

// ----------------------------------------------------------------------
//
func NewPatient() *Patient { return new(Patient) }

// ----------------------------------------------------------------------
//
func (p *Patient) SetDefaults() {
	//
	p.Age = -1
	p.Sex = SexMale
	p.RocCFG = PatientROC_def()
}

// ----------------------------------------------------------------------
//
func (p *Patient) Vc_roc() Volume {
	//
	return Volume{38.0 * p.WeightKG, ML}
}
