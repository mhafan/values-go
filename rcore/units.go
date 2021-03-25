package rcore


type Double = float64

// ----------------------------------------------------------------------
//
func conv(inval Double, inunit int, outunit int) Double {
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
type Weight struct {
  //
  Value Double
  Unit int
}

// ----------------------------------------------------------------------
//
const (
  Kg= 1
  G = 0
  Mg= -1
  Ug= -2
  Ng= -3
)

// ----------------------------------------------------------------------
//
type Volume struct {
  //
  Value Double
  Unit int
}

// ----------------------------------------------------------------------
//
const (
  L = 0
  ML= -1
  UL= -2
  NL= -3
)


// ----------------------------------------------------------------------
//
func _gg(val Double) Weight { return Weight{val, 0} }
func _mg(val Double) Weight { return Weight{val, -1} }
func _ug(val Double) Weight { return Weight{val, -2} }
func _ng(val Double) Weight { return Weight{val, -3} }

// ----------------------------------------------------------------------
//
func (w Weight) In(outunit int) Weight {
  //
  return Weight{ conv(w.Value, w.Unit, outunit), outunit }
}

// ----------------------------------------------------------------------
//
func (v Volume) In(outunit int) Volume {
  //
  return Volume{ conv(v.Value, v.Unit, outunit), outunit }
}

//
func Volume_0() Volume { return Volume{0, ML} }
func Weight_0() Weight { return Weight{0, ML} }


// ----------------------------------------------------------------------
//
func RocWSOL(w Weight) Volume {
  //
  return Volume{ w.In(Mg).Value / 10.0, ML }
}

//
func RocSOLW(v Volume) Weight {
  //
  return Weight{ v.In(ML).Value * 10.0, Mg }
}
