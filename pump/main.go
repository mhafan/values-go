package main

//
//import "log"
import (
	"flag"
	"rcore"
)

// ----------------------------------------------------------------------
//
var flBolusNoise = flag.Float64("b", 0.0, "Bolus noise")
var flInfusionNoise = flag.Float64("i", 0.0, "Infusion noise")

// ----------------------------------------------------------------------
//
func anyInfluence() bool {
	//
	return *flBolusNoise > 0 || *flInfusionNoise > 0
}

// ----------------------------------------------------------------------
// PUMP main:
// - add noise to bolus/infusion if set
func pumpMain() {
	/*
	  if anyInfluence() == true {
	    //
	    if rcore.EntityExpRecReload() == false {
	      //
	      return
	    }

	    // TODO:
	  }*/

	rcore.CurrentExp.Say(rcore.CallPatMod)
}

// ----------------------------------------------------------------------
//
func main() {
	//
	flag.Parse()

	ent := rcore.Entity{
		rcore.CallPump, true, nil,
		pumpMain,
		func() {}, func() {}, func() {}}

	//
	rcore.EntityCore(&ent)
}
