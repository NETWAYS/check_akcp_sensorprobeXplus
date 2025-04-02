package akcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/akcp/sensorProbePlus"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/utils"
	"github.com/NETWAYS/go-check"
	"github.com/gosnmp/gosnmp"
)

// Sensor statuses from the AKCP MIB (probably do not apply to all sensors)
// noStatus = 1,normal = 2,highWarning = 3,highCritical = 4,lowWarning = 5,lowCritical = 6,sensorError = 7,
type snsrStts uint64

const (
	NoStatus     snsrStts = 1
	Normal       snsrStts = 2
	HighWarning  snsrStts = 3
	HighCritical snsrStts = 4
	LowWarning   snsrStts = 5
	LowCritical  snsrStts = 6
	SensorError  snsrStts = 7
)

type SensorType uint64

type MayThreshold struct {
	Present bool
	Val     check.Threshold
}

type SensorDetails struct {
	SensorType   uint64
	Name         string
	Value        float64
	Unit         string
	Status       snsrStts
	Acknowledged bool
	Warning      MayThreshold
	Critical     MayThreshold
	Description  string
}

const akcpBaseOID = ".1.3.6.1.4.1.3854"

const (
	SensorProbeType     = 1
	SecurityProbeType   = 2
	SensorProbePlusType = 3
)

func GetSensorTypeInt(typeString string, deviceType int) (uint32, error) {
	switch deviceType {
	case SensorProbePlusType:
		{
			val, ok := sensorProbePlus.SensorsTypes[typeString]
			if !ok {
				return 0, errors.New("sensor type not found for this device")
			}

			if val <= math.MaxUint32 {
				return uint32(val), nil //nolint:gosec
			}

			return 0, errors.New("value out of range")
		}
	default:
		// TODO
		return 0, errors.New("other devices not yet implemented")
	}
}

// Fetches the IDs of all sensors
// This ID consists of four positive integers, separated by dots (aka usable as an OID)
func QuerySensorList(params *gosnmp.GoSNMP, deviceType int) (sensors []string, err error) {
	var oid string

	switch deviceType {
	case SensorProbePlusType:
		{
			oid = akcpBaseOID + sensorProbePlus.SensorIdListBase
		}
	default:
		{
			return nil, errors.New("not yet implemented")
		}
	}

	return GetSensorsIDsFromTable(params, oid)
}

// nolint: gocognit
func QueryTemperatureTable(snmp *gosnmp.GoSNMP, deviceType int) ([]SensorDetails, error) {
	var oid string

	var sensors []SensorDetails

	switch deviceType {
	case SensorProbePlusType:
		{
			oid = akcpBaseOID + sensorProbePlus.TemperatureTable
		}
	default:
		{
			return nil, errors.New("not yet implemented")
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
		for _, cell := range row {
			if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureDescription+".") { // nolint: gocritic, nestif
				sensors[counter].Name = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureType+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].SensorType = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureDegree+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Value = float64(tmp)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureUnit+".") {
				sensors[counter].Unit = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureLowWarning+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				// Thresholds are off by a factor of 10 to fake decimal point numbers
				sensors[counter].Warning.Val.Lower = float64(tmp) / 10
				sensors[counter].Warning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureHighWarning+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				// Thresholds are off by a factor of 10 to fake decimal point numbers
				sensors[counter].Warning.Val.Upper = float64(tmp) / 10
				sensors[counter].Warning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureLowCritical+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				// Thresholds are off by a factor of 10 to fake decimal point numbers
				sensors[counter].Critical.Val.Lower = float64(tmp) / 10
				sensors[counter].Critical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureHighCritical+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}
				// Thresholds are off by a factor of 10 to fake decimal point numbers
				sensors[counter].Critical.Val.Upper = float64(tmp) / 10
				sensors[counter].Critical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureAcknowledge+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				if tmp == 1 {
					sensors[counter].Acknowledged = true
				} else {
					sensors[counter].Acknowledged = false
				}
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorTemperatureStatus+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Status = snsrStts(tmp)
			}
		}

		counter++
	}

	return sensors, nil
}

// Fetches the IDs of all temperature sensors in the temperature table
func GetIDsFromTemperatureTable(params *gosnmp.GoSNMP, deviceType int) (sensors []string, err error) {
	var oid string

	switch deviceType {
	case SensorProbePlusType:
		{
			oid = akcpBaseOID + sensorProbePlus.TemperatureTable + ".1.1"
		}
	default:
		{
			return nil, errors.New("not yet implemented")
		}
	}

	sensors, err = GetSensorsIDsFromTable(params, oid)

	return sensors, err
}

func GetSensorsIDsFromTable(params *gosnmp.GoSNMP, tableOID string) (sensors []string, err error) {
	results, err := params.BulkWalkAll(tableOID)
	if err != nil {
		return nil, err
	}

	for _, variable := range results {
		sensors = append(sensors, ValueToString(variable))
	}

	return sensors, nil
}

// Fetches the IDs of all humidity sensors
// This ID consists of four positive integers, separated by dots (aka usable as an OID)

