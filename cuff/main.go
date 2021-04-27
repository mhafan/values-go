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
	defer rcore.CurrentExp.Say(rcore.CallTCM)
}

// ----------------------------------------------------------------------
//
func main() {
	// --------------------------------------------------------------------
	//
	flag.Parse()

	//
	ent := rcore.MakeNewEntity()

	//
	ent.MyTurn = rcore.CallSensor
	ent.What = cuffMain

	// --------------------------------------------------------------------
	// listening on channels vm.*
	// - standard behavior on start/end
	// - SENSOR -> my action
	rcore.EntityCore(ent)
}
