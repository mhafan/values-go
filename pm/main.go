package main

//
import "fmt"
import "flag"

// ----------------------------------------------------------------------
//
var flServer = flag.Bool("s", false, "Run in server mode")


// ----------------------------------------------------------------------
//
func main() {
  //
  fmt.Println("patmod in golang")

  //
  flag.Parse()

  //
  if *flServer == true {
    //
    rserverMain(); return;
  }

  //
  patient := NewPatient()

  //
  patient.setDefaults()
  patient.weightKG = 100

  //
  ss0 := emptySIMS(patient)
  ss0.bolus = volume{20, mL}

  //
  var ns SIMS = ss0

  for ns.time < 100 {
    //
    ns.rocSimStep()

    //
    fmt.Println(ns.time, " ", ns.rocs.yROC, " ", ns.rocs.TOF0)

    //
    ns = ns.next1S()
  }
}
