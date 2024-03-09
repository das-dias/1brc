package main

import "testing"

func Test_parse_station_name(t *testing.T) {
	test_table := map[string]struct {
		given string
		expected_key string
		expected_value uint8
	}{
		"easy": {
			given: "Buenos Aires;22.5",
			expected_key: "Buenos Aires",
			expected_value: 12,
		},
		"key with special chars": {
			given: "St. John's;15.2",
			expected_key: "St. John's",
			expected_value: 10,
		},
	}
	for name, test := range test_table {
		t.Run(name, func(t *testing.T)) {
			result, length := parse_station_name(MAX_STATION_NAME_LENGTH, test.given)
			if result != test.expected_key {
				t.Errorf("Expected %s, but got %s", test.expected_key, result)
			}
			if length != test.expected_value {
				t.Errorf("Expected %d, but got %d", test.expected_value, length)
			}
		}
	}
}

func Test_parse_station_temperature(t *testing.T) {
	test_table := map[string]struct {
		given string
		expected int16
	}{
		"positive": {
			given: "22.5",
			expected: 225,
		},
		"negative": {
			given: "-22.5",
			expected: -225,
		},
	}
	for name, test := range test_table {
		t.Run(name, func(t *testing.T)) {
			result := parse_station_temperature(MAX_TEMP_VALUE_LENGTH, test.given)
			if result != test.expected {
				t.Errorf("Expected %d, but got %d", test.expected, result)
			}
		}
	}
}

func Test_compute_min_mean_max(t *testing.T) {
	test_table := map[string]struct {
		given_min int16
		given_mean int16
		given_max int16
		given_new_temperature int16
		expected_min int16
		expected_mean int16
		expected_max int16
	}{
		"new min": {
			given_min: 0,
			given_mean: 0,
			given_max: 0,
			given_new_temperature: -225,
			expected_min: -225,
			expected_mean: -113,
			expected_max: 0,
		},
		"new max": {
			given_min: 0,
			given_mean: 0,
			given_max: 0,
			given_new_temperature: 225,
			expected_min: 0,
			expected_mean: 113,
			expected_max: 225,
		},
		"new mean": {
			given_min: 122,
			given_mean: 100,
			given_max: 150,
			given_new_temperature: 125,
			expected_min: 122,
			expected_mean: 112,
			expected_max: 150,
		},
	}
	for name, test := range test_table {
		t.Run(name, func(t *testing.T)) {
			min, mean, max := compute_min_mean_max(test.given_min, test.given_mean, test.given_max, test.given_new_temperature)
			if min != test.expected_min {
				t.Errorf("Expected %d, but got %d", test.expected_min, min)
			}
			if mean != test.expected_mean {
				t.Errorf("Expected %d, but got %d", test.expected_mean, mean)
			}
			if max != test.expected_max {
				t.Errorf("Expected %d, but got %d", test.expected_max, max)
			}
		}
	}
}

func Test_parse_line(t *testing.T) {
	test_table := map[string]struct {
		given string
		expected_key string
		expected_value int16
	}{
		"positive temperature": {
			given: "Buenos Aires;22.5",
			expected_key: "Buenos Aires",
			expected_value: 225,
		},
		"key with special chars and positive temperature": {
			given: "St. John's;15.2",
			expected_key: "St. John's",
			expected_value: 152,
		},
		"negative temperature": {
			given: "St. John's;-15.2",
			expected_key: "St. John's",
			expected_value: -152,
		},
	}
	for name, test := range test_table {
		t.Run(name, func(t *testing.T)) {
			result_key, result_value := parse_line(MAX_STATION_NAME_LENGTH, MAX_TEMP_VALUE_LENGTH, test.given)
			if result_key != test.expected_key {
				t.Errorf("Expected %s, but got %s", test.expected_key, result_key)
			}
			if result_value != test.expected_value {
				t.Errorf("Expected %d, but got %d", test.expected_value, result_value)
			}
		}
	}
}
