package rcore

// ----------------------------------------------------------------------
// 3 compartment model
type COMP_X [3 + 1]Double

//
var _TOFbounds = Bounds{0, 100}

// ----------------------------------------------------------------------
// Simulation state:
// ----------------------------------------------------------------------
type SIMS struct {
	// assuming this drug CONST
	Drug string

	// --------------------------------------------------------------------
	// Weight [kg]
	// Volume of Distribution in Central Compartment (const)
	Weight    Weight
	VdCentral Volume

	// --------------------------------------------------------------------
	// simulation internal data
	Time int

	// --------------------------------------------------------------------
	// Continuous Simulator state variables (integrators)
	YROC    COMP_X
	RocHill Hill

	//
	Effect Double
	TOF0   int

	//
	BolusConsumptionML Double

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
	out.Bolus = Volume_0()
	out.Infusion = Volume_0()

	//
	out.BolusConsumptionML = 0

	//
	out.YROC = COMP_X{0, 0, 0, 0}
	out.Effect = 0
	out.TOF0 = 0
	out.RocHill = RocDefHill()

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
	r.Infusion = Volume{Double(e.Infusion), ML}

	r.BolusConsumptionML = 0

	//
	if e.EC50 > 0 {
		//
		r.RocHill.ec50 = e.EC50
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

// --------------------------------------------------------------------
//
func (from *SIMS) SimSteps(till int) *SIMS {
	// reach Mtime in 1s simulation steps
	for from.Time <= till {
		// h = 1s, continuous simulation step
		from.RocSimStep()

		// +1s
		if from.Time+1 <= till {
			//
			from = from.Next1S()
		} else {
			//
			break
		}
	}

	//
	return from
}
