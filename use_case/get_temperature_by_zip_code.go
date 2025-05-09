package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com.br/sk8sta13/temperatures/internal/dto"
	"github.com.br/sk8sta13/temperatures/internal/entity"
)

type Address struct {
	Code         string `json:"cep"`
	State        string `json:"estado"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
	Region       string `json:"regiao"`
}

func Get(data *dto.ZipCode) (*dto.Temperature, error) {
	address, err := getLocal(data.ZipCode)
	if err != nil {
		return nil, err
	}

	temperature, err := getTemperature(&address.City)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}

func getLocal(zipcode string) (*Address, error) {
	resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode))
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}
	defer resp.Body.Close()

	var a Address
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	if a.City == "" {
		return nil, entity.ErrCanNotFindZipcode
	}

	return &a, nil
}

func getTemperature(region *string) (*dto.Temperature, error) {
	regionScape := url.QueryEscape(*region)
	resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", regionScape, os.Getenv("WEATHER_API_KEY")))
	if err != nil {
		println("a")
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}
	defer resp.Body.Close()

	var d struct {
		Current dto.Temperature `json:"current"`
	}

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	d.Current.Temp_K = d.Current.Temp_C + 273.15

	return &d.Current, nil
}
