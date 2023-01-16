package sensorProbePlus

import ()

//"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp"

// Taken from the AKCP MIB
// temperature = 1,humidity-dual = 2,temperature-dual = 3,four-20mA = 4,dcvoltage = 5,airflow = 6,dry-inout = 7,dry-in = 8,motion = 9,water = 10,security = 11,siren = 12,acvoltage = 14,relay = 19,thermocouple = 20,smoke = 21,drycontact-array = 22,temperature-array = 23,waterrope = 24,fuellevel = 25,tanksender = 26,five-drycontacts = 31,irms = 34,vrms = 35,watt = 36,energy = 37,powerfactor = 38,reactive = 39,cbstatus = 40,handlelock = 49,air-pressure = 51,ir-remote = 52,digital-amp = 53,digital-watt = 54,valve-status = 55,lcd = 86,buzzer = 87,tower-led = 88,pulse-counter = 89,flow = 90,edge-counter = 91,tanklevel-height = 92,tanklevel-volume = 93,diff-pressure = 94,tanklevel-2m = 95,tanklevel-5m = 96,tanklevel-10m = 97,tanklevel-15m = 98,tanklevel-20m = 99,thermostat = 128,virtual = 129,sound = 130,software-motion = 131,board-state = 133,power-meter = 134,access = 137,door = 138,reader = 139,
const (
	Temperature       = 1
	Humidity_dual     = 2
	Temperature_dual  = 3
	Four_20mA         = 4
	Dcvoltage         = 5
	Airflow           = 6
	Dry_inout         = 7
	Dry_in            = 8
	Motion            = 9
	Water             = 10
	Security          = 11
	Siren             = 12
	Acvoltage         = 14
	Relay             = 19
	Thermocouple      = 20
	Smoke             = 21
	Drycontact_array  = 22
	Temperature_array = 23
	Waterrope         = 24
	Fuellevel         = 25
	Tanksender        = 26
	Five_drycontacts  = 31
	Irms              = 34
	Vrms              = 35
	Watt              = 36
	Energy            = 37
	Powerfactor       = 38
	Reactive          = 39
	Cbstatus          = 40
	Handlelock        = 49
	Air_pressure      = 51
	Ir_remote         = 52
	Digital_amp       = 53
	Digital_watt      = 54
	Valve_status      = 55
	Lcd               = 86
	Buzzer            = 87
	Tower_led         = 88
	Pulse_counter     = 89
	Flow              = 90
	Edge_counter      = 91
	Tanklevel_height  = 92
	Tanklevel_volume  = 93
	Diff_pressure     = 94
	Tanklevel_2m      = 95
	Tanklevel_5m      = 96
	Tanklevel_10m     = 97
	Tanklevel_15m     = 98
	Tanklevel_20m     = 99
	Thermostat        = 128
	Virtual           = 129
	Sound             = 130
	Software_motion   = 131
	Board_state       = 133
	Power_meter       = 134
	Access            = 137
	Door              = 138
	Reader            = 139
)

type SensorType uint64

var SensorsTypes = map[string]SensorType{
	"temperature":       Temperature,
	"humidity_dual":     Humidity_dual,
	"temperature_dual":  Temperature_dual,
	"four_20mA":         Four_20mA,
	"dcvoltage":         Dcvoltage,
	"airflow":           Airflow,
	"dry_inout":         Dry_inout,
	"dry_in":            Dry_in,
	"motion":            Motion,
	"water":             Water,
	"security":          Security,
	"siren":             Siren,
	"acvoltage":         Acvoltage,
	"relay":             Relay,
	"thermocouple":      Thermocouple,
	"smoke":             Smoke,
	"drycontact_array":  Drycontact_array,
	"temperature_array": Temperature_array,
	"waterrope":         Waterrope,
	"fuellevel":         Fuellevel,
	"tanksender":        Tanksender,
	"five_drycontacts":  Five_drycontacts,
	"irms":              Irms,
	"vrms":              Vrms,
	"watt":              Watt,
	"energy":            Energy,
	"powerfactor":       Powerfactor,
	"reactive":          Reactive,
	"cbstatus":          Cbstatus,
	"handlelock":        Handlelock,
	"air_pressure":      Air_pressure,
	"ir_remote":         Ir_remote,
	"digital_amp":       Digital_amp,
	"digital_watt":      Digital_watt,
	"valve_status":      Valve_status,
	"lcd":               Lcd,
	"buzzer":            Buzzer,
	"tower_led":         Tower_led,
	"pulse_counter":     Pulse_counter,
	"flow":              Flow,
	"edge_counter":      Edge_counter,
	"tanklevel_height":  Tanklevel_height,
	"tanklevel_volume":  Tanklevel_volume,
	"diff_pressure":     Diff_pressure,
	"tanklevel_2m":      Tanklevel_2m,
	"tanklevel_5m":      Tanklevel_5m,
	"tanklevel_10m":     Tanklevel_10m,
	"tanklevel_15m":     Tanklevel_15m,
	"tanklevel_20m":     Tanklevel_20m,
	"thermostat":        Thermostat,
	"virtual":           Virtual,
	"sound":             Sound,
	"software_motion":   Software_motion,
	"board_state":       Board_state,
	"power_meter":       Power_meter,
	"access":            Access,
	"door":              Door,
	"reader":            Reader,
}

