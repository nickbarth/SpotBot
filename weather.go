package main

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"log"
	"net/http"
)

type WeatherJSON struct {
	ForecastGroup struct {
		Forecast []struct {
			Period      string `xml:"period"`
			TextSummary string `xml:"textSummary"`
		} `xml:"forecast"`
	} `xml:"forecastGroup"`
}

type Weather struct {
	url string
}

func (w Weather) Get() string {
	resp, err := http.Get(w.url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// body = []byte(`<siteData> <license>test</license> <forecastGroup><forecast> <period>Monday</period> <textSummary>A mix of sun and cloud. High 11. UV index 3 or moderate.</textSummary></forecast></forecastGroup> </siteData>`)

	var weather WeatherJSON
	// err = xml.Unmarshal(body, &weather)
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&weather)

	if err != nil {
		log.Fatal(err)
	}

	return weather.ForecastGroup.Forecast[0].TextSummary
}
