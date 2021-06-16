package main

import (
	"github.com/gosnmp/gosnmp"
)

// Taken from the AKCP MIB
// temperature = 1,humidity-dual = 2,temperature-dual = 3,four-20mA = 4,dcvoltage = 5,airflow = 6,dry-inout = 7,dry-in = 8,motion = 9,water = 10,security = 11,siren = 12,acvoltage = 14,relay = 19,thermocouple = 20,smoke = 21,drycontact-array = 22,temperature-array = 23,waterrope = 24,fuellevel = 25,tanksender = 26,five-drycontacts = 31,irms = 34,vrms = 35,watt = 36,energy = 37,powerfactor = 38,reactive = 39,cbstatus = 40,handlelock = 49,air-pressure = 51,ir-remote = 52,digital-amp = 53,digital-watt = 54,valve-status = 55,lcd = 86,buzzer = 87,tower-led = 88,pulse-counter = 89,flow = 90,edge-counter = 91,tanklevel-height = 92,tanklevel-volume = 93,diff-pressure = 94,tanklevel-2m = 95,tanklevel-5m = 96,tanklevel-10m = 97,tanklevel-15m = 98,tanklevel-20m = 99,thermostat = 128,virtual = 129,sound = 130,software-motion = 131,board-state = 133,power-meter = 134,access = 137,door = 138,reader = 139,
const (
	temperature = 1
	humidity_dual = 2
	temperature_dual = 3
	four_20mA = 4
	dcvoltage = 5
	airflow = 6
	dry_inout = 7
	dry_in = 8
	motion = 9
	water = 10
	security = 11
	siren = 12
	acvoltage = 14
	relay = 19
	thermocouple = 20
	smoke = 21
	drycontact_array = 22
	temperature_array = 23
	waterrope = 24
	fuellevel = 25
	tanksender = 26
	five_drycontacts = 31
	irms = 34
	vrms = 35
	watt = 36
	energy = 37
	powerfactor = 38
	reactive = 39
	cbstatus = 40
	handlelock = 49
	air_pressure = 51
	ir_remote = 52
	digital_amp = 53
	digital_watt = 54
	valve_status = 55
	lcd = 86
	buzzer = 87
	tower_led = 88
	pulse_counter = 89
	flow = 90
	edge_counter = 91
	tanklevel_height = 92
	tanklevel_volume = 93
	diff_pressure = 94
	tanklevel_2m = 95
	tanklevel_5m = 96
	tanklevel_10m = 97
	tanklevel_15m = 98
	tanklevel_20m = 99
	thermostat = 128
	virtual = 129
	sound = 130
	software_motion = 131
	board_state = 133
	power_meter = 134
	access = 137
	door = 138
	reader = 139
)

// Sensor statuses from the AKCP MIB (probably do not apply to all sensors)
// noStatus = 1,normal = 2,highWarning = 3,highCritical = 4,lowWarning = 5,lowCritical = 6,sensorError = 7,
const (
	noStatus = 1
	normal = 2
	highWarning = 3
	highCritical = 4
	lowWarning = 5
	lowCritical = 6
	sensorError = 7
)

type SensorType uint64

var SensorsTypes = map[string]SensorType {
	"temperature" : temperature,
	"humidity_dual" : humidity_dual,
	"temperature_dual" : temperature_dual,
	"four_20mA" : four_20mA,
	"dcvoltage" : dcvoltage,
	"airflow" : airflow,
	"dry_inout" : dry_inout,
	"dry_in" : dry_in,
	"motion" : motion,
	"water" : water,
	"security" : security,
	"siren" : siren,
	"acvoltage" : acvoltage,
	"relay" : relay,
	"thermocouple" : thermocouple,
	"smoke" : smoke,
	"drycontact_array" : drycontact_array,
	"temperature_array" : temperature_array,
	"waterrope" : waterrope,
	"fuellevel" : fuellevel,
	"tanksender" : tanksender,
	"five_drycontacts" : five_drycontacts,
	"irms" : irms,
	"vrms" : vrms,
	"watt" : watt,
	"energy" : energy,
	"powerfactor" : powerfactor,
	"reactive" : reactive,
	"cbstatus" : cbstatus,
	"handlelock" : handlelock,
	"air_pressure" : air_pressure,
	"ir_remote" : ir_remote,
	"digital_amp" : digital_amp,
	"digital_watt" : digital_watt,
	"valve_status" : valve_status,
	"lcd" : lcd,
	"buzzer" : buzzer,
	"tower_led" : tower_led,
	"pulse_counter" : pulse_counter,
	"flow" : flow,
	"edge_counter" : edge_counter,
	"tanklevel_height" : tanklevel_height,
	"tanklevel_volume" : tanklevel_volume,
	"diff_pressure" : diff_pressure,
	"tanklevel_2m" : tanklevel_2m,
	"tanklevel_5m" : tanklevel_5m,
	"tanklevel_10m" : tanklevel_10m,
	"tanklevel_15m" : tanklevel_15m,
	"tanklevel_20m" : tanklevel_20m,
	"thermostat" : thermostat,
	"virtual" : virtual,
	"sound" : sound,
	"software_motion" : software_motion,
	"board_state" : board_state,
	"power_meter" : power_meter,
	"access" : access,
	"door" : door,
	"reader" : reader,
}

