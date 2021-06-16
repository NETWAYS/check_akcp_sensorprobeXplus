package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
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
	sensor			string
	// authProtocol		string
	// authPassword 	string
	// privProtocol		string
	// privPassword 	string
}

// Modes
const (
	listSensors uint = iota
)

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	fs.StringVarP(&c.hostname, "host", "h", "", "Hostname or IP of the targeted device (required)")
	fs.StringVarP(&c.snmp_version_param, "snmp_version", "", "2c", "Version of SNMP to use (1|2c, default: 2c)")
	fs.StringVarP(&c.community, "community", "c", "public", "SNMP Community string (default: \"public\")")
	fs.Uint16VarP(&c.port, "port", "p", 161, "SNMP Port (default: \"161\")")
}

func (c *Config) Validate() error {
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
	sensors, err := querySensorList(params)
	if err != nil {
		check.ExitError(err)
	}

	for _, sensor := range sensors {
		//fmt.Printf("%d: %s\n", num, sensor)
		details, err := querySensorDetails(params, sensor)
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

		sensorString := fmt.Sprintf("%s: %d", details.name, details.value)
		if details.unit == "%" {
			sensorString += "%%"
		}
		//fmt.Println(sensorString)
		if details.status == 2 {
			overall.AddOK(sensorString)
		} else if details.status == 3 || details.status == 5 {
			overall.AddWarning(sensorString)
		} else if details.status == 6 || details.status == 4 {
			overall.AddCritical(sensorString)
		} else {
			overall.AddUnknown(sensorString)
		}
	}

	return nil
}

func ValueToString(pdu gosnmp.SnmpPDU) string {
	switch pdu.Type {
	case gosnmp.OctetString:
		return string(pdu.Value.([]byte))
	default:
		//fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
		return fmt.Sprint("%d", gosnmp.ToBigInt(pdu.Value))
	}
}

func ValueToUint64(pdu gosnmp.SnmpPDU) (uint64, error) {
	switch pdu.Type {
	case gosnmp.Integer:
		var val = gosnmp.ToBigInt(pdu.Value)
		if val.IsUint64() {
			return val.Uint64(), nil
		} else {
			return 0, errors.New("Value not in uint64")
		}
	default:
		return 0, errors.New("Value is not an integer")
	}
}

func printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	return nil
}
