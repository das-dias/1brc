package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const MAX_STATIONS uint32 = 100000
const MAX_STATION_NAME_LENGTH uint8 = 100
const MAX_TEMP_VALUE_LENGTH uint8 = 5
const MAX_TEMP_VALUE_DECIMALS uint8 = 1

func parse_station_name(max_length uint8, line string) (uint8, string) {
	var i uint8 = 0
	temp_station_name := [max_length]byte{0}
	for ; (i < max_length) && (line[i] != ';'); i++ {
		temp_station_name[i] = line[i]
	}
	return i, string(temp_station_name)
}

func parse_station_temperature(max_length uint8, line string) int16 {
	var i uint8 = 0
	temperature_buffer := [max_length]byte{0}
	for ; (i < max_length) && (line[i] != 0); i++ {
		temperature_buffer[i] = line[i]
	}
	var uint_temp uint16 = 0
	i = 0
	if temperature_buffer[0] == '-' {
		i = 1
	}
	for ; (temperature_buffer[i] != '.') && (temperature_buffer[i] != 0); i++ {
		uint_temp = (uint_temp * 10) + uint16(temperature_buffer[i]-'0')
	}
	i++ // ignore the '.' separation point
	for ; (i < max_length) && (temperature_buffer[i] != 0); i++ {
		uint_temp = (uint_temp * 10) + uint16(temperature_buffer[i]-'0')
	}
	var int_temp int16 = int16(uint_temp)
	if temperature_buffer[0] == '-' {
		int_temp = -int_temp
	}
	return int_temp
}

func compute_min_mean_max(min int16, mean int16, max int16, new_temperature int16) (int16, int16, int16) {
	var new_mean int16 = (mean + new_temperature) >> 1
	var new_min int16 = min
	var new_max int16 = max
	if new_temperature < min {
		new_min = new_temperature
	} else { //if int_temp > max[station_idx] { // if it is not the new minimum it is the new maximum or equal
		new_max = new_temperature
	}
	return new_min, new_mean, new_max
}

func parse_line(max_station_name_length uint8, max_temp_value_length uint8, line string) (string, int16) {
	var (line_idx uint8 station_name string) = parse_station_name(max_station_name_length, line)
	var station_temperature int16 = parse_station_temperature(max_temp_value_length, line[line_idx:])
	return station_name, station_temperature
}

func build_output_string(station_name string, min int16, mean int16, max int16) string {
	min := float32(min) / 10.0
	max := float32(max) / 10.0
	mean := float32(mean) / 10.0
	return fmt.Sprintf("%s;%.1f;%.1f;%.1f", station_name, min, mean, max)
}

func main() {

	start_time := time.Now()

	min := [MAX_STATIONS]int16{0}
	max := [MAX_STATIONS]int16{0}
	mean := [MAX_STATIONS]int16{0}
	var current_station_idx uint16 = 0
	station_name_idx := make(map[string]uint16)

	fp, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		var (station_name string temperature int16) = 
			parse_line(MAX_STATION_NAME_LENGTH, MAX_TEMP_VALUE_LENGTH, line)
		
		if _, exists := station_name_idx[station_name]; !exists {
			station_name_idx[station_name] = current_station_idx
			mean[current_station_idx] = temperature
			min[current_station_idx] = temperature
			max[current_station_idx] = temperature
			current_station_idx++
		}
		var station_idx uint16 = station_name_idx[station_name]
		var new_min, new_mean, new_max int16 = 
			compute_min_mean_max(min[station_idx], mean[station_idx], max[station_idx], temperature)
		min[station_idx] = new_min
		mean[station_idx] = new_mean
		max[station_idx] = new_max
	}

	for station_name, station_idx := range station_name_idx {
		fmt.Printf("%s\n", build_output(station_name, min[station_idx], mean[station_idx], max[station_idx]))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Now().Sub(start_time)
	fmt.Printf("Runtime: %s \n", elapsed.String())
}
