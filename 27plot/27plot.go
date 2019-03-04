package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Output struct {
	Now                     string
	NowTime                 time.Time
	Date                    string `json:"forecast_date"`
	Min                     int    `json:"temperature_min"`
	Max                     int    `json:"temperature_max"`
	Precipitation           int    `json:"precipitation_amount"`
	PrecipitationPercentage int    `json:"precipitation_percentage"`
	WindDirection           string `json:"wind_direction"`
	WindSpeed               int    `json:"wind_speed_bft"`
}

func main() {
	// 2019-02-21 11:00:02.651709762
	//layout := "2014-09-12 11:45:26.123"
	//layout := "2014-09-12 "

	items := linecount("27.log") - 1
	file, err := os.Open("27.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	min := make(plotter.XYs, items)
	max := make(plotter.XYs, items)
	precipationPerc := make(plotter.XYs, items)
	WindSpeed := make(plotter.XYs, items)

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		o := &Output{}
		t := scanner.Text()
		if strings.Index(t, "{") != 0 {
			continue
		}
		err := json.Unmarshal([]byte(t), o)
		if err != nil {
			log.Fatal(err)
		}
		t3 := strings.Split(o.Now, " ")
		t3 = t3[:len(t3)-1]
		ti := strings.Join(t3, " ")

		o.NowTime, _ = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", ti)

		min[i].Y = float64(o.Min)
		min[i].X = float64(o.NowTime.Unix())

		max[i].Y = float64(o.Max)
		max[i].X = float64(o.NowTime.Unix())

		precipationPerc[i].Y = float64(o.PrecipitationPercentage)
		precipationPerc[i].X = float64(o.NowTime.Unix())

		WindSpeed[i].Y = float64(o.WindSpeed)
		WindSpeed[i].X = float64(o.NowTime.Unix())

		log.Printf("Min: %+v", min)
		//fmt.Println(scanner.Text())
		log.Printf("output: %+v", o)
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	log.Printf("Adding Plot\n")
	err = plotutil.AddLinePoints(p,
		"Min", min,
		"Max", max,
		"Windspeed", WindSpeed,
		"PrecipationPerc", precipationPerc)
	//"Second", randomPoints(15),
	//"Third", randomPoints(15))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	log.Printf("Writing to file points.png\n")
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

func linecount(file string) int {
	fileh, _ := os.Open(file)
	fileScanner := bufio.NewScanner(fileh)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}
