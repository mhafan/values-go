package main

//
import "fmt"
import "os"
import "flag"
import "rcore"

// ----------------------------------------------------------------------
//
var flServer = flag.Bool("s", false, "Run in server mode")

// ----------------------------------------------------------------------
//
func main() {
  //
  flag.Parse()

  //
  if *flServer == false {
    //
    fmt.Println("run me in -s mode, por favor...")

    //
    os.Exit(2)
  }

  //
  rcore.EntityCore(rcore.CallCNT, func() {
    //
    rcore.CurrentExp.Say(rcore.CallPump)
  })
}
