package main


// ----------------------------------------------------------------------
//
type PatientROC struct {
  // [mL/kg]
  Vd_kg     double
  // [ug/mL]
  ec50      double
}

//
func PatientROC_def() PatientROC {
  //
  return PatientROC{ Vd_kg: 38, ec50: 0.823 }
}

// ----------------------------------------------------------------------
// Patient descriptor
// ----------------------------------------------------------------------
type Patient struct {
  // [kg]
  weightKG  double

  //
  age       int
  sex       string

  //
  rocCFG    PatientROC
}


// ----------------------------------------------------------------------
//
func NewPatient() *Patient { return new(Patient) }

// ----------------------------------------------------------------------
//
func (p *Patient) setDefaults() {
  //
  p.age = -1
  p.sex = sexMale
  p.rocCFG = PatientROC_def()
}

// ----------------------------------------------------------------------
//
const (
  //
  drugRocuronium = "Roc"
  drugCisatracurium = "CisAtra"

  //
  sexMale = "male"
  sexFemale = "female"
)

// ----------------------------------------------------------------------
//
func (p *Patient) Vc_roc() volume { return volume{ 38.0 * p.weightKG, mL }}
