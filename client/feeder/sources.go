package feeder

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Price struct {
	Denom string
	TargetPrice float64
	CurrentPrice float64
}

type Source interface {
	getData() ([]Price, error)
}

type FixedSource struct {
	prices []Price
}

func (fs FixedSource) getData() ([]Price , error) {
	return fs.prices, nil
}

func CreateJsonSource(priceDataJson string) (*FixedSource, error) {
	var prices []Price

	err := json.Unmarshal([]byte(priceDataJson), &prices)
	if err != nil {
		return nil, err
	}
	return &FixedSource{prices}, nil
}


type URLSource struct {
	url string
}

func (src URLSource) getData() (prices []Price , err error)  {
	resp, err := http.Get(src.url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &prices)
	return
}

type FileSource struct {
	filename string
}

func (src FileSource) getData() (prices []Price , err error)  {
	content, err := ioutil.ReadFile(src.filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &prices)
	return
}
