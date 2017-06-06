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
	"github.com/patrickmn/go-cache"
	"github.com/modmuss50/CurseMapper/dataUtil"
	"github.com/wcharczuk/go-chart/drawing"
)

const NewLine = "\n"
var pngCache  = cache.New(15*time.Minute, 15*time.Minute)

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html lang=\"en\"><body>"+NewLine)
	io.WriteString(w, "<p>Curse mapper</p>")

	if r.FormValue("dlhour") == "true" {
		io.WriteString(w, "Showing downloads per hour")
	}
	files, _ := ioutil.ReadDir("./curseData")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_export.csv") {
			username := strings.Replace(f.Name(), "_export.csv", "", -1)
			io.WriteString(w, "<p><a href=\"/user/"+username+"\">View all projects by "+username+"</a></p>"+NewLine)
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
			projectName := strings.Replace(f.Name(), "_export.csv", "", -1)
			io.WriteString(w, "<p>"+projectName+"</p>"+NewLine)
			io.WriteString(w, "<img src=\"/chart/projects/"+username+"/"+projectName+"\">"+NewLine)
		}

	}
	io.WriteString(w, "</body></body>")
}

func drawChart(w http.ResponseWriter, r *http.Request) {
	png, found := pngCache.Get(r.URL.String())
	if found {
		fmt.Println("Using cache copy of :"+ r.URL.String())
		w.Header().Set("Content-Type", chart.ContentTypeSVG)
		w.Write(png.([]byte))
		return
	}
	fmt.Println("generating new copy of :"+ r.URL.String())

	fmt.Println(r.URL.String())
	username := strings.Replace(r.URL.String(), "/chart/", "", -1)

	xvalues, yvalues := dataUtil.ReadDataHour(username + "_export.csv")
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
	linRegSeries := &chart.LinearRegressionSeries{
		Name: "Liner Average",
		InnerSeries: mainSeries,
	}
	smaSeries := &chart.SMASeries{
		Name: "SM Average",
		InnerSeries: mainSeries,
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
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
			Name: "Date",
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		Series: []chart.Series{
			mainSeries,
			maxSeries,
			//chart.LastValueAnnotation(maxSeries),
			linRegSeries,
			smaSeries,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.SVG, buffer)
	w.Header().Set("Content-Type", chart.ContentTypeSVG)
	pngCache.Set(r.URL.String(), buffer.Bytes(), cache.DefaultExpiration)
	w.Write(buffer.Bytes())
}



func main() {
	http.HandleFunc("/chart/", drawChart)
	http.HandleFunc("/user/", userpage)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
	fmt.Println("Server started at http://localhost:8000")
}

