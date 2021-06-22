// --------------------------------------------------------------------
// Entity = one component in the distributed system Values-NMT-HIL
// --------------------------------------------------------------------
package rcore

// --------------------------------------------------------------------
// ...
import (
	"log"
	"os"
)

// ----------------------------------------------------------------------
// Entity descriptor:
// Entities: CNT, PM, PUMP, ...
// ----------------------------------------------------------------------
type Entity interface {
	// ------------------------------------------------------------------
	// activating on this message in the context of HiL platform
	MyTurn() string

	// ------------------------------------------------------------------
	//
	ResetState()

	// ------------------------------------------------------------------
	// Behavior in the HiL platform protocol
	// start command
	// end command
	// iteration (MyTurn) command
	// default function - for all inputs except for start, end and MyTurn
	CycleFunction()
	StartFunction()
	EndFunction()
	DefaultFunction(msg Rmsg)
}

// ----------------------------------------------------------------------
//
func EntityStartSequence(ent Entity, expIDChannel string) {
	//
	log.Println("New experiment started: ", expIDChannel)

	//
	CurrentExp = MakeExpID(expIDChannel)

	//
	CurrentExp.LoadAll()

	//
	ent.StartFunction()
}

// ----------------------------------------------------------------------
//
func EntityEndSequence(ent Entity, expIDChannel string) {
	////
	if CurrentExp != nil {
		//
		if CurrentExp.ischannel(expIDChannel) {
			//
			log.Println("Experiment ended: ", expIDChannel)

			//
			ent.EndFunction()
		}
	}

	//
	CurrentExp = nil
}

// ----------------------------------------------------------------------
//
func EntityRoundSequence(ent Entity, expIDChannel string) {
	//
	if CurrentExp != nil {
		//
		if CurrentExp.ischannel(expIDChannel) {
			//
			ent.CycleFunction()
		}
	}
}

// ----------------------------------------------------------------------
//
func EntityMasterChannel(msg Rmsg) {
	//
	switch msg.Message {
	case "quit":
		//
		Global.Running = false
		//
		return
	}
}

// ----------------------------------------------------------------------
// Main function of the entity's life cycle
func EntityCore(ent Entity) {
	// --------------------------------------------------------------------
	// initiate the r-sysem library (sender|listener)
	//
	_rglobal := RServerInit()

	// some error
	if _rglobal == nil {
		//
		log.Println("R-system library start failure")

		//
		os.Exit(1)
	}

	// --------------------------------------------------------------------
	// become a new follower (receiver of messages from vm.*)
	_meFollower := NewFollower()

	// --------------------------------------------------------------------
	// Entity's runLoop
	for Global.Running == true {
		// waiting for an input from Listener
		msg := <-_meFollower.Inputs

		//
		if msg.Channel == MasterChannel {
			//
			EntityMasterChannel(msg)
			//
		} else {
			//
			switch msg.Message {
			//
			case CallStart:
				//
				EntityStartSequence(ent, msg.Channel)

			//
			case CallEnd:
				//
				EntityEndSequence(ent, msg.Channel)

				//
			case ent.MyTurn():
				//
				EntityRoundSequence(ent, msg.Channel)

			default:
				//
				if CurrentExp != nil {
					//
					if CurrentExp.ischannel(msg.Channel) {
						//
						ent.DefaultFunction(msg)
					}
				}
			}
		}
	}
}
