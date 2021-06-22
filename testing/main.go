package main

import (
	"fmt"
	"rcore"
)

//
func mainaaa() {
	//
	drug := rcore.Rocuronium{}
	hills := drug.DefHill4Coefs()

	//
	for c := 0.0; c < 5; c += 0.1 {
		//
		out := rcore.LinScaleModel4(c, hills)

		//
		fmt.Println(c, out.TOFcount(), out.TOFratio(), out.TOF4Amplitude[0], out.TOF4Amplitude[1], out.TOF4Amplitude[2], out.TOF4Amplitude[3])
	}
}

// ----------------------------------------------------------------------
//
func totrevpred(sims *rcore.SIMS) bool {
	//
	return sims.Effect.TOFSimpleAmplitude() < 90
}

//
func main() {
	//
	drug := rcore.Rocuronium{}

	//
	ss := rcore.EmptySIMS()
	ss.Drug = drug
	ss.Weight = rcore.Weight{100, rcore.Kg}
	ss.VdCentral = drug.DefVd(0, 0, ss.Weight)
	ss.Bolus = drug.InitialBolus(100, 1, 0.6)

	//
	fmt.Println(ss.Bolus.Value)

	//
	for tm := 0; tm < 30*60; tm++ {
		//
		ss.Drug.SimStep(ss)

		//
		trec := ss.Clone().SimStepsWhile(totrevpred)

		//
		fmt.Println(tm/60, tm, ss.Effect.TOFSimpleAmplitude(), ss.Cinp(), "totre", trec.Time-tm)

		//
		ss = ss.Next1S()
	}
}
