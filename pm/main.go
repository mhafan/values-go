package main

//
import "fmt"


// ----------------------------------------------------------------------
//
func main() {
  //
  fmt.Println("patmod in golang")

  //
  patient := &Patient{ 100 }

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
