package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Output struct {
	Now                     string
	Date                    string `json:"forecast_date"`
	Min                     int    `json:"temperature_min"`
	Max                     int    `json:"temperature_max"`
	Precipitation           int    `json:"precipitation_amount"`
	PrecipitationPercentage int    `json:"precipitation_percentage"`
	WindDirection           string `json:"wind_direction"`
	WindSpeed               int    `json:"wind_speed_bft"`
}

func main() {
	resp, err := http.Get("https://www.weeronline.nl/Europa/Nederland/Schoorl/4058246")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading body: %s", err)
		return
	}
	//body, _ := ioutil.ReadFile("index.html")

	r, err := regexp.Compile("2019-02-26.*({.*2019-02-27T.+?})")
	if err != nil {
		panic(err)
	}

	res := r.FindStringSubmatch(string(body))
	/*
		for id, re := range res {
			log.Printf("id: %d re: %s", id, re)
		}
	*/

	w := &Output{}
	err = json.Unmarshal([]byte(res[1]), w)
	if err != nil {
		panic(err)
	}
	w.Now = time.Now().String()

	//log.Printf("Weer: %+v", w)

	wJSON, _ := json.Marshal(w)
	fmt.Printf("%s\n", wJSON)

	/*
		for i := 0; i < 20; i++ {
			i := strings.Index(string(body), "{")
			body = body[i+1:]
		}
	*/
	//log.Printf("Body: %s", body)

}
