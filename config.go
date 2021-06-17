package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp/sensorProbePlus"
	"github.com/gosnmp/gosnmp"
	"github.com/spf13/pflag"
)


type Config struct {
	hostname		string
	snmp_version_param string
	snmp_version	gosnmp.SnmpVersion
	community		string
	port			uint16
	mode			string
	device			string
	sensorPort			string
	// authProtocol		string
	// authPassword 	string
	// privProtocol		string
	// privPassword 	string
}

// Modes
const (
	queryAllSensors uint = iota
	listPossibleSensors
	single
	temperaturSensors
	humiditySensors
	airflowSensors
	drycontactSensors
	current4to20mA
	dcVoltage
	acVoltage
	waterRope
	power
	fuel
	tankSender
	temperatureArray
	towerLED
)

var modes = map[string] uint {
	"queryAllSensors" : queryAllSensors,
	"single" :	single,
	"temperatureSensors" : temperaturSensors,
	"humiditySensors" : humiditySensors,
}

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	fs.StringVarP(&c.hostname, "host", "h", "", "Hostname or IP of the targeted device (required)")
	fs.StringVarP(&c.snmp_version_param, "snmp_version", "", "2c", "Version of SNMP to use (1|2c, default: 2c)")
	fs.StringVarP(&c.community, "community", "c", "public", "SNMP Community string (default: \"public\")")
	fs.Uint16VarP(&c.port, "port", "p", 161, "SNMP Port (default: \"161\")")
	fs.StringVarP(&c.device, "device", "", "sensorProbe+", `Device type, may be one of:
	- sensorProbe
	- securityProbe
	- sensorProbe+`)
	fs.StringVarP(&c.mode, "mode", "m", "queryAllSensors", `Usage mode (default: queryAllSensors)
	Possible modes:
	- queryAllSensors: Query all the sensors and show their value and state
	- single: Query a single sensor (sensorPort must be set)
	- temperaturSensors: Query all the temperature sensors
	- humiditySensors: Query all the humidity sensors

	The following modes will query the respective sensor types:
	- temperature
	- humidity_dual
	- temperature_dual
	- four_20mA
	- dcvoltage
	- airflow
	- dry_inout
	- dry_in
	- motion
	- water
	- security
	- siren
	- acvoltage
	- relay
	- thermocouple
	- smoke
	- drycontact_array
	- temperature_array
	- waterrope
	- fuellevel
	- tanksender
	- five_drycontacts
	- irms
	- vrms
	- watt
	- energy
	- powerfactor
	- reactive
	- cbstatus
	- handlelock
	- air_pressure
	- ir_remote
	- digital_amp
	- digital_watt
	- valve_status
	- lcd
	- buzzer
	- tower_led
	- pulse_counter
	- flow
	- edge_counter
	- tanklevel_height
	- tanklevel_volume
	- diff_pressure
	- tanklevel_2m
	- tanklevel_5m
	- tanklevel_10m
	- tanklevel_15m
	- tanklevel_20m
	- thermostat
	- virtual
	- sound
	- software_motion
	- board_state
	- power_meter
	- access
	- door
	- reader`)
	fs.StringVarP(&c.sensorPort, "sensorPort", "", "", "Sensor Port (required for single mode)")
}

func (c *Config) Validate() error {

	val , ok := modes[c.mode]
	if !ok {
		if c.device == "sensorProbePlus" {
			_, ok = sensorProbePlus.SensorsTypes[c.mode]
			if !ok {
				return errors.New("Mode is not a valid value")
			}
		}
	} else {
		if val == single && c.sensorPort == "" {
				return errors.New("No sensorPort was given")
		}
	}

	if c.snmp_version_param == "1" {
		c.snmp_version = gosnmp.Version1
	} else if c.snmp_version_param == "2c" {
		c.snmp_version = gosnmp.Version2c
	} else if c.snmp_version_param == "3" {
		return errors.New("SNMP v3 not yet implemented")
		//c.snmp_version = gosnmp.Version3
	} else {
		return errors.New("Invalid SNMP version string")
	}
	return nil
}

func (c *Config) Run(overall *result.Overall ) (err error) {

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		check.ExitError(err)
	}

	params := &gosnmp.GoSNMP{
		Target:		c.hostname,
		Port:		c.port,
		Community:	c.community,
		Version:	c.snmp_version,
		Timeout:	timeout,
		Retries:	3,
	}

	err = params.Connect()
	if err != nil {
		check.ExitError(err)
	}
	//fmt.Println("Connection succesful")
	defer params.Conn.Close()

	// Get all sensors
	sensors, err := akcp.QuerySensorList(params)
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		//fmt.Printf("%d: %s\n", num, sensor)
		details, err := akcp.QuerySensorDetails(params, sensor)
		if err != nil {
			check.ExitError(err)
		}
		/*
		fmt.Printf("Name: %s\n", details.name)
		fmt.Printf("Sensor Type: %d\n", details.sensortype)
		fmt.Printf("Sensor Value: %d\n", details.value)
		fmt.Printf("Unit: %s\n", details.unit)
		fmt.Printf("Status: %d\n", details.status)
		*/

		sensorString := fmt.Sprintf("%s: %d", details.Name, details.Value)
		if details.Unit == "%" {
			sensorString += "%%"
		}
		//fmt.Println(sensorString)
		if details.Status == 2 {
			overall.AddOK(sensorString)
		} else if details.Status == 3 || details.Status == 5 {
			overall.AddWarning(sensorString)
		} else if details.Status == 6 || details.Status == 4 {
			overall.AddCritical(sensorString)
		} else {
			overall.AddUnknown(sensorString)
		}
	}

	return nil
}
