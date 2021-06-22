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
	"fmt"
	"rcore"
)

// ----------------------------------------------------------------------
//
type CUFFEntity struct {
	//
	TOFIntervals rcore.Scheduling
	PTCIntervals rcore.Scheduling
}

// ----------------------------------------------------------------------
//
func MakeCUFFEntity() *CUFFEntity {
	//
	out := &CUFFEntity{}

	//
	return out
}

// ----------------------------------------------------------------------
//
func sensorKeys() []string {
	//
	return []string{"SensorStatus", "SensorTOF0", "SensorTOF1", "SensorTOF2", "SensorTOF3", "SensorPTC"}
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
// TOF Measurement
func TOFMeas() rcore.LinScale {
	//
	out := rcore.LinScale{}

	//
	return out
}

// ----------------------------------------------------------------------
// Main model of TOF/PTC measurements
func (e *CUFFEntity) CNTCuffDoMeasurement(tnow int) {
	//
	if e.TOFIntervals.IsTimeToFire(tnow) {
		//
		rcore.CurrentExp.SensorStatus += rcore.CuffCommandTOF

		//
		fmt.Println("Doing TOF", e.TOFIntervals.LastExecutedEvent, e.TOFIntervals.PlannedInterval)

		//
		e.TOFIntervals = e.TOFIntervals.PlanNext(tnow)
	}

	//
	if e.PTCIntervals.IsTimeToFire(tnow) {
		//
		rcore.CurrentExp.SensorStatus += rcore.CuffCommandPTC

		//
		e.PTCIntervals = e.PTCIntervals.PlanNext(tnow)
	}
}

// ----------------------------------------------------------------------
//
func (e *CUFFEntity) MyTurn() string {
	//
	return rcore.CallSensor
}

// ----------------------------------------------------------------------
//
func (e *CUFFEntity) ResetState() {
	//
}

// ----------------------------------------------------------------------
//
func (e *CUFFEntity) StartFunction() {
	//
	t0 := rcore.CurrentExp.Mtime

	//
	e.TOFIntervals = rcore.MakeScheduling(2 * 60).PlanNext(t0)
	e.PTCIntervals = rcore.MakeScheduling(5 * 60).PlanNext(t0)

	//
}

// ----------------------------------------------------------------------
//
func (e *CUFFEntity) EndFunction() {
	//
}

// ----------------------------------------------------------------------
//
func (e *CUFFEntity) DefaultFunction(msg rcore.Rmsg) {
	//
}

// ----------------------------------------------------------------------
// invokek when SENSOR command received
func (e *CUFFEntity) CycleFunction() {
	// ------------------------------------------------------------------
	// ... first update
	if rcore.EntityExpRecReload([]string{"SensorCommand", "mtime"}) == false {
		//
		return
	}

	// reset outputs from the sensor, dont save tho
	resetSensorStatus(rcore.CurrentExp, false)

	//
	tnow := rcore.CurrentExp.Mtime

	//
	e.CNTCuffDoMeasurement(tnow)

	// dont forget to call next in the line (TCM)
	defer rcore.CurrentExp.Say(rcore.CallTCM)

	// save
	saveSensorStatus(rcore.CurrentExp)
}
