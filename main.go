package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	cep       = "01153000"
	brasilAPI = "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCEP    = "http://viacep.com.br/ws/" + cep + "/json/"
	timeout   = 1 * time.Second
)

type BrasilAPIAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCEPAddress struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type UnifiedAddress struct {
	Cep          string
	Street       string
	Neighborhood string
	City         string
	State        string
}

func fetchBrasilAPI(ch chan<- map[string]interface{}) {
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(brasilAPI)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var address BrasilAPIAddress
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return
	}

	unifiedAddress := UnifiedAddress{
		Cep:          address.Cep,
		Street:       address.Street,
		Neighborhood: address.Neighborhood,
		City:         address.City,
		State:        address.State,
	}

	ch <- map[string]interface{}{
		"api":     "BrasilAPI",
		"address": unifiedAddress,
	}
}

func fetchViaCEP(ch chan<- map[string]interface{}) {
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(viaCEP)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var address ViaCEPAddress
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return
	}

	unifiedAddress := UnifiedAddress{
		Cep:          address.Cep,
		Street:       address.Logradouro,
		Neighborhood: address.Bairro,
		City:         address.Localidade,
		State:        address.Uf,
	}

	ch <- map[string]interface{}{
		"api":     "ViaCEP",
		"address": unifiedAddress,
	}
}

func main() {
	ch := make(chan map[string]interface{}, 2)

	go fetchBrasilAPI(ch)
	go fetchViaCEP(ch)

	select {
	case res := <-ch:
		address := res["address"].(UnifiedAddress)
		fmt.Println("----")
		fmt.Println("API:", res["api"])
		fmt.Println("Endereço:")
		fmt.Println("  CEP:", address.Cep)
		fmt.Println("  Logradouro:", address.Street)
		fmt.Println("  Bairro:", address.Neighborhood)
		fmt.Println("  Cidade:", address.City)
		fmt.Println("  Estado:", address.State)
	case <-time.After(timeout):
		fmt.Println("Timeout waiting response")
	}

}
