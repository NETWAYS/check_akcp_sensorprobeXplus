package akcp

import (
	"github.com/gosnmp/gosnmp"
	"fmt"
	"errors"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp/sensorProbePlus"
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


type SensorDetails struct {
	Sensortype	uint64
	Name		string
	Value		uint64
	Unit		string
	Status		uint64
}


const akcpBaseOID = ".1.3.6.1.4.1.3854"

const (
	SensorProbe_type = iota
	SecurityProbe_type
	SensorProbePlus_type
)

func GetSensorTypeInt (type_string string, device_type int) (uint32, error) {
	switch device_type {
		case SensorProbePlus_type: {
			val, ok := sensorProbePlus.SensorsTypes[type_string]
			if !ok {
				return 0, errors.New("Sensor type not found for this device")
			}
			return uint32(val), nil
		}
		default:
			// TODO
			return 0, errors.New("Other devices not yet implemented")
	}
}

func QuerySensorList(params *gosnmp.GoSNMP, device_type int) (sensors []string, err error) {
	// Fetches the IDs of all sensors
	// This ID consists of four positive integers, separated by dots (aka usable as an OID)

	var oid string

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.SensorIdListBase
		}
		default : {
			return nil, errors.New("Not yet implemented")
		}
	}
	sensors, err = GetSensorsIDsFromTable(params, oid)
	return sensors, err
}

func QueryTemperatureTable(params *gosnmp.GoSNMP, device_type int) (sensors []string, err error) {
	// Fetches the IDs of all sensors
	// This ID consists of four positive integers, separated by dots (aka usable as an OID)

	var oid string

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.TemeratureTable + ".1.1"
		}
		default : {
			return nil, errors.New("Not yet implemented")
		}
	}
	sensors, err = GetSensorsIDsFromTable(params, oid)
	return sensors, err
}

func GetSensorsIDsFromTable(params *gosnmp.GoSNMP, tableOID string) (sensors[]string, err error) {
	results, err := params.BulkWalkAll(tableOID)
	if err != nil {
		return nil, err
	}
	for _, variable := range results{
		//printValue(variable)
		sensors = append(sensors, ValueToString(variable))
		//fmt.Println(variable.Name)
	}

	//fmt.Println(sensors)
	return sensors, nil
}

func QueryHumidityTable(params *gosnmp.GoSNMP, device_type int) (sensors []string, err error) {
	// Fetches the IDs of all sensors
	// This ID consists of four positive integers, separated by dots (aka usable as an OID)

	var oid string

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.HumidityTable + ".1.1"
		}
		default : {
			return nil, errors.New("Not yet implemented")
		}
	}

	sensors, err = GetSensorsIDsFromTable(params, oid)
	return sensors, err
}


func QuerySensorDetails (params *gosnmp.GoSNMP, sensorIndex string, device_type int) (SensorDetails, error) {
	var details SensorDetails
	var tmp_oid string
	var oids = make([]string, 5, 5)

	switch device_type {
		case SensorProbePlus_type: {
			tmp_oid = akcpBaseOID
			oids[0] = tmp_oid + sensorProbePlus.SensorNameBase + "." + sensorIndex
			oids[1] = tmp_oid + sensorProbePlus.SensorTypeBase + "." + sensorIndex
			oids[2] = tmp_oid + sensorProbePlus.SensorValueBase + "." + sensorIndex
			oids[3] = tmp_oid + sensorProbePlus.SensorUnitBase + "." + sensorIndex
			oids[4] = tmp_oid + sensorProbePlus.SensorStatusBase + "." + sensorIndex
		}
		default : {
			return details, errors.New("Not yet implemented")
		}
	}

	//fmt.Println(oids)
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

	// Name
	details.Name = ValueToString(query.Variables[0])

	// Sensor type
	details.Sensortype, err = ValueToUint64(query.Variables[1])
	if err != nil {
		return details, err
}

	// The sensor Value (as seen in the interface)
	details.Value, err = ValueToUint64(query.Variables[2])
	if err != nil {
		return details, err
	}

	// The measuring unit (if any)
	details.Unit = ValueToString(query.Variables[3])

	// Sensor status (is the value inside the thresholds configured on the device
	details.Status, err = ValueToUint64(query.Variables[4])
	if err != nil {
		return details, err
	}

	return details, nil
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
