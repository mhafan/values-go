// ----------------------------------------------------------------------
// PatMod for R-system
// ----------------------------------------------------------------------
package main

//
import "rcore"
import "log"

//
var rpatient *Patient = nil
var rsims *SIMS = nil

// ----------------------------------------------------------------------
// rcore.Exprec transger to SIMS (patmod internal data struct)
func (r *SIMS) updateFrom(e *rcore.Exprec) {
  // patient info
  r.patient.weightKG = rcore.Double(e.Weight)
  r.patient.age = e.Age
  r.patient.sex = sexMale
  // TODO: Vd_kg, ec50
  // r.patient.rocCFG

  //
  r.time = e.Mtime

  //
  r.bolus = rcore.Volume { rcore.Double(e.Bolus), rcore.ML }
  r.infusion = rcore.Volume_0()
}

// ----------------------------------------------------------------------
// For every cycle of distributed simulation.
func rserverCycle() {
  //
  tol := []string { "mtime", "cycle", "bolus", "infusion" }

  // --------------------------------------------------------------------
  // load Redis record
  if rcore.EntityExpRecReload(tol) == false {
    //
    panic("REDIS record not found")
  }

  // ...
  var _c = rcore.CurrentExp

  // --------------------------------------------------------------------
  // REDIS record -> SIMS (patmod simulation state)
  rsims.updateFrom(_c)

  //
  log.Println("PM; cycle: ", _c.Cycle, "mtime", _c.Mtime,
      "rtime", rsims.time, "bolus ", _c.Bolus, " ", rsims.bolus)

  // --------------------------------------------------------------------
  // reach Mtime in 1s simulation steps
  for rsims.time <= _c.Mtime {
    // h = 1s, continuous simulation step
    rsims.rocSimStep()

    //
    log.Println(rsims.time, " ", rsims.rocs.yROC, " ", rsims.rocs.TOF0)

    // +1s
    rsims = rsims.next1S();

    //
    if rsims.time == _c.Mtime { break }
  }

  //
  _c.Say(rcore.CallSensor)
}

// ----------------------------------------------------------------------
//
func rserverStart() {
  //
  rpatient = NewPatient()

  //
  rpatient.setDefaults()
  rpatient.weightKG = rcore.Double(rcore.CurrentExp.Weight)

  //
  rsims = emptySIMS(rpatient)

  //
  log.Println("Starting new patient: ", rsims)
}

// ----------------------------------------------------------------------
//
func rserverEnd() {
  //
  rpatient = nil
  rsims = nil

  //
  log.Println("Ending the experiment")
}


// ----------------------------------------------------------------------
// main function (called from main() when arg is -s)
func rserverMain() {
  //
  rcore.EntityCore(rcore.CallPatMod, rserverCycle, rserverStart, rserverEnd);
}