type SensorDetails struct {
	sensortype	uint64
	name		string
	value		uint64
	unit		string
	status		uint64
}


const akcpBaseOID = ".1.3.6.1.4.1.3854"

// AKCP Subtree
const (
	akcpSensorProbe = akcpBaseOID + ".1"
	akcpSecurityProbe = akcpBaseOID + ".2"
	akcpPlusSeries = akcpBaseOID + ".3"
)

// TODO sensorProbe
// TODO securityProbe

// plusSeries
const (
	akcpPlusSeriesSensors = akcpPlusSeries + ".5"
)

const (
	PlsSrs_commonTable = akcpPlusSeriesSensors + ".1"
	PlsSrs_temeratureTable = akcpPlusSeriesSensors + ".2"
	PlsSrs_humidityTable = akcpPlusSeriesSensors + ".3"
	PlsSrs_drycontactTable = akcpPlusSeriesSensors + ".4"
	PlsSrs_current4to20mATable = akcpPlusSeriesSensors + ".5"
	PlsSrs_dcVoltageTable = akcpPlusSeriesSensors + ".6"
	PlsSrs_airflowTable = akcpPlusSeriesSensors + ".7"
	PlsSrs_motionTable = akcpPlusSeriesSensors + ".8"
	PlsSrs_waterTable = akcpPlusSeriesSensors + ".9"
	PlsSrs_securityTable = akcpPlusSeriesSensors + ".10"
	PlsSrs_sirenTable = akcpPlusSeriesSensors + ".11"
	PlsSrs_relayTable = akcpPlusSeriesSensors + ".12"
	PlsSrs_acVoltageTable = akcpPlusSeriesSensors + ".13"
	PlsSrs_smokeTable = akcpPlusSeriesSensors + ".14"
	PlsSrs_waterRopeTable = akcpPlusSeriesSensors + ".21"
	PlsSrs_powerTable = akcpPlusSeriesSensors + ".22"
	PlsSrs_fuelTable = akcpPlusSeriesSensors + ".24"
	PlsSrs_tankSenderTable = akcpPlusSeriesSensors + ".26"
	PlsSrs_doorTable = akcpPlusSeriesSensors + ".27"
	PlsSrs_temperatureArrayTable = akcpPlusSeriesSensors + ".28"
	PlsSrs_towerLEDTable = akcpPlusSeriesSensors + ".29"
	PlsSrs_enumTable = akcpPlusSeriesSensors + ".31"
)

const (
	PlsSrs_sensorBase = akcpBaseOID + ".3.5.1.1"
	sensorIdListBase = akcpBaseOID + ".3.5.1.1.1"
	sensorNameBase = akcpBaseOID + ".3.5.1.1.2"
	sensorValueBase = akcpBaseOID + ".3.5.1.1.4"
	sensorUnitBase = akcpBaseOID + ".3.5.1.1.5"
)

func querySensorList(params *gosnmp.GoSNMP) (sensors []string, err error) {
	// Fetches the IDs of all sensors
	// This ID consists of four positive integers, separated by dots (aka usable as an OID)

	results, err := params.BulkWalkAll(sensorIdListBase)
	if err != nil {
		return nil, err
	}
	for _, variable := range results{
		//printValue(variable)
		sensors = append(sensors, ValueToString(variable))
		//fmt.Println(variable.Name)
	}

	return sensors, nil
}

func querySensorDetails (params *gosnmp.GoSNMP, sensorIndex string) (SensorDetails, error) {
	var details SensorDetails

	var oids = make([]string, 5, 5)
	oids[0] = PlsSrs_sensorBase + ".2." + sensorIndex
	oids[1] = PlsSrs_sensorBase + ".3." + sensorIndex
	oids[2] = PlsSrs_sensorBase + ".4." + sensorIndex
	oids[3] = PlsSrs_sensorBase + ".5." + sensorIndex
	oids[4] = PlsSrs_sensorBase + ".6." + sensorIndex
	/*
	var oids = [...]string {
		PlsSrs_sensorBase + ".2" + sensorIndex, // Description
		PlsSrs_sensorBase + ".3" + sensorIndex, // Type
		PlsSrs_sensorBase + ".4" + sensorIndex,  // Value
		PlsSrs_sensorBase + ".5" + sensorIndex,  // Unit
		PlsSrs_sensorBase + ".6" + sensorIndex,  // Status
	}
	*/
	query, err := params.Get(oids)
	if err != nil {
		return details, err
	}

	/*
	fmt.Println("querySensorDetails:")
	for i, value := range query.Variables {
		fmt.Printf("%d: %s\n", i, value)
	}
	*/

	details.name = ValueToString(query.Variables[0])
	details.sensortype, err = ValueToUint64(query.Variables[1])
	if err != nil {
		return details, err
	}
	details.value, err = ValueToUint64(query.Variables[2])
	if err != nil {
		return details, err
	}
	details.unit = ValueToString(query.Variables[3])
	details.status, err = ValueToUint64(query.Variables[4])
	if err != nil {
		return details, err
	}

	return details, nil
}