const (
	PlusSeriesID = ".3"
	Sensors      = PlusSeriesID + ".5"

	SensorBase                  = Sensors + ".1.1"
	SensorIdListBase            = SensorBase + ".1"
	SensorNameBase              = SensorBase + ".2"
	SensorTypeBase              = SensorBase + ".3"
	SensorValueBase             = SensorBase + ".4"
	SensorUnitBase              = SensorBase + ".5"
	SensorStatusBase            = SensorBase + ".6"
	SensorsOnDescriptionBase    = SensorBase + ".52"
	SensorsOffDescriptionBase   = SensorBase + ".53"
	SensorsValueFormatFloatBase = SensorBase + ".99"

	CommonTable           = Sensors + ".1"
	TemperatureTable      = Sensors + ".2"
	HumidityTable         = Sensors + ".3"
	DrycontactTable       = Sensors + ".4"
	Current4to20mATable   = Sensors + ".5"
	DcVoltageTable        = Sensors + ".6"
	AirflowTable          = Sensors + ".7"
	MotionTable           = Sensors + ".8"
	WaterTable            = Sensors + ".9"
	SecurityTable         = Sensors + ".10"
	SirenTable            = Sensors + ".11"
	RelayTable            = Sensors + ".12"
	AcVoltageTable        = Sensors + ".13"
	SmokeTable            = Sensors + ".14"
	WaterRopeTable        = Sensors + ".21"
	PowerTable            = Sensors + ".22"
	FuelTable             = Sensors + ".24"
	TankSenderTable       = Sensors + ".26"
	DoorTable             = Sensors + ".27"
	TemperatureArrayTable = Sensors + ".28"
	TowerLEDTable         = Sensors + ".29"
	NumTable              = Sensors + ".31"
)

const (
	TemperatureTableEntry = TemperatureTable + ".1"

	SensorTemperatureIndex             = TemperatureTableEntry + ".1"
	SensorTemperatureDescription       = TemperatureTableEntry + ".2"
	SensorTemperatureType              = TemperatureTableEntry + ".3"
	SensorTemperatureDegree            = TemperatureTableEntry + ".4"
	SensorTemperatureUnit              = TemperatureTableEntry + ".5"
	SensorTemperatureStatus            = TemperatureTableEntry + ".6"
	SensorTemperatureGoOffline         = TemperatureTableEntry + ".8"
	SensorTemperatureLowCritical       = TemperatureTableEntry + ".9"
	SensorTemperatureLowWarning        = TemperatureTableEntry + ".10"
	SensorTemperatureHighWarning       = TemperatureTableEntry + ".11"
	SensorTemperatureHighCritical      = TemperatureTableEntry + ".12"
	SensorTemperatureRearm             = TemperatureTableEntry + ".13"
	SensorTemperatureDelayError        = TemperatureTableEntry + ".14"
	SensorTemperatureDelayNormal       = TemperatureTableEntry + ".15"
	SensorTemperatureDelayLowCritical  = TemperatureTableEntry + ".16"
	SensorTemperatureDelayLowWarning   = TemperatureTableEntry + ".17"
	SensorTemperatureDelayHighWarning  = TemperatureTableEntry + ".18"
	SensorTemperatureDelayHighCritical = TemperatureTableEntry + ".19"
	SensorTemperatureRaw               = TemperatureTableEntry + ".20"
	SensorTemperatureOffset            = TemperatureTableEntry + ".21"
	SensorTemperaturePort              = TemperatureTableEntry + ".35"
	SensorTemperatureSubPort           = TemperatureTableEntry + ".36"
	SensorTemperatureDisplayStyle      = TemperatureTableEntry + ".45"
	SensorTemperatureHighCriticalDesc  = TemperatureTableEntry + ".46"
	SensorTemperatureLowCriticalDesc   = TemperatureTableEntry + ".47"
	SensorTemperatureNormalDesc        = TemperatureTableEntry + ".48"
	SensorTemperatureLowWarningDesc    = TemperatureTableEntry + ".49"
	SensorTemperatureHighWarningDesc   = TemperatureTableEntry + ".50"
	SensorTemperatureSensorErrorDesc   = TemperatureTableEntry + ".51"
	SensorTemperatureHighCriticalColor = TemperatureTableEntry + ".54"
	SensorTemperatureLowCriticalColor  = TemperatureTableEntry + ".55"
	SensorTemperatureNormalColor       = TemperatureTableEntry + ".56"
	SensorTemperatureLowWarningColor   = TemperatureTableEntry + ".57"
	SensorTemperatureHighWarningColor  = TemperatureTableEntry + ".57"
	SensorTemperatureSensorErrorColor  = TemperatureTableEntry + ".59"
	SensorTemperatureAcknowledge       = TemperatureTableEntry + ".70"
	SensorTemperatureSerialNumber      = TemperatureTableEntry + ".71"
	SensorTemperatureId                = TemperatureTableEntry + ".1000"
)

