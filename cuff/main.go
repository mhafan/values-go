package main

//
import "flag"
import "rcore"

// ----------------------------------------------------------------------
// CUF Behavior:
// 1) read TOF values set by PatMod
// 2) add some noise if configured
func cuffMain() {
  //
}

// ----------------------------------------------------------------------
//
func main() {
  // --------------------------------------------------------------------
  //
  flag.Parse()

  // --------------------------------------------------------------------
  // listening on channels vm.*
  // - standard behavior on start/end
  // - SENSOR -> my action
  rcore.EntityCore(rcore.CallSensor, func() {
    // do your job
    cuffMain()

    // next in the loop: TCM
    rcore.CurrentExp.Say(rcore.CallTCM)
    //
  }, func() {}, func() {})
}
