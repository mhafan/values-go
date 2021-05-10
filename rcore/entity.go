// --------------------------------------------------------------------
// Entity = one component in the distributed system Values-NMT-HIL
// --------------------------------------------------------------------
package rcore

// --------------------------------------------------------------------
// ...
import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// --------------------------------------------------------------------
//
var Global MConn = MConn{}

// --------------------------------------------------------------------
// hostname + AUTH passwd
var flHostname = flag.String("h", ":6379", "REDIS hostname")
var flAuth = flag.String("a", "", "REDIS auth password")

// --------------------------------------------------------------------
// make connection
func dial() (redis.Conn, error) {
	// open connection OPTIONS, passwd AUTH
	opta := redis.DialPassword(*flAuth)

	// open the connection with given optiions
	c, e := redis.Dial("tcp", *flHostname, opta)

	// ...
	if e != nil {
		//
		log.Println("Dial error: ", e)
		log.Println("Cannot connect REDIS. Exiting.")

		//
		os.Exit(-1)
	}

	//
	return c, e
}

// --------------------------------------------------------------------
// Message
type Rmsg struct {
	///
	Channel string
	Message string
}

// --------------------------------------------------------------------
// Follower is defined by its channel for receiving messages
type Follower struct {
	//
	Inputs chan Rmsg
}

// --------------------------------------------------------------------
// create a new follower and make him registered
func NewFollower() *Follower {
	//
	p := &Follower{
		//
		Inputs: make(chan Rmsg, 100),
	}

	//
	Global._followers.Lock()
	Global.followers = append(Global.followers, p)
	Global._followers.Unlock()

	//
	return p
}

// --------------------------------------------------------------------
// Global singleton for REDIS connection and listening
type MConn struct {
	// sockets for output and input
	// (the input socket is a message subscriber)
	handleOUT redis.Conn
	handleIN  redis.Conn

	// listening
	followers  []*Follower
	_followers sync.Mutex

	// channel for publishing
	topublish chan Rmsg

	//
	Running bool
}

// --------------------------------------------------------------------
// thread that publishes messages
func sender() {
	//
	for {
		// get a message from the channel
		msg, _ := <-Global.topublish

		// send it
		Global.handleOUT.Send("publish", msg.Channel, msg.Message)
		Global.handleOUT.Flush()
	}
}

// --------------------------------------------------------------------
// listening thread
func listener() {
	//
	psc := redis.PubSubConn{Conn: Global.handleIN}

	// set your subscriptions
	psc.PSubscribe("vm.*")

	//
	for {
		// wait & receive a message
		switch v := psc.Receive().(type) {

		//
		case redis.Message:
			// construct a message
			_rmsg := Rmsg{v.Channel, string(v.Data)}

			//
			Global._followers.Lock()

			// ... and distribute it among the followers
			for _, v := range Global.followers {
				// ...
				v.Inputs <- _rmsg
			}

			//
			Global._followers.Unlock()

		//
		case redis.Subscription:
			//

		case error:
			return
		}
	}
}

// --------------------------------------------------------------------
// main system procedure. Initialization of the core.
// --------------------------------------------------------------------
func RServerInit() *MConn {
	// ----------------------------------------------------------------
	//
	var err, err2 error

	// open the sockets
	Global.handleOUT, err = dial()
	Global.handleIN, err2 = dial()

	// make published-messages channel
	Global.topublish = make(chan Rmsg, 100)

	// ...
	Global.Running = true

	//
	if err != nil || err2 != nil {
		//
		log.Println(err)

		//
		return nil
	}

	// ----------------------------------------------------------------
	// start threads....
	go sender()
	go listener()

	// ----------------------------------------------------------------
	//
	return &Global
}

// ----------------------------------------------------------------------
// Entity descriptor
type Entity struct {
	// activating on this message
	MyTurn string

	//
	Slave *Entity

	//
	What, WhatStart, WhatEnd func()

	//
	WhatDefault func(Rmsg)
}

// --------------------------------------------------------------------
//
func MakeNewEntity() *Entity {
	//
	ent := Entity{}

	//
	ent.MyTurn = ""
	ent.Slave = nil

	//
	ent.What = func() {}
	ent.WhatStart = func() {}
	ent.WhatEnd = func() {}
	ent.WhatDefault = func(a Rmsg) {}

	//
	return &ent
}

// --------------------------------------------------------------------
//
func RPublish(channel, message string) {
	// ...
	Global.topublish <- Rmsg{channel, message}
}

// --------------------------------------------------------------------
//
func EntityExpRecReload(keys []string) bool {
	//
	if CurrentExp == nil {
		//
		log.Println("EntityExpRecReload() failed")

		//
		return false
	}

	//
	return CurrentExp.Load(keys, len(keys) == 0)
}

// ----------------------------------------------------------------------
//
func EntityStartSequence(ent *Entity, expIDChannel string) {
	//
	log.Println("New experiment started: ", expIDChannel)

	//
	CurrentExp = MakeExpID(expIDChannel)

	//
	CurrentExp.LoadAll()

	//
	ent.WhatStart()

	//
	if ent.Slave != nil {
		//
		ent.Slave.WhatStart()
	}
}

// ----------------------------------------------------------------------
//
func EntityEndSequence(ent *Entity, expIDChannel string) {
	////
	if CurrentExp != nil {
		//
		if CurrentExp.ischannel(expIDChannel) {
			//
			log.Println("Experiment ended: ", expIDChannel)

			//
			ent.WhatEnd()

			//
			if ent.Slave != nil {
				//
				ent.Slave.WhatEnd()
			}
		}
	}

	//
	CurrentExp = nil
}

// ----------------------------------------------------------------------
//
func EntityRoundSequence(ent *Entity, expIDChannel string) {
	//
	if CurrentExp != nil {
		//
		if CurrentExp.ischannel(expIDChannel) {
			//
			ent.What()
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
//
func EntityCore(ent *Entity) {
	// --------------------------------------------------------------------
	// initiate the r-sysem library (sender|listener)
	//
	_rglobal := RServerInit()

	// some errror
	if _rglobal == nil {
		//
		log.Println("R-system library start failure")

		//
		os.Exit(1)
	}

	// --------------------------------------------------------------------
	// become a new follower (receiver of messages from vm.*)
	_meFollower := NewFollower()
	_secOpt := "never never never"

	//
	if ent.Slave != nil {
		//
		_secOpt = ent.Slave.MyTurn
	}

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
			case ent.MyTurn:
				//
				EntityRoundSequence(ent, msg.Channel)

			case _secOpt:
				//
				EntityRoundSequence(ent.Slave, msg.Channel)

			default:
				//
				if CurrentExp != nil {
					//
					if CurrentExp.ischannel(msg.Channel) {
						//
						ent.WhatDefault(msg)
					}
				}
			}
		}
	}
}
