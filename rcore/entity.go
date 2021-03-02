//
package rcore

// --------------------------------------------------------------------
// ...
import (
	"fmt"
	"sync"
	"log"
	"os"
	"github.com/gomodule/redigo/redis"
)

// --------------------------------------------------------------------
//
var global MConn
var err error

// --------------------------------------------------------------------
//
func dial() (redis.Conn, error) {
	//
	return redis.Dial("tcp", ":6379")
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
}

// --------------------------------------------------------------------
//
func addFollower(who *Follower) {
	//
	global._followers.Lock()

	//
	global.followers = append(global.followers, who)

	//
	global._followers.Unlock()
}

// --------------------------------------------------------------------
//
func sender() {
	//
	for {
		//
		msg, _ := <- global.topublish

		//
		global.handleOUT.Send("publish", msg.Channel, msg.Message)
		global.handleOUT.Flush()
	}
}

// --------------------------------------------------------------------
//
func listener() {
	//
	psc := redis.PubSubConn{Conn: global.handleIN}

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
					global._followers.Lock()

					//
					for _, v := range global.followers {
						//
						v.Inputs <- _rmsg
					}

					//
					global._followers.Unlock()
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
func demo() {
	//
	me := newFollower(); addFollower(me)

	//
	for {
		//
		msg := <- me.Inputs

		//
		fmt.Println("Nekdo nam pise: ", msg.Channel, " ", msg.Message)
	}
}


// --------------------------------------------------------------------
//
func RServerInit() *MConn {
	//
	global.handleOUT, err = dial()
	global.handleIN, err = dial()

	//
	global.topublish = make(chan rmsg, 100)

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
	return &global
}

// --------------------------------------------------------------------
//
func RPublish(channel, message string) {
	//
	global.topublish <- rmsg{ channel, message }
}

// ----------------------------------------------------------------------
//
func EntityCore(myTurn string, what func ()) {
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

  // --------------------------------------------------------------------
  //
	for {
    //
		msg := <- _meFollower.Inputs

    //
    switch msg.Message {
      //
      case CallStart:
        //
        log.Println("New experiment started: ", msg.Channel)
        CurrentExp = MakeExpID(msg.Channel)

      //
      case CallEnd:
        //
        if CurrentExp != nil {
          //
          if CurrentExp.Varname == msg.Channel {
            //
            log.Println("Experiment ended: ", msg.Channel)
            CurrentExp = nil
          }
        }

      //
    case myTurn:
        //
        if CurrentExp != nil {
          //
          if CurrentExp.Varname == msg.Channel {
            //
            what()
          }
        }

      default:
        ;
    }
	}
}

// --------------------------------------------------------------------
//
func mainRCore() {
	//
	global.handleOUT, err = dial()
	global.handleIN, err = dial()

	//
	global.topublish = make(chan rmsg, 100)

	//
	if err != nil {
		//
		fmt.Println(err)

		//
		return
	}

	//
	defer global.handleIN.Close()
	defer	global.handleOUT.Close()

	//
	//go valuesMedicalRoo()
	/*
	c.Do("SET", "k1", 1)
	n, _ := redis.Int(c.Do("GET", "k1"))
	fmt.Printf("%#v\n", n)
	n, _ = redis.Int(c.Do("INCR", "k1"))
	fmt.Printf("%#v\n", n)
	*/

	//
	go sender()
	go listener()
	go valuesMedicalRoot()

	//
	demo()

	//
}
