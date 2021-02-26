package main

//
type Patient struct {
  //
  weightKG double
}

//
const (
  drugRocuronium = "Roc"
  drugCisatracurium = "CisAtra"
)

//
func (p *Patient) Vc_roc() volume { return volume{ 38.0 * p.weightKG, mL }}
