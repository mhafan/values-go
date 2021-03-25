//
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
var Global MConn

// REDIS instance for long term experimenting
var redis_host = "pchrubym.fit.vutbr.cz"

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
		os.Exit(-1)
	}

	//
	return c, e
}

// --------------------------------------------------------------------
// Message
type rmsg struct {
	///
	Channel string
	Message string
}

// --------------------------------------------------------------------
// Follower is defined by its channel for receiving messages
type Follower struct {
	//
	Inputs chan rmsg
}

// --------------------------------------------------------------------
// create a new follower and make him registered
func NewFollower() *Follower {
	//
	p := &Follower{
		//
		Inputs: make(chan rmsg, 100),
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
	topublish chan rmsg

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
			_rmsg := rmsg{v.Channel, string(v.Data)}

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
	Global.topublish = make(chan rmsg, 100)

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

// --------------------------------------------------------------------
//
func RPublish(channel, message string) {
	// ...
	Global.topublish <- rmsg{channel, message}
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
func EntityStartSequence(expIDChannel string, whatStart func()) {
	//
	log.Println("New experiment started: ", expIDChannel)

	//
	CurrentExp = MakeExpID(expIDChannel)

	//
	CurrentExp.Load([]string{}, true)

	//
	whatStart()
}

// ----------------------------------------------------------------------
//
func EntityEndSequence(expIDChannel string, whatEnd func()) {
	////
	if CurrentExp != nil {
		//
		if CurrentExp.Varname == expIDChannel {
			//
			log.Println("Experiment ended: ", expIDChannel)

			//
			whatEnd()
		}

		//
		CurrentExp = nil
	}
}

// ----------------------------------------------------------------------
//
func EntityRoundSequence(expIDChannel string, what func()) {
	//
	if CurrentExp != nil {
		//
		if CurrentExp.Varname == expIDChannel {
			//
			what()
		}
	}
}

// ----------------------------------------------------------------------
//
func EntityMasterChannel(msg rmsg) {
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
func EntityCore(myTurn string, what, whatStart, whatEnd func()) {
	// --------------------------------------------------------------------
	// initiate the r-sysem library (sender|listener)
	_rglobal := RServerInit()

	// some errror
	if _rglobal == nil {
		//
		log.Println("R-system library start failure")
		os.Exit(1)
	}

	// --------------------------------------------------------------------
	// become a new follower (receiver of messages from vm.*)
	_meFollower := NewFollower()

	//
	defer Global.handleIN.Close()
	defer Global.handleOUT.Close()

	// --------------------------------------------------------------------
	//
	for Global.Running == true {
		//
		msg := <-_meFollower.Inputs

		//
		if msg.Channel == MasterChannel {
			//
			EntityMasterChannel(msg)
		} else {
			//
			switch msg.Message {
			//
			case CallStart:
				//
				EntityStartSequence(msg.Channel, whatStart)

			//
			case CallEnd:
				//
				EntityEndSequence(msg.Channel, whatEnd)

				//
			case myTurn:
				//
				EntityRoundSequence(msg.Channel, what)

			default:

			}
		}
	}
}
