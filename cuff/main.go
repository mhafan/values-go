package main

//
import (
	"flag"
	"rcore"
)

// ----------------------------------------------------------------------
// CUF Behavior:
// 1) read TOF values set by PatMod
// 2) add some noise if configured
func cuffMain() {
	// next in the loop: TCM
	rcore.CurrentExp.Say(rcore.CallTCM)
}

// ----------------------------------------------------------------------
//
func main() {
	// --------------------------------------------------------------------
	//
	flag.Parse()

	ent := rcore.Entity{rcore.CallSensor, true, nil,
		cuffMain,
		func() {}, func() {}, func() {}}

	// --------------------------------------------------------------------
	// listening on channels vm.*
	// - standard behavior on start/end
	// - SENSOR -> my action
	rcore.EntityCore(&ent)
}