func GetIDsFromHumidityTable(params *gosnmp.GoSNMP, deviceType int) (sensors []string, err error) {
	var oid string

	switch deviceType {
	case SensorProbePlusType:
		{
			oid = akcpBaseOID + sensorProbePlus.HumidityTable + ".1.1"
		}
	default:
		{
			return nil, errors.New("not yet implemented")
		}
	}

	sensors, err = GetSensorsIDsFromTable(params, oid)

	return sensors, err
}

// nolint: gocognit
func QueryHumidityTable(snmp *gosnmp.GoSNMP, deviceType int) ([]SensorDetails, error) {
	var oid string

	var sensors []SensorDetails

	switch deviceType {
	case SensorProbePlusType:
		{
			oid = akcpBaseOID + sensorProbePlus.HumidityTable
		}
	default:
		{
			return nil, errors.New("not yet implemented")
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
		for _, cell := range row {
			if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityDescription+".") { // nolint: gocritic,nestif
				sensors[counter].Name = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityType+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].SensorType = tmp
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityPercent+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Value = float64(tmp)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityUnit+".") {
				sensors[counter].Unit = ValueToString(cell.Pdu)
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityLowWarning+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Warning.Val.Lower = float64(tmp)
				sensors[counter].Warning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityHighWarning+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Warning.Val.Upper = float64(tmp)
				sensors[counter].Warning.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityLowCritical+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Critical.Val.Lower = float64(tmp)
				sensors[counter].Critical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityHighCritical+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Critical.Val.Upper = float64(tmp)
				sensors[counter].Critical.Present = true
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityAcknowledge+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				if tmp == 1 {
					sensors[counter].Acknowledged = true
				} else {
					sensors[counter].Acknowledged = false
				}
			} else if strings.HasPrefix(cell.Pdu.Name, akcpBaseOID+sensorProbePlus.SensorHumidityStatus+".") {
				tmp, err := ValueToUint64(cell.Pdu)
				if err != nil {
					return sensors, err
				}

				sensors[counter].Status = snsrStts(tmp)
			}
		}

		counter++
	}

	return sensors, nil
}

func QuerySensorDetails(params *gosnmp.GoSNMP, sensorIndex string, deviceType int) (SensorDetails, error) {
	var details SensorDetails

	var tmpOID string

	var oids = make([]string, 8)

	switch deviceType {
	case SensorProbePlusType:
		{
			tmpOID = akcpBaseOID
			oids[0] = tmpOID + sensorProbePlus.SensorNameBase + "." + sensorIndex
			oids[1] = tmpOID + sensorProbePlus.SensorTypeBase + "." + sensorIndex
			oids[2] = tmpOID + sensorProbePlus.SensorValueBase + "." + sensorIndex
			oids[3] = tmpOID + sensorProbePlus.SensorUnitBase + "." + sensorIndex
			oids[4] = tmpOID + sensorProbePlus.SensorStatusBase + "." + sensorIndex
			// common on description
			oids[5] = tmpOID + sensorProbePlus.SensorsOnDescriptionBase + "." + sensorIndex
			// common off description
			oids[6] = tmpOID + sensorProbePlus.SensorsOffDescriptionBase + "." + sensorIndex
			oids[7] = tmpOID + sensorProbePlus.SensorsValueFormatFloatBase + "." + sensorIndex
		}
	default:
		{
			return details, errors.New("not yet implemented")
		}
	}

	query, err := params.Get(oids)
	if err != nil {
		return details, err
	}

	// Name
	details.Name = ValueToString(query.Variables[0])

	// Sensor type
	details.SensorType, err = ValueToUint64(query.Variables[1])
	if err != nil {
		return details, err
	}

	// The sensor Value (as seen in the interface)
	details.Value, err = ValueIEEE754ToFloat64(query.Variables[7])
	if err != nil {
		return details, err
	}

	// The measuring unit (if any)
	details.Unit = ValueToString(query.Variables[3])

	// Sensor status (is the value inside the thresholds configured on the device
	tmp, err := ValueToUint64(query.Variables[4])
	if err != nil {
		return details, err
	}

	details.Status = snsrStts(tmp)

	if details.Status != Normal {
		details.Description = ValueToString(query.Variables[5])
	} else {
		details.Description = ValueToString(query.Variables[6])
	}

	return details, nil
}

func ValueToString(pdu gosnmp.SnmpPDU) string {
	switch pdu.Type { //nolint: exhaustive
	case gosnmp.OctetString:
		return string(pdu.Value.([]byte))
	default:
		return fmt.Sprintf("%d", gosnmp.ToBigInt(pdu.Value))
	}
}

func ValueToUint64(pdu gosnmp.SnmpPDU) (uint64, error) {
	switch pdu.Type { //nolint: exhaustive
	case gosnmp.Integer:
		var val = gosnmp.ToBigInt(pdu.Value)
		if val.IsUint64() {
			return val.Uint64(), nil
		}

		return 0, errors.New("value not in uint64")
	default:
		return 0, errors.New("value is not an integer")
	}
}

func ValueIEEE754ToFloat64(pdu gosnmp.SnmpPDU) (float64, error) {
	switch pdu.Type { //nolint: exhaustive
	case gosnmp.Opaque:
		tmp := pdu.Value.([]uint8)
		bla := binary.LittleEndian.Uint32(tmp)
		tmp2 := math.Float32frombits(bla)

		return float64(tmp2), nil
	default:
		return 0, errors.New("value is not an Opaque")
	}
}
