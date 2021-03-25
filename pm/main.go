// ----------------------------------------------------------------------
// Patient Model (to be integrated in RServer)
// Pharmacodynamic/Pharmacokinetic model for Rocuronium & Cisatracurium
// ----------------------------------------------------------------------

package main

//
import (
	"flag"
	"log"
	"rcore"
)

// ----------------------------------------------------------------------
//
var flServer = flag.Bool("s", false, "Run in server mode")
var flExpID = flag.String("E", "", "Experiment ID")

// ----------------------------------------------------------------------
//
func simulateRserver(expID string) {
	//
	log.Println("ExpID ", expID, " starting sequence")

	//
	rcore.EntityStartSequence(expID, rserverStart)

	//
	var _r = rcore.CurrentExp

	//
	_r.Mtime = 0
	_r.Cycle = 0

	//
	rcore.EntityEndSequence(expID, rserverEnd)
}

// ----------------------------------------------------------------------
//
func main() {
	// --------------------------------------------------------------------
	//
	flag.Parse()

	// --------------------------------------------------------------------
	// Program works either in "server" mode
	if *flServer == true {
		//
		rserverMain()
		return
	}

	// --------------------------------------------------------------------
	//
	if *flExpID != "" {
		//
		simulateRserver(*flExpID)
	}
}
