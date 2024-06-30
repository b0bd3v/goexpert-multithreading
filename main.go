package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	urlViaCep  = "https://viacep.com.br/ws/%s/json"
	urlFindCep = "https://brasilapi.com.br/api/cep/v1/%s.json"
)

func main() {
	var viaCEPChannel = make(chan interface{})
	var brasilAPIChannel = make(chan interface{})

	var viaCEP ViaCEP
	go FetchAndDecode(urlViaCep, &viaCEP, viaCEPChannel)

	var brasilAPI BrasilAPI
	go FetchAndDecode(urlFindCep, &brasilAPI, brasilAPIChannel)

	select {
	case response := <-viaCEPChannel:
		prettyJSON, _ := json.MarshalIndent(response, "", "    ")
		fmt.Printf("ViaCEP: %s\n", string(prettyJSON))
	case response := <-brasilAPIChannel:
		prettyJSON, _ := json.MarshalIndent(response, "", "    ")
		fmt.Printf("BrasilAPI: %s\n", string(prettyJSON))
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

	fmt.Printf("Done")
}

func FetchAndDecode(url string, target interface{}, channel chan<- interface{}) {
	resp, err := http.Get(fmt.Sprintf(url, "91450147"))

	if err != nil {
		panic("failed to fetch data")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("failed to fetch data")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic("failed to read data")
	}

	err = json.Unmarshal(body, target)

	if err != nil {
		panic("failed to decode data")
	}

	channel <- target
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}
