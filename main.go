package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Address struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	API          string `json:"api"`
}

type viaCepResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"estado"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
}

func main() {

	viaCepAddress := make(chan Address)
	brasilApiAddress := make(chan Address)

	cep := "13503100"

	go func() {
		viaCepAddress <- searchWithViaCep(cep)
	}()

	go func() {
		brasilApiAddress <- searchWithBrasilApi(cep)
	}()

	select {
	case address := <-viaCepAddress:
		printAddress(address)
	case address := <-brasilApiAddress:
		printAddress(address)
	case <-time.After(time.Second):
		println("Timeout")
	}

}

func printAddress(address Address) {
	println("CEP:", address.CEP)
	println("State:", address.State)
	println("City:", address.City)
	println("Neighborhood:", address.Neighborhood)
	println("Street:", address.Street)
	println("API:", address.API)
}

func convertViaCepResponseToAddress(response viaCepResponse) Address {

	return Address{
		CEP:          response.CEP,
		State:        response.State,
		City:         response.City,
		Neighborhood: response.Neighborhood,
		Street:       response.Street,
		API:          "Via CEP",
	}
}

func searchWithViaCep(cep string) Address {

	viaCepApiUrl := "https://viacep.com.br/ws/" + cep + "/json/"

	response, err := http.Get(viaCepApiUrl)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var viaCepResponse viaCepResponse

	err = json.NewDecoder(response.Body).Decode(&viaCepResponse)

	if err != nil {
		panic(err)
	}

	return convertViaCepResponseToAddress(viaCepResponse)
}

func searchWithBrasilApi(cep string) Address {

	brasilApiUrl := "https://brasilapi.com.br/api/cep/v1/" + cep

	response, err := http.Get(brasilApiUrl)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var address Address

	err = json.NewDecoder(response.Body).Decode(&address)

	if err != nil {
		panic(err)
	}

	address.API = "Brasil API"

	return address
}
