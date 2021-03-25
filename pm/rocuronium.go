package main

import "rcore"

//
func rocDefHill() Hill {
  //
  return Hill { 100.0, 0.823, 4.79 }
}

// ----------------------------------------------------------------------
//
const (
  rocK12 = 0.259/60.0
  rocK21 = 0.163/60.0
  rocK13 = 0.060/60.0
  rocK31 = 0.012/60.0
  rocK10 = 0.119/60.0
)


// ----------------------------------------------------------------------
//
func rocInputs(y COMP_X, infConc rcore.Double) COMP_X {
  //
  var out COMP_X

  //
  out[1] = y[2] * rocK21 + y[3] * rocK31 - y[1] * (rocK10 + rocK12 + rocK13)
  out[2] = y[1] * rocK12 - y[2] * rocK21
  out[3] = y[1] * rocK13 - y[3] * rocK31

  //
  out[1] += infConc

  //
  return out
}

// ----------------------------------------------------------------------
//
func rocSimStep1H(yin COMP_X, infConc rcore.Double) COMP_X {
  //
  f := rocInputs(yin, infConc)

  //
  var out COMP_X

  //
  for i := 1; i < 4; i++ {
    //
    out[i] = yin[i] + f[i];
  }

  //
  return out
}


// ----------------------------------------------------------------------
// PK/PD Model for Rocuronium
// ----------------------------------------------------------------------
func (ss *SIMS) rocSimStep() {
  // Volume of distribution
  Vd := ss.patient.Vc_roc()

  // eventual infusion input
  ic := 0.0

  // [ml/hr] => effective weight / hr => per s =>
  if ss.infusion.Value > 0.0 {
    // [ug/hr]
    inWeightHour := rcore.RocSOLW(ss.infusion).In(rcore.Ug)

    // [ug/s]
    inWeightS := inWeightHour.Value / 3600.0

    //
    ic = inWeightS / Vd.In(rcore.ML).Value
  }

  //
  if ss.time > 0 {
    //
    ss.rocs.yROC = rocSimStep1H(ss.rocs.yROC, ic)
  }

  //
  if ss.bolus.Value > 0.0 {
    //
    ss.rocs.yROC[1] += rcore.RocSOLW(ss.bolus).In(rcore.Ug).Value / Vd.In(rcore.ML).Value;
  }

  //
  ss.rocs.effect = ss.rocs.rocHill.value(ss.rocs.yROC[1])
  ss.rocs.TOF0 = int(_TOFbounds.bound(100.0 - ss.rocs.effect))
}
