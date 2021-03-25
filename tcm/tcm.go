//
package main

// --------------------------------------------------------------------
// ...
import (
	"flag"
	"log"
	"os"
	"rcore"
	"time"
)

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

	// name of TEST-CASE must be given
	if *flTC == "" {
		// print help and exit
		flag.PrintDefaults()

		//
		return
	}

	// --------------------------------------------------------------------
	// initiate the r-sysem library (sender|listener)
	_rglobal := rcore.RServerInit()

	// some errror
	if _rglobal == nil {
		//
		log.Println("R-system library start failure")
	}

	// --------------------------------------------------------------------
	// become a new follower (receiver of messages from vm.*)
	_meFollower := rcore.NewFollower()

	// --------------------------------------------------------------------
	// assign a new simulation experiment name
	// format: vm.(flTC).increasedCounter
	_expID := rcore.NewExpID(*flTC)

	// create a REDIS record for that name
	rcore.CurrentExp = rcore.MakeExpID(_expID)

	// fill the record with initial data (prompt etc)
	mydefs(rcore.CurrentExp)

	// REDIS save, publish first msg -> START
	rcore.CurrentExp.Save([]string{}, true)
	// START !!!
	// - all participants (CNT, PUMP, PM, SENSOR) are supposed to
	// get ready for this experiment
	rcore.CurrentExp.Say(rcore.CallStart)

	//
	log.Println("ExperimentID=", _expID, "; starting")

	// --------------------------------------------------------------------
	// The experiment goes in cycles with a predefined number of iterations
	for {
		// if the experiment was interrupted (...)
		if rcore.Global.Running == false {
			//
			break
		}

		// ------------------------------------------------------------------
		// ending condition
		if rcore.CurrentExp.Cycle > *flCycles || rcore.CurrentExp.Mtime > *flTMAX {
			//
			break
		}

		//
		log.Println("Cycle: ", rcore.CurrentExp.Cycle)

		// ------------------------------------------------------------------
		// next cycle, save the record and call CNT out
		rcore.CurrentExp.Save([]string{"cycle", "mtime"}, false)
		rcore.CurrentExp.Say(rcore.CallCNT)

		// ------------------------------------------------------------------
		//
		var _waiting = true

		// ------------------------------------------------------------------
		// waiting for the loop to go around
		// CNT -> PUMP -> PM -> CUFF -> TCM
		// ------------------------------------------------------------------
		// TCM passed the token to CNT. Now TCM needs to wait
		// until the loop gets back to TCM
		for _waiting == true {
			//
			select {
			// some input message, now check the channel
			case msg := <-_meFollower.Inputs:
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
					// TCM,
					if msg.Message == rcore.CallTCM {
						// end of waiting...
						_waiting = false
						//
						break
					}
				}

			// timeout
			case <-time.After(time.Second * 2):
				//
				log.Println("Timeout. Ending")

				//
				os.Exit(1)
			}
		}

		// ------------------------------------------------------------------
		// next cycle, next time moment
		rcore.CurrentExp.Cycle++
		rcore.CurrentExp.Mtime += *flTimeStep
	}

	// --------------------------------------------------------------------
	// publish end
	rcore.RPublish(_expID, rcore.CallEnd)
}
