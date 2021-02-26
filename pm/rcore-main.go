// ----------------------------------------------------------------------
// PatMod for R-system
// ----------------------------------------------------------------------
package main

//
import "rcore"
import "log"

// ----------------------------------------------------------------------
//
func rserverCycle() {
  //
  log.Println("PM; cycle")

  //
  rcore.CurrentExp.Say(rcore.CallSensor)
}

// ----------------------------------------------------------------------
// main function (called from main() when arg is -s)
func rserverMain() {
  // --------------------------------------------------------------------
  // initiate the r-sysem library (sender|listener)
  _rglobal := rcore.RServerInit()

  // some errror
  if _rglobal == nil {
    //
    log.Println("PM; R-system library start failure");
  }

  // --------------------------------------------------------------------
  // become a new follower (receiver of messages from vm.*)
  _meFollower := rcore.NewFollower()


  // --------------------------------------------------------------------
  //
	for {
		//
		msg := <- _meFollower.Inputs

    //
    switch msg.Message {
      //
      case rcore.CallStart:
        //
        log.Println("New experiment started: ", msg.Channel)
        rcore.CurrentExp = rcore.MakeExpID(msg.Channel)

      //
      case rcore.CallEnd:
        //
        if rcore.CurrentExp != nil {
          //
          if rcore.CurrentExp.Varname == msg.Channel {
            //
            log.Println("Experiment ended: ", msg.Channel)
            rcore.CurrentExp = nil
          }
        }

      //
      case rcore.CallPatMod:
        //
        if rcore.CurrentExp != nil { rserverCycle() }

      default:
        ;
    }
	}
}
