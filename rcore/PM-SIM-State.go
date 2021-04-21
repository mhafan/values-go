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
	//
	Drug string
	// --------------------------------------------------------------------
	// Weight [kg]
	// Volume of Distribution in Central Compartment (const)
	Weight    Weight
	VdCentral Volume

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
func EmptySIMS() *SIMS {
	//
	out := SIMS{}

	//
	out.Drug = DrugRocuronium
	out.Weight = Weight{0, Kg}
	out.VdCentral = Volume_0()
	out.Rocs = ROCS{COMP_X{}, RocDefHill(), 0, 0}
	out.Bolus = Volume_0()
	out.Infusion = Volume_0()

	//
	return &out
}

// ----------------------------------------------------------------------
// Initial construction of SIMS from Exprec
func (r *SIMS) SetupFrom(e *Exprec) {
	//
	r.Drug = e.Drug
	r.Weight = Weight{float64(e.Weight), Kg}
	r.VdCentral = VdFor(r.Drug, float64(e.AbsoluteVd), float64(e.UnitVd), r.Weight)
}

// ----------------------------------------------------------------------
// rcore.Exprec transger to SIMS (patmod internal data struct)
func (r *SIMS) UpdateFrom(e *Exprec) {
	//
	r.Bolus = Volume{Double(e.Bolus), ML}
	r.Infusion = Volume_0()

	//
	if e.EC50 > 0 {
		//
		r.Rocs.rocHill.ec50 = e.EC50
	}
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