const (
	HumidityTableEntry = HumidityTable + ".1"

	SensorHumidityIndex             = HumidityTableEntry + ".1"
	SensorHumidityDescription       = HumidityTableEntry + ".2"
	SensorHumidityType              = HumidityTableEntry + ".3"
	SensorHumidityPercent           = HumidityTableEntry + ".4"
	SensorHumidityUnit              = HumidityTableEntry + ".5"
	SensorHumidityStatus            = HumidityTableEntry + ".6"
	SensorHumidityGoOffline         = HumidityTableEntry + ".8"
	SensorHumidityLowCritical       = HumidityTableEntry + ".9"
	SensorHumidityLowWarning        = HumidityTableEntry + ".10"
	SensorHumidityHighWarning       = HumidityTableEntry + ".11"
	SensorHumidityHighCritical      = HumidityTableEntry + ".12"
	SensorHumidityRearm             = HumidityTableEntry + ".13"
	SensorHumidityDelayError        = HumidityTableEntry + ".14"
	SensorHumidityDelayNormal       = HumidityTableEntry + ".15"
	SensorHumidityDelayLowCritical  = HumidityTableEntry + ".16"
	SensorHumidityDelayLowWarning   = HumidityTableEntry + ".17"
	SensorHumidityDelayHighWarning  = HumidityTableEntry + ".18"
	SensorHumidityDelayHighCritical = HumidityTableEntry + ".19"
	SensorHumidityRaw               = HumidityTableEntry + ".20"
	SensorHumidityOffset            = HumidityTableEntry + ".21"
	SensorHumidityPort              = HumidityTableEntry + ".35"
	SensorHumiditySubPort           = HumidityTableEntry + ".36"
	SensorHumidityDisplayStyle      = HumidityTableEntry + ".45"
	SensorHumidityHighCriticalDesc  = HumidityTableEntry + ".46"
	SensorHumidityLowCriticalDesc   = HumidityTableEntry + ".47"
	SensorHumidityNormalDesc        = HumidityTableEntry + ".48"
	SensorHumidityLowWarningDesc    = HumidityTableEntry + ".49"
	SensorHumidityHighWarningDesc   = HumidityTableEntry + ".50"
	SensorHumiditySensorErrorDesc   = HumidityTableEntry + ".51"
	SensorHumidityHighCriticalColor = HumidityTableEntry + ".54"
	SensorHumidityLowCriticalColor  = HumidityTableEntry + ".55"
	SensorHumidityNormalColor       = HumidityTableEntry + ".56"
	SensorHumidityLowWarningColor   = HumidityTableEntry + ".57"
	SensorHumidityHighWarningColor  = HumidityTableEntry + ".57"
	SensorHumiditySensorErrorColor  = HumidityTableEntry + ".59"
	SensorHumidityAcknowledge       = HumidityTableEntry + ".70"
	SensorHumidityId                = HumidityTableEntry + ".1000"
)
