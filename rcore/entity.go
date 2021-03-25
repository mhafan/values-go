//
package rcore

// --------------------------------------------------------------------
// ...
import (
	"sync"
	"log"
	"os"
	"flag"
	"github.com/gomodule/redigo/redis"
)

// --------------------------------------------------------------------
//
var Global MConn
var err error

//
var redis_host = "pchrubym.fit.vutbr.cz"

//
var flHostname = flag.String("h", ":6379", "REDIS hostname")
var flAuth = flag.String("a", "", "REDIS auth password")

// --------------------------------------------------------------------
//
func dial() (redis.Conn, error) {
	//
	opta := redis.DialPassword(*flAuth)

	//
	c, e := redis.Dial("tcp", *flHostname, opta)

	//
	if e != nil {
		//
		log.Println("Dial error: ", e); os.Exit(-1)
	}

	//
	return c, e
}

// --------------------------------------------------------------------
//
type rmsg struct {
		Channel string
		Message string
}

// --------------------------------------------------------------------
//
type Follower struct {
	//
	Inputs chan rmsg
}


// --------------------------------------------------------------------
//
func NewFollower() *Follower {
	//
	p := newFollower()

	//
	if p != nil {
		//
		addFollower(p)
	}

	//
	return p
}

//
func newFollower() *Follower {
	//
	f := Follower{}
	f.Inputs = make(chan rmsg, 100)

	//
	return &f
}

// --------------------------------------------------------------------
//
type MConn struct {
	//
	handleOUT redis.Conn
	handleIN redis.Conn

	//
	followers [] *Follower
	_followers sync.Mutex

	//
	topublish chan rmsg

	//
	Running bool
}

// --------------------------------------------------------------------
//
func addFollower(who *Follower) {
	//
	Global._followers.Lock()

	//
	Global.followers = append(Global.followers, who)

	//
	Global._followers.Unlock()
}

// --------------------------------------------------------------------
//
func sender() {
	//
	for {
		//
		msg, _ := <- Global.topublish

		//
		Global.handleOUT.Send("publish", msg.Channel, msg.Message)
		Global.handleOUT.Flush()
	}
}

// --------------------------------------------------------------------
//
func listener() {
	//
	psc := redis.PubSubConn{Conn: Global.handleIN}

	//
	psc.PSubscribe("vm.*")

	//
	for {
			//
	    switch v := psc.Receive().(type) {
	    case redis.Message:
					//
					_rmsg := rmsg{v.Channel, string(v.Data)}

					//
					Global._followers.Lock()

					//
					for _, v := range Global.followers {
						//
						v.Inputs <- _rmsg
					}

					//
					Global._followers.Unlock()
					//
	    case redis.Subscription:
				;
	    case error:
	        return
	    }
	}
}



// --------------------------------------------------------------------
//
func RServerInit() *MConn {
	//
	Global.handleOUT, err = dial()
	Global.handleIN, err = dial()

	//
	Global.topublish = make(chan rmsg, 100)

	//
	Global.Running = true

	//
	if err != nil {
		//
		log.Println(err)

		//
		return nil
	}

	//
	go sender()
	go listener()

	//
	return &Global
}

// --------------------------------------------------------------------
//
func RPublish(channel, message string) {
	//
	Global.topublish <- rmsg{ channel, message }
}


// --------------------------------------------------------------------
//
func EntityExpRecReload(keys [] string) bool {
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
func EntityStartSequence(expIDChannel string, whatStart func ()) {
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
func EntityEndSequence(expIDChannel string, whatEnd func ()) {
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
func EntityRoundSequence(expIDChannel string, what func ()) {
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
		return;
	}
}

// ----------------------------------------------------------------------
//
func EntityCore(myTurn string, what, whatStart, whatEnd func ()) {
  // --------------------------------------------------------------------
  // initiate the r-sysem library (sender|listener)
  _rglobal := RServerInit()

  // some errror
  if _rglobal == nil {
    //
    log.Println("R-system library start failure"); os.Exit(1)
  }

  // --------------------------------------------------------------------
  // become a new follower (receiver of messages from vm.*)
  _meFollower := NewFollower()

	//
	defer Global.handleIN.Close()
	defer	Global.handleOUT.Close()

  // --------------------------------------------------------------------
  //
	for Global.Running == true {
    //
		msg := <- _meFollower.Inputs

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
	        ;
	    }
		}
	}
}
