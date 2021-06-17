package sensorProbePlus

import (
	//"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp"
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
type SensorType uint32

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

const (
	PlusSeriesID = ".3"
	Sensors = PlusSeriesID + ".5"

	SensorBase = Sensors + ".1.1"
	SensorIdListBase = SensorBase + ".1"
	SensorNameBase = SensorBase + ".2"
	SensorTypeBase = SensorBase + ".3"
	SensorValueBase = SensorBase + ".4"
	SensorUnitBase = SensorBase + ".5"
	SensorStatusBase = SensorBase + ".6"

	CommonTable = Sensors + ".1"
	TemeratureTable = Sensors + ".2"
	HumidityTable = Sensors + ".3"
	DrycontactTable = Sensors + ".4"
	Current4to20mATable = Sensors + ".5"
	DcVoltageTable = Sensors + ".6"
	AirflowTable = Sensors + ".7"
	MotionTable = Sensors + ".8"
	WaterTable = Sensors + ".9"
	SecurityTable = Sensors + ".10"
	SirenTable = Sensors + ".11"
	RelayTable = Sensors + ".12"
	AcVoltageTable = Sensors + ".13"
	SmokeTable = Sensors + ".14"
	WaterRopeTable = Sensors + ".21"
	PowerTable = Sensors + ".22"
	FuelTable = Sensors + ".24"
	TankSenderTable = Sensors + ".26"
	DoorTable = Sensors + ".27"
	TemperatureArrayTable = Sensors + ".28"
	TowerLEDTable = Sensors + ".29"
	NumTable = Sensors + ".31"
)
