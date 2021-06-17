package akcp

import (
	"github.com/gosnmp/gosnmp"
	"fmt"
	"errors"
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

func QuerySensorList(params *gosnmp.GoSNMP) (sensors []string, err error) {
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

func QuerySensorDetails (params *gosnmp.GoSNMP, sensorIndex string) (SensorDetails, error) {
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

	details.Name = ValueToString(query.Variables[0])
	details.Sensortype, err = ValueToUint64(query.Variables[1])
	if err != nil {
		return details, err
	}
	details.Value, err = ValueToUint64(query.Variables[2])
	if err != nil {
		return details, err
	}
	details.Unit = ValueToString(query.Variables[3])
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
