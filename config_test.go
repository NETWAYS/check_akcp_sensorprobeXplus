package main

import (
	"strings"
	"testing"

	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/akcp"
	"github.com/NETWAYS/check_akcp_sensorprobeXplus/internal/akcp/sensorProbePlus"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

func TestSensorStatus(t *testing.T) {

	testcases := map[string]struct {
		sensor   akcp.SensorDetails
		expected string
	}{
		"normal": {
			sensor: akcp.SensorDetails{
				Unit:   "C",
				Status: akcp.Normal,
			},
			expected: "states: ok=1\n\\_ [OK] : 0.0",
		},
		"HighWarning": {
			sensor: akcp.SensorDetails{
				Unit:   "Foo",
				Status: akcp.HighWarning,
			},
			expected: "[WARNING] : 0.0Foo",
		},
		"LowCritical": {
			sensor: akcp.SensorDetails{
				SensorType: sensorProbePlus.Motion,
				Unit:       "Bar",
				Value:      1,
				Warning: akcp.MayThreshold{
					Present: true,
					Val:     check.Threshold{Lower: 10},
				},
				Status: akcp.LowCritical,
			},
			expected: "[CRITICAL] : Bar",
		},
		"Error": {
			sensor: akcp.SensorDetails{
				Unit:   "Error",
				Status: akcp.SensorError,
			},
			expected: "[CRITICAL]  ERROR!",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			overall := &result.Overall{}

			err := mapSensorStatus(tc.sensor, overall)

			if err != nil {
				t.Error("Expected no error, got %w", err)
			}

			actual := overall.GetOutput()

			if !strings.Contains(actual, tc.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", tc.expected)
			}
		})
	}
}
