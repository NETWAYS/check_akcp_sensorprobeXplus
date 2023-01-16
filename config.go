package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/akcp"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/akcp/sensorProbePlus"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/gosnmp/gosnmp"
	"github.com/spf13/pflag"
)

type Config struct {
	hostname              string
	snmp_version_param    string
	snmp_version          gosnmp.SnmpVersion
	community             string
	port                  uint16
	mode                  string
	device                string
	device_type           int
	sensorPort            string
	excludeSensorType     []string
	excludeSensorType_int []uint32
	// authProtocol		string
	// authPassword 	string
	// privProtocol		string
	// privPassword 	string
}

// Modes
const (
	queryAllSensors uint64 = iota
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
	runTestSuccess
)

var modes = map[string]uint64{
	"queryAllSensors":    queryAllSensors,
	"single":             single,
	"temperatureSensors": temperaturSensors,
	"humiditySensors":    humiditySensors,
	"run_test_success":   runTestSuccess,
}

//var default_excluded_types []string

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	fs.StringVarP(&c.hostname, "host", "h", "", "Hostname or IP of the targeted device (required)")
	fs.StringVarP(&c.snmp_version_param, "snmp_version", "", "2c", "Version of SNMP to use (1|2c)")
	fs.StringVarP(&c.community, "community", "c", "public", "SNMP Community string")
	fs.Uint16VarP(&c.port, "port", "p", 161, "SNMP Port")
	fs.StringVarP(&c.device, "device", "", "sensorProbe+", `Device type, may be one of:
	- sensorProbe
	- securityProbe
	- sensorProbe+
	`)
	fs.StringVarP(&c.mode, "mode", "m", "queryAllSensors", `Usage mode (default: queryAllSensors)
	Possible modes:
	- queryAllSensors: Query all the sensors and show their value and state
	- single: Query a single sensor (sensorPort must be set)
	- temperatureSensors: Query all the temperature sensors
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
	- reader
	`)
	fs.StringVarP(&c.sensorPort, "sensorPort", "", "", "Sensor Port (required for single mode)")
	fs.StringArrayVarP(&c.excludeSensorType, "exclude", "e", nil, "Exclude specific sensor type, valid types are the same as are available for querying above, can be used multiple times")
}

func (c *Config) Validate() error {

	val, ok := modes[c.mode]
	if !ok {
		if c.device == "sensorProbePlus" {
			_, ok = sensorProbePlus.SensorsTypes[c.mode]
			if !ok {
				return errors.New("Mode is not a valid value")
			}
		}
	} else if val == runTestSuccess {
		check.ExitRaw(0, "It seems like you can execute this programm")
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

	switch c.device {
	case "sensorProbe":
		{
			// TODO
			return errors.New("Not yet implemented.")
		}
	case "securityProbe":
		{
			// TODO
			return errors.New("Not yet implemented.")
		}
	case "sensorProbe+":
		{
			c.device_type = akcp.SensorProbePlus_type
		}
	default:
		{
			return errors.New("Invalid device type.")
		}
	}

	if c.excludeSensorType == nil {
		c.excludeSensorType = append(c.excludeSensorType, "buzzer")
	}
	if c.excludeSensorType[0] != "" {
		for _, tmp := range c.excludeSensorType {
			if tmp == "" {
				continue
			}
			val, err := akcp.GetSensorTypeInt(tmp, c.device_type)
			if err != nil {
				return err
			}
			c.excludeSensorType_int = append(c.excludeSensorType_int, val)
		}
	}
	return nil
}

func (c *Config) Run(overall *result.Overall) (err error) {

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		check.ExitError(err)
	}

	params := &gosnmp.GoSNMP{
		Target:    c.hostname,
		Port:      c.port,
		Community: c.community,
		Version:   c.snmp_version,
		Timeout:   timeout,
		Retries:   3,
	}

	err = params.Connect()
	if err != nil {
		check.ExitError(err)
	}
	defer params.Conn.Close()

	// Get name, location and type
	query, err := params.Get([]string{".1.3.6.1.4.1.3854.3.2.1.9.0", ".1.3.6.1.4.1.3854.3.2.1.10.0", ".1.3.6.1.4.1.3854.3.2.1.8.0"})
	if err != nil {
		check.ExitError(err)
	}

	name := query.Variables[0].Value
	location := query.Variables[1].Value
	devType := query.Variables[2].Value

	overall.Summary = fmt.Sprintf("Device %s at location %s (%s)", name, location, devType)

	val, ok := modes[c.mode]
	if !ok {
		// not one of the main modes, look for specifics
		if c.device_type == akcp.SensorProbePlus_type {
			val, ok := sensorProbePlus.SensorsTypes[c.mode]
			if !ok {
				return errors.New("Mode is not a valid value")
			} else {
				err = querySensorByType(params, c, overall, c.device_type, uint64(val))
				if err != nil {
					check.ExitError(err)
				}
				return nil
			}
			//return errors.New("Mode not yet implemented.")
		} else {
			// TODO
			return errors.New("Device not yet implemented.")
		}
	} else {
		if val == queryAllSensors {
			err = queryAllSensorsMode(params, c, overall, c.device_type)
			if err != nil {
				check.ExitError(err)
			}
			return nil
		} else if val == temperaturSensors {
			err = queryTemperatureSensors(params, c, overall, c.device_type)
			if err != nil {
				check.ExitError(err)
			}
			return nil
		} else if val == humiditySensors {
			err = queryHumiditySensors(params, c, overall, c.device_type)
			if err != nil {
				check.ExitError(err)
			}
			return nil
		} else {
			// TODO
			return errors.New("Not yet implemented.")
		}
	}

}

