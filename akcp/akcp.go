package akcp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NETWAYS/check_akcp_sensorprobeXplus/akcp/sensorProbePlus"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/utils"
	"github.com/gosnmp/gosnmp"
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

type MayUint64 struct {
	Present bool
	Val 	uint64
}


type SensorDetails struct {
	SensorType	uint64
	Name		string
	Value		uint64
	Unit		string
	Status		uint64
	Acknowledged bool
	LowWarning 	MayUint64
	HighWarning	MayUint64
	LowCritical 	MayUint64
	HighCritical	MayUint64
}


const akcpBaseOID = ".1.3.6.1.4.1.3854"

const (
	SensorProbe_type = 1
	SecurityProbe_type = 2
	SensorProbePlus_type = 3
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

func QueryTemperatureTable(snmp *gosnmp.GoSNMP, device_type int) ([]SensorDetails, error) {
	var oid string
	var sensors []SensorDetails

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.TemperatureTable
		}
		default : {
			return nil, errors.New("Not yet implemented")
		}
	}

	tempTable, err := snmp.BulkWalkAll(oid)
	if err != nil {
		return sensors, err
	}

	foo, err := utils.ParseSnmpTable(&tempTable, 12)
	if err != nil {
		return sensors, err
	}

	sensors = make([]SensorDetails, len(*foo))

	counter := 0
	for _, row := range *foo {
		// Every row is a temperature sensor
		//fmt.Printf("%#v:\n", row)
		for _, cell := range row {
			//fmt.Printf("%#v\n", cell)
			if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureDescription + ".") {
				sensors[counter].Name = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureType + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].SensorType = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureDegree + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].Value = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureUnit + ".") {
				sensors[counter].Unit = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureLowWarning + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].LowWarning.Val = tmp
				sensors[counter].LowWarning.Present =  true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureHighWarning + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].HighWarning.Val = tmp
				sensors[counter].HighWarning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureLowCritical + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].LowCritical.Val = tmp
				sensors[counter].LowCritical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureHighCritical + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].HighCritical.Val = tmp
				sensors[counter].HighCritical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureAcknowledge + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				if tmp == 1 {
					sensors[counter].Acknowledged = true
				} else {
					sensors[counter].Acknowledged = false
				}
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorTemperatureStatus + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].Status = tmp
			}
		}
		counter ++
	}

	return sensors, nil
}

func GetIDsFromTemperatureTable(params *gosnmp.GoSNMP, device_type int) (sensors []string, err error) {
	// Fetches the IDs of all temperature sensors in the temperature table

	var oid string

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.TemperatureTable + ".1.1"
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

func GetIDsFromHumidityTable(params *gosnmp.GoSNMP, device_type int) (sensors []string, err error) {
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

func QueryHumidityTable(snmp *gosnmp.GoSNMP, device_type int) ([]SensorDetails, error) {
	var oid string
	var sensors []SensorDetails

	switch device_type {
		case SensorProbePlus_type: {
			oid = akcpBaseOID + sensorProbePlus.HumidityTable
		}
		default : {
			return nil, errors.New("Not yet implemented")
		}
	}

	tempTable, err := snmp.BulkWalkAll(oid)
	if err != nil {
		return sensors, err
	}

	foo, err := utils.ParseSnmpTable(&tempTable, 12)
	if err != nil {
		return sensors, err
	}

	sensors = make([]SensorDetails, len(*foo))

	counter := 0
	for _, row := range *foo {
		// Every row is a temperature sensor
		//fmt.Printf("%#v:\n", row)
		for _, cell := range row {
			//fmt.Printf("%#v\n", cell)
			if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityDescription + ".") {
				sensors[counter].Name = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityType + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].SensorType = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityPercent+ ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].Value = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityUnit + ".") {
				sensors[counter].Unit = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityLowWarning + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].LowWarning.Val = tmp
				sensors[counter].LowWarning.Present =  true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityHighWarning + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].HighWarning.Val = tmp
				sensors[counter].HighWarning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityLowCritical + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].LowCritical.Val = tmp
				sensors[counter].LowCritical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityHighCritical + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].HighCritical.Val = tmp
				sensors[counter].HighCritical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityAcknowledge + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				if tmp == 1 {
					sensors[counter].Acknowledged = true
				} else {
					sensors[counter].Acknowledged = false
				}
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID + sensorProbePlus.SensorHumidityStatus + ".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				sensors[counter].Status = tmp
			}
		}
		counter ++
	}

	return sensors, nil
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
	details.SensorType, err = ValueToUint64(query.Variables[1])
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
		return fmt.Sprintf("%d", gosnmp.ToBigInt(pdu.Value))
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
