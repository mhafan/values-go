//
package main

// --------------------------------------------------------------------
// ...
import "log"
import "rcore"
import "flag"
import "time"
import "os"

// ----------------------------------------------------------------------
//
var flTC = flag.String("t", "", "Name of test-case")

//
var flTimeStep = flag.Int("T", 15, "Time step [s]")
var flCycles = flag.Int("c", 100, "Number of cycles")

// --------------------------------------------------------------------
//
func main() {
	//
  flag.Parse()

	//
	if *flTC == "" {
		//
		flag.PrintDefaults(); return
	}

  // --------------------------------------------------------------------
  // initiate the r-sysem library (sender|listener)
  _rglobal := rcore.RServerInit()

  // some errror
  if _rglobal == nil {
    //
    log.Println("R-system library start failure");
  }

  // --------------------------------------------------------------------
  // become a new follower (receiver of messages from vm.*)
  _meFollower := rcore.NewFollower()

	// --------------------------------------------------------------------
	//
	_expID := rcore.NewExpID(*flTC)

	//
	rcore.CurrentExp = rcore.MakeExpID(_expID)
	rcore.CurrentExp.Save()
	rcore.CurrentExp.Say(rcore.CallStart)

	//
	log.Println("ExperimentID=", _expID, "; starting")

  // --------------------------------------------------------------------
  //
	for {
		//
		if rcore.CurrentExp.Cycle >= *flCycles {
			//
			break;
		}

		//
		rcore.CurrentExp.Save()
		rcore.RPublish(_expID, rcore.CallCNT)

		//
		for _waiting := true; _waiting == true; {
			//
			select {
				//
				case msg := <- _meFollower.Inputs:
					//
					if msg.Message == rcore.CallTCM {
						//
						log.Println("TCM; going to the next cycle");

						_waiting = false; break
					}

				case <- time.After(time.Second * 1):
					//
					log.Println("Timeout. Ending");

					//
					os.Exit(1)
			}
		}

		//
		rcore.CurrentExp.Cycle++;
		rcore.CurrentExp.Mtime += *flTimeStep
	}

	// --------------------------------------------------------------------
	//
	rcore.RPublish(_expID, rcore.CallEnd)
}
