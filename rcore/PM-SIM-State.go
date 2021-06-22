package rcore

// ----------------------------------------------------------------------
// 3 compartment model
type COMP_X [3 + 1]Double

// ----------------------------------------------------------------------
// Simulation Predicate
// SIMS -> bool
type SimPredicate = func(sim *SIMS) bool

// ----------------------------------------------------------------------
// Simulation state:
// ----------------------------------------------------------------------
type SIMS struct {
	// assuming this drug CONST
	Drug DrugDef

	// --------------------------------------------------------------------
	// - Weight [kg]
	// - Volume of Distribution in Central Compartment (const)
	// - hill coefs
	Weight    Weight
	VdCentral Volume
	HillCoefs Hill

	// --------------------------------------------------------------------
	// simulation internal data
	Time int

	// --------------------------------------------------------------------
	// Continuous Simulator state variables (integrators)
	YROC     COMP_X
	CinpStat MStat

	//
	Effect LinScale

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

	// default drug driver (will be replaced later on)
	// default weight and others
	out.Drug = Rocuronium{}
	out.Weight = Weight{0, Kg}
	out.VdCentral = Volume_0()
	out.Bolus = Volume_0()
	out.Infusion = Volume_0()

	//
	out.BolusConsumptionML = 0
	out.CinpStat.Reset()

	//
	out.YROC = COMP_X{0, 0, 0, 0}
	out.Effect = LinScale{}
	out.HillCoefs = out.Drug.DefHillCoefs()

	//
	return &out
}

// ----------------------------------------------------------------------
//
func (r *SIMS) Cinp() Double {
	//
	return r.YROC[1]
}

// ----------------------------------------------------------------------
//
func (r *SIMS) Clone() *SIMS {
	//
	out := *r

	//
	return &out
}

// ----------------------------------------------------------------------
// Initial construction of SIMS from Exprec
func (r *SIMS) SetupFrom(e *Exprec) {
	//
	r.Drug = MakeDrugDef(e.Drug)
	r.Weight = Weight{Double(e.Weight), Kg}

	//
	r.VdCentral = r.Drug.DefVd(Double(e.AbsoluteVd), Double(e.UnitVd), r.Weight)
}

// ----------------------------------------------------------------------
// rcore.Exprec transger to SIMS (patmod internal data struct)
func (r *SIMS) UpdateFrom(e *Exprec) {
	//
	r.Bolus = Volume{e.Bolus, ML}
	r.Infusion = Volume{e.Infusion, ML}

	//
	r.BolusConsumptionML = 0
	r.CinpStat.Reset()

	//
	if e.EC50 > 0 {
		//
		r.HillCoefs.EC50 = e.EC50
	}
}

// ----------------------------------------------------------------------
// next state by shifting time
func (from *SIMS) NextState(at int) *SIMS {
	// ... copy ...
	ns := from.Clone()

	// shift time
	ns.Time = at

	// reset inputs
	ns.Bolus = Volume_0()
	ns.Infusion = Volume_0()

	//
	return ns
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
		from.Drug.SimStep(from)

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

// --------------------------------------------------------------------
//
func (from *SIMS) SimStepsWhile(pred SimPredicate) *SIMS {
	// reach Mtime in 1s simulation steps
	for pred(from) {
		// h = 1s, continuous simulation step
		from.Drug.SimStep(from)

		//
		from = from.Next1S()
	}

	//
	return from
}

// ----------------------------------------------------------------------
//
func (sims *SIMS) CinpAboveLow(expID *Exprec) bool {
	//
	return sims.Cinp() >= expID.TargetCinpLow
}
