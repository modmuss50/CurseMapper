package main

import (
	"io"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"github.com/wcharczuk/go-chart"
	"time"
	"strings"
	"github.com/wcharczuk/go-chart/util"
	"strconv"
	"github.com/patrickmn/go-cache"
)

const NewLine = "\n"
var pngCache  = cache.New(15*time.Minute, 15*time.Minute)

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html lang=\"en\"><body>"+NewLine)
	files, _ := ioutil.ReadDir("./curseData")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_export.csv") {
			username := strings.Replace(f.Name(), "_export.csv", "", -1)
			io.WriteString(w, "<p><a href=\"/user/"+username+"\">View all projects by "+username+"</a></p>"+NewLine)
			io.WriteString(w, "<img src=\"/chart/"+f.Name()+"\">"+NewLine)
		}

	}
	io.WriteString(w, "</body></body>")
}

func userpage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html lang=\"en\"><body>"+NewLine)
	io.WriteString(w, "<p><a href=\"/\">Go home</a></p>"+NewLine)
	fmt.Println(r.URL)
	username := strings.Replace(r.URL.String(), "/user/", "", -1)
	fmt.Println(username)
	files, _ := ioutil.ReadDir("./curseData/projects/" + username)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_export.csv") {
			io.WriteString(w, "<p>"+f.Name()+"</p>"+NewLine)
			io.WriteString(w, "<img src=\"/chart/projects/"+username+"/"+f.Name()+"\">"+NewLine)
		}

	}
	io.WriteString(w, "</body></body>")
}

func drawChart(w http.ResponseWriter, r *http.Request) {
	png, found := pngCache.Get(r.URL.String())
	if found {
		fmt.Println("Using cache copy of :"+ r.URL.String())
		w.Write(png.([]byte))
		return
	}
	fmt.Println("generating new copy of :"+ r.URL.String())

	fmt.Println(r.URL.String())
	username := strings.Replace(r.URL.String(), "/chart/", "", -1)

	xvalues, yvalues := readData(username)
	mainSeries := chart.TimeSeries{
		Name: "Downloads",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlue,
		},
		XValues: xvalues,
		YValues: yvalues,
	}
	maxSeries := &chart.MaxSeries{
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorAlternateGray,
		},
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Title:  "Total Downloads for <USER>",
		Width:  1280,
		Height: 720,
		YAxis: chart.YAxis{
			Name:      "Total Downloads",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d dl", int(v.(float64)))
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		Series: []chart.Series{
			mainSeries,
			maxSeries,
			chart.LastValueAnnotation(maxSeries),
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)
	w.Header().Set("Content-Type", chart.ContentTypePNG)
	pngCache.Set(r.URL.String(), buffer.Bytes(), cache.DefaultExpiration)
	w.Write(buffer.Bytes())
}

func readData(name string) ([]time.Time, []float64) {
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

func main() {
	http.HandleFunc("/chart/", drawChart)
	http.HandleFunc("/user/", userpage)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}
