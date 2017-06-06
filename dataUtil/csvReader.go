package dataUtil

import (
	"time"
	"strings"
	"strconv"
	"fmt"
	"github.com/wcharczuk/go-chart/util"
)

func ReadDataSimple(name string) ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64
	err := util.File.ReadByLines("./curseData/"+name, func(line string) error {
		parts := strings.Split(line, ",")
		download := parseFloat64(parts[1])
		timeStr := parts[0]
		timeSplit := strings.Split(timeStr, " ")
		dateSplit := strings.Split(timeSplit[0], "/")
		hourSplit := strings.Split(timeSplit[1], ":")
		day, _ := strconv.Atoi(dateSplit[0])
		month, _ := strconv.Atoi(dateSplit[1])
		year, _ := strconv.Atoi(dateSplit[2])
		hour, _ := strconv.Atoi(hourSplit[0])
		min, _ := strconv.Atoi(hourSplit[1])
		xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC))
		yvalues = append(yvalues, download)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return xvalues, yvalues
}

func ReadDataHour(name string) ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64
	var timeCount = 0
	var lastDlCount float64
	var firstCount = true
	var lastValue float64 = 0
	err := util.File.ReadByLines("./curseData/"+name, func(line string) error {
		timeCount ++
		if timeCount == 4 * 24  {
			fmt.Println(line)
			parts := strings.Split(line, ",")
			download := parseFloat64(parts[1])
			if firstCount {
				lastDlCount = download
				firstCount = false
				timeCount = 0
				return nil
			}
			timeStr := parts[0]
			timeSplit := strings.Split(timeStr, " ")
			dateSplit := strings.Split(timeSplit[0], "/")
			hourSplit := strings.Split(timeSplit[1], ":")
			day, _ := strconv.Atoi(dateSplit[0])
			month, _ := strconv.Atoi(dateSplit[1])
			year, _ := strconv.Atoi(dateSplit[2])
			hour, _ := strconv.Atoi(hourSplit[0])
			min, _ := strconv.Atoi(hourSplit[1])
			xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC))
			yvalues = append(yvalues, download - lastDlCount)
			lastValue = download

			lastDlCount = download
			timeCount = 0
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return xvalues, yvalues
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}
