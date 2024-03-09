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

func main() {

	start_time := time.Now()

	min := [MAX_STATIONS]int16{0}
	max := [MAX_STATIONS]int16{0}
	mean := [MAX_STATIONS]int16{0}
	var current_station_idx uint16 = 0
	station_name_idx := make(map[string]uint16)
	temp_station_name := [MAX_STATION_NAME_LENGTH]byte{0}
	station_temperature_buffer := [MAX_TEMP_VALUE_LENGTH]byte{0} // string value = '-99.9' ; '99.9'
	fp, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		var i uint8 = 0
		// read station's name
		for ; line[i] != ';'; i++ {
			temp_station_name[i] = line[i]
		}
		i++ // ignore the ';' utf-8 char and read station's temperature byte representation
		var k uint8 = 0
		for _, ch := range line[i:] {
			station_temperature_buffer[k] = byte(ch)
			k++
		}
		// convert the temperature bytes to its integer representation, considering the decimal parts ar integers as well
		var uint_temp uint16 = 0
		k = 0
		if station_temperature_buffer[0] == '-' {
			k = 1
		}
		// convert the temperature to an integer
		for ; station_temperature_buffer[k] != '.'; k++ {
			uint_temp = (uint_temp * 10) + uint16(station_temperature_buffer[k]-'0')
		}
		k++ // ignore the '.' separation point
		for ; (k < MAX_TEMP_VALUE_LENGTH) && (station_temperature_buffer[k] != 0); k++ {
			uint_temp = (uint_temp * 10) + uint16(station_temperature_buffer[k]-'0')
		}
		// setup the sign bits
		var int_temp int16 = int16(uint_temp)
		if station_temperature_buffer[0] == '-' {
			int_temp = -int_temp
		}
		sname := string(temp_station_name[:])
		if _, exists := station_name_idx[sname]; !exists {
			station_name_idx[sname] = current_station_idx
			mean[current_station_idx] = int_temp
			min[current_station_idx] = int_temp
			max[current_station_idx] = int_temp
			current_station_idx++
		}
		var station_idx uint16 = station_name_idx[sname]
		// get the float_temp with only the mantissa to perform arithmetic operations
		temp_station_name = [MAX_STATION_NAME_LENGTH]byte{0}
		station_temperature_buffer = [MAX_TEMP_VALUE_LENGTH]byte{0}
		var sum int16 = mean[station_idx] + int_temp
		mean[station_idx] = sum >> 1

		if int_temp < min[station_idx] {
			min[station_idx] = int_temp
		} else { //if int_temp > max[station_idx] { // if it is not the new minimum it is the new maximum or equal
			max[station_idx] = int_temp
		}
	}

	for sname, sidx := range station_name_idx {
		min := float32(int16(min[sidx])) / 10.0
		max := float32(int16(max[sidx])) / 10.0
		mean := float32(int16(mean[sidx])) / 10.0
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", sname, min, mean, max)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Now().Sub(start_time)
	fmt.Printf("Runtime: %s \n", elapsed.String())
}
