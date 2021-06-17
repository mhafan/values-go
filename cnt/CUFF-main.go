// ----------------------------------------------------------------------
// Main procedure for TOF/PTC Cuff sensor
// ----------------------------------------------------------------------
// Activity description:
// 1) Patient Model outputs Cinp and its PTC/TOF stimations
// ----------------------------------------------------------------------

//
package main

//
import (
	"rcore"
)

// ----------------------------------------------------------------------
//
func sensorKeys() []string {
	//
	return []string{"SesorStutus", "SensorTOF0", "SensorTOF1", "SensorTOF2", "SensorTOF3", "SensorPTC"}
}

// ----------------------------------------------------------------------
//
func saveSensorStatus(exp *rcore.Exprec) {
	//
	exp.Save(sensorKeys(), false)
}

// ----------------------------------------------------------------------
//
func resetSensorStatus(exp *rcore.Exprec, andSave bool) {
	//
	exp.SensorStatus = 0
	exp.SensorTOF0 = 0
	exp.SensorTOF1 = 0
	exp.SensorTOF2 = 0
	exp.SensorTOF3 = 0
	exp.SensorPTC = 0

	//
	if andSave {
		//
		saveSensorStatus(exp)
	}
}

// ----------------------------------------------------------------------
// invokek when SENSOR command received
func CNTCuffMain(msg rcore.Rmsg) {
	// reset outputs from the sensor, dont save tho
	resetSensorStatus(rcore.CurrentExp, false)

	// dont forget to call next in the line (TCM)
	defer rcore.CurrentExp.Say(rcore.CallTCM)

	// save
	saveSensorStatus(rcore.CurrentExp)
}
