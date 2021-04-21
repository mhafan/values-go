package rcore

// ----------------------------------------------------------------------
// 3 compartment model
type COMP_X [3 + 1]Double

//
var _TOFbounds = Bounds{0, 100}

// ----------------------------------------------------------------------
//
type ROCS struct {
	// --------------------------------------------------------------------
	// model, internal
	yROC    COMP_X
	rocHill Hill

	//
	effect Double
	TOF0   int
}

// ----------------------------------------------------------------------
// Simulation state:
// ----------------------------------------------------------------------
type SIMS struct {
	// --------------------------------------------------------------------
	//
	Patient *Patient

	// --------------------------------------------------------------------
	// simulation internal data
	Time int

	//
	Rocs ROCS

	// --------------------------------------------------------------------
	// inputs
	Bolus    Volume
	Infusion Volume
}

// ----------------------------------------------------------------------
// time(0) zero simstate constructor
func EmptySIMS(pat *Patient) *SIMS {
	//
	return &SIMS{
		//
		pat, 0,
		//
		ROCS{COMP_X{}, RocDefHill(), 0, 0},
		// bolus & infusion
		Volume_0(), Volume_0()}
}

// ----------------------------------------------------------------------
// next state by shifting time
func (from *SIMS) NextState(at int) *SIMS {
	// ... copy ...
	ns := *from

	// shift time
	ns.Time = at
	// reset inputs
	ns.Bolus = Volume_0()
	ns.Infusion = Volume_0()

	//
	return &ns
}

// ----------------------------------------------------------------------
//
func (from *SIMS) Next1S() *SIMS {
	//
	return from.NextState(from.Time + 1)
}

// ----------------------------------------------------------------------
// rcore.Exprec transger to SIMS (patmod internal data struct)
func (r *SIMS) UpdateFrom(e *Exprec) {
	// patient info
	r.Patient.WeightKG = Double(e.Weight)
	r.Patient.Age = e.Age
	r.Patient.Sex = SexMale
	// TODO: Vd_kg, ec50
	// r.patient.rocCFG

	//
	r.Bolus = Volume{Double(e.Bolus), ML}
	r.Infusion = Volume_0()
}
