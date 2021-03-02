package main


type double = float64

// ----------------------------------------------------------------------
//
func conv(inval double, inunit int, outunit int) double {
  //
  for inunit != outunit {
    //
    if inunit < outunit {
      inval /= 1000.0; inunit++;
    } else {
      inval *= 1000.0; inunit--;
    }
  }

  //
  return inval
}

// ----------------------------------------------------------------------
//
type weight struct {
  //
  value double
  unit int
}

// ----------------------------------------------------------------------
//
const (
  kg= 1
  g = 0
  mg= -1
  ug= -2
  ng= -3
)

// ----------------------------------------------------------------------
//
type volume struct {
  //
  value double
  unit int
}

// ----------------------------------------------------------------------
//
const (
  L = 0
  mL= -1
  uL= -2
  nL= -3
)

// ----------------------------------------------------------------------
//
type conc struct {
  //
}

// ----------------------------------------------------------------------
//
func _gg(val double) weight { return weight{val, 0} }
func _mg(val double) weight { return weight{val, -1} }
func _ug(val double) weight { return weight{val, -2} }
func _ng(val double) weight { return weight{val, -3} }

// ----------------------------------------------------------------------
//
func (w weight) in(outunit int) weight {
  //
  return weight{ conv(w.value, w.unit, outunit), outunit }
}

// ----------------------------------------------------------------------
//
func (v volume) in(outunit int) volume {
  //
  return volume{ conv(v.value, v.unit, outunit), outunit }
}

//
func volume_0() volume { return volume{0, mL} }
func weight_0() weight { return weight{0, mL} }
