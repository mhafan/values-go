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
var flTC = flag.String("C", "", "Name of test-case")
var flSex = flag.String("s", rcore.SexMale, "Sex of patient {male/female}")
var flAge = flag.Int("A", 42, "Age of patient")
var flWeight = flag.Int("w", 100, "Weight [kg] of patient")

//
var flTimeStep = flag.Int("t", 15, "Time step [s]")
var flTMAX = flag.Int("T", 1000000, "Max Time [s]")
var flCycles = flag.Int("c", 100000, "Number of cycles")


// --------------------------------------------------------------------
//
func mydefs(_c *rcore.Exprec) {
  //
  _c.Weight = *flWeight
  _c.Age = *flAge

  //
  _c.Drug = rcore.DrugRocuronium
}

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
	// assign a new simulation experiment name
	_expID := rcore.NewExpID(*flTC)

	// create a REDIS record for that name
	rcore.CurrentExp = rcore.MakeExpID(_expID)

  // fill the record with initial data (prompt etc)
  mydefs(rcore.CurrentExp)

  // REDIS save, publish first msg -> START
	rcore.CurrentExp.Save([]string{}, true)
	rcore.CurrentExp.Say(rcore.CallStart)

	//
	log.Println("ExperimentID=", _expID, "; starting")

  // --------------------------------------------------------------------
  //
	for {
    //
    if rcore.Global.Running == false {
      //
      break
    }
    // ------------------------------------------------------------------
    //
    var _r *rcore.Exprec = rcore.CurrentExp
    var _waiting = true

		// ending condition
		if _r.Cycle > *flCycles || _r.Mtime > *flTMAX {
			//
			break;
		}

    //
    log.Println("Cycle: ", rcore.CurrentExp.Cycle)

    // ------------------------------------------------------------------
		// next cycle, save the record and call CNT out
		rcore.CurrentExp.Save([]string{ "cycle", "mtime"}, false)
		rcore.RPublish(_expID, rcore.CallCNT)

    // ------------------------------------------------------------------
		// waiting for the loop to go around
    // CNT -> PUMP -> PM -> CUFF -> TCM
		for _waiting == true {
			//
			select {
				// input messages
				case msg := <- _meFollower.Inputs:
					// calling me, ...
          if msg.Channel == rcore.MasterChannel {
            //
            rcore.EntityMasterChannel(msg)

            //
            if rcore.Global.Running == false {
              //
              break
            }
            //
          } else {
            //
            if msg.Message == rcore.CallTCM {
              //
              _waiting = false; break
            }
          }

        // timeout
        case <- time.After(time.Second * 2):
					//
					log.Println("Timeout. Ending");

					//
					os.Exit(1)
			}
		}

    // ------------------------------------------------------------------
		// next cycle, next time moment
		rcore.CurrentExp.Cycle++;
		rcore.CurrentExp.Mtime += *flTimeStep
	}

	// --------------------------------------------------------------------
	// publish end
	rcore.RPublish(_expID, rcore.CallEnd)
}
