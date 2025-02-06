package utils

import (
	"strings"

	"github.com/gosnmp/gosnmp"
)

type Cell struct {
	Pdu gosnmp.SnmpPDU
	ID  string
}

// Get a list of PDUs and try to form a table from it
// prefixLength parts of the oid are ignored from the beginning of the oid.
// The suffix determines the id of entries from the end of the oid
// its length is computed by (number of oid points - 1 - prefixLength)
//
// Assumptions:
// All OIDs in the list _table_ have the same length (= number of separators)
func ParseSnmpTable(table *[]gosnmp.SnmpPDU, prefixLength uint) (*map[string][]Cell, error) {
	var result = make(map[string][]Cell)
	// compute parameters
	for _, value := range *table {
		oid := value.Name
		tmp := strings.Split(oid, ".")
		id := tmp[prefixLength]
		rowID := strings.Join(tmp[prefixLength+1:len(tmp)-1], ".")

		entry, ok := result[rowID]
		cellVal := Cell{
			ID:  id,
			Pdu: value,
		}

		if !ok {
			result[rowID] = make([]Cell, 1)
			result[rowID][0] = cellVal
		} else {
			result[rowID] = append(entry, cellVal)
		}
	}

	return &result, nil
}
