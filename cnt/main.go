package main

//
import "log"
import "flag"
import "rcore"

// ----------------------------------------------------------------------
// direct dosing of NMT blockator:
// time = 0: initial bolus
// then in "flBolusInterval" intervals, bolus "flBolusAmount"
var flBolusInterval = flag.Int("b", 200, "Interval between boluses [s]")
var flBolusAmount = flag.Int("B", 5, "Bolus volume [mL of solution]")

// ----------------------------------------------------------------------
// status
var _lastTimeBolus = 0
var _scheduledBolusAt = -1
var _initialBolusGiven = false


// ----------------------------------------------------------------------
// Initil bolus as defined by manufacturer
func initialBolus(drug string, wkg int) rcore.Volume {
  //
  switch drug {
  case rcore.DrugRocuronium:
    // 0.6 mg per [kg] of patient's weight
    return rcore.RocWSOL(rcore.Weight{ 0.6 * 100.0, rcore.Mg }).In(rcore.ML)
  case rcore.DrugCiatracurium:
    // TODO
    return rcore.Volume { 0, rcore.ML }
  }

  // default value if drug is set incorrectly
  return rcore.Volume { 0, rcore.ML }
}

// ----------------------------------------------------------------------
// with every START msg, do reset internals
func startupWithExperiment() {
  //
  _lastTimeBolus = 0
  _scheduledBolusAt = -1
  _initialBolusGiven = false

  //
  if *flBolusInterval > 0 {
    //
    _scheduledBolusAt = *flBolusInterval
  }
}

// ----------------------------------------------------------------------
// Direct MODE:
// time == 0 => INITIAL bolus
// intervals
func regulationInDirectMode(_r *rcore.Exprec) bool {
  // by defualt, set both zero
  _r.Bolus = 0
  _r.Infusion = 0

  // --------------------------------------------------------------------
  // time == 0 || Cycle == 0
  // --------------------------------------------------------------------
  if _r.Mtime <= 0 || _r.Cycle == 0 {
    //
    if _initialBolusGiven == true {
      //
      log.Println("Will not set the initial bolus multiple times!");

      // error
      return false
    }

    // initial bolus in recommended volume 0.6mg/kg
    _r.Bolus = int(initialBolus(_r.Drug, _r.Weight).Value)
    _initialBolusGiven = true

    //
    log.Println("CNT:initial bolus [mL]: ", _r.Bolus)

    //
    return true
  }

  // --------------------------------------------------------------------
  // repetitive bolus, if enabled:
  // the time has reached scheduled moment
  if _r.Mtime >= _scheduledBolusAt && _scheduledBolusAt > 0 {
    // now
    _lastTimeBolus = _r.Mtime
    // schedule the next moment
    _scheduledBolusAt = _r.Mtime + (*flBolusInterval)
    // set the bolus
    _r.Bolus = *flBolusAmount

    //
    log.Println("CNT:repetitive bolus [mL]:", _r.Bolus, " time=", _r.Mtime)

    //
    return true
  }

  //
  return true
}

// ----------------------------------------------------------------------
// Regulation cycle =>
// 1) direct mode
// 2) feedback mode
func cycle() {
  // --------------------------------------------------------------------
  //
  if rcore.EntityExpRecReload([]string { "cycle", "mtime" }) == false {
    //
    return
  }

  // --------------------------------------------------------------------
  //
  var _r = rcore.CurrentExp

  //
  if regulationInDirectMode(_r) == true {
    //
    _r.Save([]string{ "bolus", "infusion" }, false)
  }
}

// ----------------------------------------------------------------------
//
func main() {
  //
  flag.Parse()

  //
  rcore.EntityCore(rcore.CallCNT, func() {
    //
    cycle()

    //
    rcore.CurrentExp.Say(rcore.CallPump)
    //
  }, startupWithExperiment, func() {})
}
