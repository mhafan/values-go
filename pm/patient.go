package main

import "rcore"

// ----------------------------------------------------------------------
//
type PatientROC struct {
  // [mL/kg]
  Vd_kg     rcore.Double
  // [ug/mL]
  ec50      rcore.Double
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
  weightKG  rcore.Double

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
func (p *Patient) Vc_roc() rcore.Volume {
  //
  return rcore.Volume{ 38.0 * p.weightKG, rcore.ML }
}
