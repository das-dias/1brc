package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

const MAX_STATIONS uint32 = 100000
const MAX_STATION_NAME_LENGTH uint8 = 100
const MAX_TEMP_VALUE_LENGTH uint8 = 5

func parse_station_name(line string) (uint8, string) {
	var i uint8 = 0
	temp_station_name := [MAX_STATION_NAME_LENGTH]byte{0}
	for ; (i < MAX_STATION_NAME_LENGTH) && (line[i] != ';'); i++ {
		temp_station_name[i] = line[i]
	}
	return i + 1, string(temp_station_name[:])
}

func parse_station_temperature(line string) int16 {
	var i uint8 = 0
	temperature_buffer := [MAX_TEMP_VALUE_LENGTH]byte{0}
	for ; (i < MAX_TEMP_VALUE_LENGTH) && (i < uint8(len(line))); i++ {
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
	for ; (i < MAX_TEMP_VALUE_LENGTH) && (temperature_buffer[i] != 0); i++ {
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
	} else { // only run if it is not minimum, other wise it is a waste of time
		if new_temperature > max {
			new_max = new_temperature
		}
	}

	return new_min, new_mean, new_max
}

func parse_line(line string) (string, int16) {
	line_idx, station_name := parse_station_name(line)
	var station_temperature int16 = parse_station_temperature(line[line_idx:])
	return station_name, station_temperature
}

func build_output_string(station_name string, min int16, mean int16, max int16) string {
	float_min := float32(min) / 10.0
	float_max := float32(max) / 10.0
	float_mean := float32(mean) / 10.0
	return fmt.Sprintf("%s;%.1f;%.1f;%.1f", station_name, float_min, float_mean, float_max)
}

func main() {

	input_file := flag.String("i", "", "Input file")
	cpu_profiling := flag.String("cp", "", "Output cpu profile file for time resource utilisation report")
	flag.Parse()

	if *cpu_profiling != "" {
		fp, err := os.Create(*cpu_profiling)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
		defer pprof.StopCPUProfile()
	}

	if *input_file == "" {
		log.Fatal("Input file is required")
	}

	ifp, err := os.Open(*input_file)
	if err != nil {
		log.Fatal(err)
	}
	defer ifp.Close()

	start_time := time.Now()
	min := [MAX_STATIONS]int16{0}
	max := [MAX_STATIONS]int16{0}
	mean := [MAX_STATIONS]int16{0}
	var current_station_idx uint16 = 0
	station_name_idx := make(map[string]uint16)

	scanner := bufio.NewScanner(ifp)
	for scanner.Scan() {
		line := scanner.Text()
		station_name, temperature := parse_line(line)
		if _, exists := station_name_idx[station_name]; !exists {
			station_name_idx[station_name] = current_station_idx
			mean[current_station_idx] = temperature
			min[current_station_idx] = temperature
			max[current_station_idx] = temperature
			current_station_idx++
		}
		var station_idx uint16 = station_name_idx[station_name]
		var new_min, new_mean, new_max int16 = compute_min_mean_max(min[station_idx], mean[station_idx], max[station_idx], temperature)
		min[station_idx] = new_min
		mean[station_idx] = new_mean
		max[station_idx] = new_max
	}

	for station_name, station_idx := range station_name_idx {
		fmt.Printf("%s\n", build_output_string(station_name, min[station_idx], mean[station_idx], max[station_idx]))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Now().Sub(start_time)
	fmt.Printf("Runtime: %s \n", elapsed.String())
}
