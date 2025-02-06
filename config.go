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
	hostname                 string
	snmpVersionParam         string
	snmpVersion              gosnmp.SnmpVersion
	community                string
	port                     uint16
	mode                     string
	device                   string
	deviceType               int
	sensorPort               string
	excludeSensorType        []string
	excludeSensorTypeInteger []uint32
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

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	fs.StringVarP(&c.hostname, "host", "h", "", "Hostname or IP of the targeted device (required)")
	fs.StringVarP(&c.snmpVersionParam, "snmp_version", "", "2c", "Version of SNMP to use (1|2c)")
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
	if !ok { // nolint: gocritic,nestif
		if c.device == "sensorProbePlus" {
			_, ok = sensorProbePlus.SensorsTypes[c.mode]
			if !ok {
				return errors.New("mode is not a valid value")
			}
		}
	} else if val == runTestSuccess {
		check.ExitRaw(0, "It seems like you can execute this programm")
	} else if val == single && c.sensorPort == "" {
		return errors.New("no sensorPort was given")
	}

	switch c.snmpVersionParam {
	case "1":
		c.snmpVersion = gosnmp.Version1
	case "2c":
		c.snmpVersion = gosnmp.Version2c
	case "3":
		return errors.New("SNMP v3 not yet implemented")
	default:
		return errors.New("invalid SNMP version string")
	}

	switch c.device {
	case "sensorProbe":
		{
			// TODO
			return errors.New("not yet implemented")
		}
	case "securityProbe":
		{
			// TODO
			return errors.New("not yet implemented")
		}
	case "sensorProbe+":
		{
			c.deviceType = akcp.SensorProbePlusType
		}
	default:
		{
			return errors.New("invalid device type")
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

			val, err := akcp.GetSensorTypeInt(tmp, c.deviceType)

			if err != nil {
				return err
			}

			c.excludeSensorTypeInteger = append(c.excludeSensorTypeInteger, val)
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
		Version:   c.snmpVersion,
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
	if !ok { //nolint:nestif
		// not one of the main modes, look for specifics
		if c.deviceType == akcp.SensorProbePlusType {
			val, ok := sensorProbePlus.SensorsTypes[c.mode]
			if !ok {
				return errors.New("mode is not a valid value")
			}

			err = querySensorByType(params, c, overall, c.deviceType, uint64(val))

			if err != nil {
				check.ExitError(err)
			}

			return nil
		}

		// TODO
		return errors.New("device not yet implemented")
	}

	switch val {
	case queryAllSensors:
		err = queryAllSensorsMode(params, c, overall, c.deviceType)
		if err != nil {
			check.ExitError(err)
		}

		return nil
	case temperaturSensors:
		err = queryTemperatureSensors(params, c, overall, c.deviceType)
		if err != nil {
			check.ExitError(err)
		}

		return nil
	case humiditySensors:
		err = queryHumiditySensors(params, c, overall, c.deviceType)
		if err != nil {
			check.ExitError(err)
		}

		return nil
	default:
		return errors.New("not yet implemented")
	}
}

// nolint: gocognit
func queryAllSensorsMode(params *gosnmp.GoSNMP, c *Config, overall *result.Overall, deviceType int) (err error) {
	sensors, err := akcp.QuerySensorList(params, deviceType) // Get all sensors
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		details, err := akcp.QuerySensorDetails(params, sensor, deviceType)

		if err != nil {
			check.ExitError(err)
		}

		exclude := false

		for _, excludedType := range c.excludeSensorTypeInteger {
			if uint64(excludedType) == details.SensorType {
				exclude = true
			}
		}

		if exclude {
			continue
		}

		if details.SensorType == sensorProbePlus.Temperature ||
			details.SensorType == sensorProbePlus.Temperature_dual {
			tempSensors, err := akcp.QueryTemperatureTable(params, deviceType)

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
			humiSensors, err := akcp.QueryHumidityTable(params, deviceType)

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

	sc := result.PartialResult{}
	_ = sc.SetDefaultState(check.Unknown)
	sc.Output = sensorString

	switch sensor.Status {
	case akcp.Normal:
		_ = sc.SetState(check.OK)
		sc.Perfdata.Add(&pf)
	case akcp.HighWarning, akcp.LowWarning:
		_ = sc.SetState(check.Warning)
		sc.Perfdata.Add(&pf)
	case akcp.HighCritical, akcp.LowCritical:
		_ = sc.SetState(check.Critical)
		sc.Perfdata.Add(&pf)
	case akcp.SensorError:
		_ = sc.SetState(check.Critical)
		sc.Output = sensor.Name + " ERROR!"
		sc.Perfdata.Add(&pf)
	case akcp.NoStatus:
		_ = sc.SetState(check.Unknown)
		sc.Output = sensor.Name + " is unknown (No Status)!"
	}

	overall.AddSubcheck(sc)

	return nil
}

func querySensorByType(params *gosnmp.GoSNMP, _ *Config, overall *result.Overall, deviceType int, sensorType uint64) error { //nolint:unparam
	sensors, err := akcp.QuerySensorList(params, deviceType) // Get all sensors

	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		details, err := akcp.QuerySensorDetails(params, sensor, deviceType)
		if err != nil {
			check.ExitError(err)
		}

		if details.SensorType == sensorType {
			err = mapSensorStatus(details, overall)
			if err != nil {
				check.ExitError(err)
			}
		}
	}

	return nil
}

func queryTemperatureSensors(params *gosnmp.GoSNMP, _ *Config, overall *result.Overall, deviceType int) (err error) {
	sensors, _ := akcp.QueryTemperatureTable(params, deviceType) // Get all sensors
	// TODO: Error Handling

	for _, details := range sensors {
		err = mapSensorStatus(details, overall)
		if err != nil {
			return err
		}
	}

	return nil
}

func queryHumiditySensors(params *gosnmp.GoSNMP, _ *Config, overall *result.Overall, deviceType int) (err error) {
	sensors, err := akcp.QueryHumidityTable(params, deviceType) // Get all sensors
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