func queryAllSensorsMode(params *gosnmp.GoSNMP, c *Config, overall *result.Overall, device_type int) (err error) {

	// Get all sensors
	sensors, err := akcp.QuerySensorList(params, device_type)
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		details, err := akcp.QuerySensorDetails(params, sensor, device_type)
		if err != nil {
			check.ExitError(err)
		}

		var exclude bool = false
		for _, excluded_type := range c.excludeSensorType_int {
			if uint64(excluded_type) == details.SensorType {
				exclude = true
			}
		}
		if exclude {
			continue
		}

		if details.SensorType == sensorProbePlus.Temperature ||
			details.SensorType == sensorProbePlus.Temperature_dual {
			tempSensors, err := akcp.QueryTemperatureTable(params, device_type)
			if err != nil {
				return err
			}
			for _, tempSensor := range tempSensors {
				if details.Name == tempSensor.Name {
					details.Warning = tempSensor.Warning
					details.Critical = tempSensor.Critical
				}
			}
		}
		if details.SensorType == sensorProbePlus.Humidity_dual {
			humiSensors, err := akcp.QueryHumidityTable(params, device_type)
			if err != nil {
				return err
			}
			for _, humSensor := range humiSensors {
				if details.Name == humSensor.Name {
					details.Warning = humSensor.Warning
					details.Critical = humSensor.Critical
				}
			}
		}

		err = mapSensorStatus(details, overall)
		if err != nil {
			return err
		}
	}

	return nil
}

func mapSensorStatus(sensor akcp.SensorDetails, overall *result.Overall) error {
	var sensorString string
	if sensor.SensorType == sensorProbePlus.Motion {
		sensorString = fmt.Sprintf("%s: %s", sensor.Name, sensor.Description)
	} else {
		sensorString = fmt.Sprintf("%s: %.1f", sensor.Name, sensor.Value)
	}

	unit := ""
	if sensor.Unit == "C" {
		unit = "℃"
	} else if sensor.Unit != "" {
		unit += sensor.Unit
	}

	sensorString += unit

	var pf perfdata.Perfdata
	pf.Label = sensor.Name
	pf.Value = sensor.Value

	if sensor.Warning.Present {
		pf.Warn = &sensor.Warning.Val
	}
	if sensor.Critical.Present {
		pf.Crit = &sensor.Critical.Val
	}

	if sensor.Unit != "C" {
		pf.Uom = strings.ToLower(sensor.Unit)
	} else {
		pf.Uom = "C"
	}

	if sensor.Status == 2 {
		overall.AddOK(sensorString + " | " + pf.String())
	} else if sensor.Status == 3 || sensor.Status == 5 {
		if sensor.Warning.Present && (float64(sensor.Value) <= sensor.Warning.Val.Lower) {
			sensorString += fmt.Sprintf(" is lower than warning threshold %.1f%s", sensor.Warning.Val.Lower, unit)
		} else if sensor.Warning.Present && (float64(sensor.Value) >= sensor.Warning.Val.Upper) {
			sensorString += fmt.Sprintf(" is higher than warning threshold %.1f%s", sensor.Warning.Val.Upper, unit)
		}
		overall.AddWarning(sensorString + " | " + pf.String())
	} else if sensor.Status == 6 || sensor.Status == 4 {
		if sensor.Critical.Present && (float64(sensor.Value) <= sensor.Critical.Val.Lower) {
			sensorString += fmt.Sprintf(" is lower than critical threshold %.1f%s", sensor.Critical.Val.Lower, unit)
		} else if sensor.Critical.Present && (float64(sensor.Value) >= sensor.Critical.Val.Upper) {
			sensorString += fmt.Sprintf(" is higher than critical threshold %.1f%s", sensor.Critical.Val.Upper, unit)
		}
		overall.AddCritical(sensorString + " | " + pf.String())
	} else if sensor.Status == 7 {
		overall.AddCritical(sensor.Name + " ERROR!" + " | " + pf.String())
	} else {
		overall.AddUnknown(sensorString + " | " + pf.String())
	}

	return nil
}

func querySensorByType(params *gosnmp.GoSNMP, c *Config, overall *result.Overall, device_type int, sensor_type uint64) (err error) {
	// Get all sensors
	sensors, err := akcp.QuerySensorList(params, device_type)
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		details, err := akcp.QuerySensorDetails(params, sensor, device_type)
		if err != nil {
			check.ExitError(err)
		}

		if details.SensorType == uint64(sensor_type) {
			err = mapSensorStatus(details, overall)
			if err != nil {
				check.ExitError(err)
			}
		}
	}
	return nil
}

func queryTemperatureSensors(params *gosnmp.GoSNMP, c *Config, overall *result.Overall, device_type int) (err error) {

	// Get all sensors
	// TODO: Error Handling
	sensors, _ := akcp.QueryTemperatureTable(params, device_type)

	for _, details := range sensors {
		err = mapSensorStatus(details, overall)
		if err != nil {
			return err
		}
	}
	return nil
}

func queryHumiditySensors(params *gosnmp.GoSNMP, c *Config, overall *result.Overall, device_type int) (err error) {

	// Get all sensors
	sensors, err := akcp.QueryHumidityTable(params, device_type)
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		err = mapSensorStatus(sensor, overall)
		if err != nil {
			return err
		}
	}
	return nil
}
