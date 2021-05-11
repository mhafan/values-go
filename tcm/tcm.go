//
package main

// --------------------------------------------------------------------
// ...
import (
	"flag"
	"fmt"
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

// --------------------------------------------------------------------
//
var flTimeStep = flag.Int("t", 15, "Time step [s]")
var flTMAX = flag.Int("T", 1000000, "Max Time [s]")
var flCycles = flag.Int("c", 100000, "Number of cycles")

// --------------------------------------------------------------------
// control args
var flCNT_strategy = flag.String("S", "basic", "CNT strategy { none, basic, fwsim }")
var flCNT_targetLow = flag.Float64("r", 2, "targetCinpLow")
var flCNT_forwardRange = flag.Int("f", 3*60, "fwsim fwRange [s] forward range")

// --------------------------------------------------------------------
//
func mydefs(_c *rcore.Exprec) {
	//
	_c.Weight = *flWeight
	_c.Age = *flAge

	//
	_c.Drug = rcore.DrugRocuronium
	_c.IbolusMg = 0.6
	_c.Wcoef = 1.0
	_c.CNTStrategy = *flCNT_strategy
	_c.UnitVd = 38
	_c.AbsoluteVd = 0

	//
	_c.RepeBolus = 2
	_c.RepeStep = 10 * 6000
	_c.FwRange = *flCNT_forwardRange

	//
	_c.TargetCinpLow = rcore.Double(*flCNT_targetLow)
	_c.TargetCinpHi = _c.TargetCinpLow + 3.0

	//
	_c.Bolus = 1
}

// --------------------------------------------------------------------
//
func printState(c *rcore.Exprec) {
	//
	fmt.Println(
		c.Cycle, "/", c.Mtime,
		" INPUTS ", c.Bolus, "/", c.Infusion,
		" OUTS ", c.Cinp,
		" TOF/PTC", c.TOF, "/", c.PTC,
		" cons ", c.ConsumedTotal,
		" rtime ", c.RecoveryTime)
}

// --------------------------------------------------------------------
//
func printCSV(file *os.File, c *rcore.Exprec) {
	//
	getc := c.CsvExport()

	file.WriteString(getc)
	file.WriteString("\n")

	//
	fmt.Println(getc)
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
	//
	csvfile, csvError := os.Create("out.csv")

	//
	if csvError != nil {
		//
		log.Println("Creating output CSV failed")

		//
		os.Exit(-2)
	}

	//
	defer csvfile.Close()

	//
	csvfile.WriteString(rcore.CsvHeader())

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

	//
	_toUpdateFrom := []string{
		"TOF", "PTC", "bolus", "infusion",
		"Cinp", "ConsumedTotal", "RecoveryTime"}

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
		//
		if rcore.EntityExpRecReload(_toUpdateFrom) == false {
			//
			break
		}

		//
		printState(rcore.CurrentExp)
		printCSV(csvfile, rcore.CurrentExp)

		// ------------------------------------------------------------------
		// next cycle, next time moment
		rcore.CurrentExp.Cycle++
		rcore.CurrentExp.Mtime += *flTimeStep
	}

	// --------------------------------------------------------------------
	// publish end
	rcore.RPublish(_expID, rcore.CallEnd)
}
