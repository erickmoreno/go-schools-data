package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	var host = "https://api.educamaisbrasil.com.br"
	// joao-pessoa paraiba
	// recife pernambuco
	// campina-grande paraiba
	var estado = "paraiba"
	var cidade = "campina-grande"
	amount := getTotalEscolas(host, cidade, estado)
	escolas := getEscolas(host, cidade, estado, amount)
	var detalhes []DetalhesEscola

	for _, element := range escolas.Escolas {
		detalhes = append(detalhes, getEscolaDetails(host, element.Url)[0])
	}

	//saving csv file
	csvFile, _ := os.Create(fmt.Sprintf("%s-%s.csv", estado, cidade))
	gocsv.MarshalFile(&detalhes, csvFile) // Get all clients as CSV string
}

func getCidade(host string, cidade string, uf string) Cidade {
	var getCidade = "/api/Instituicao/ConsultarCidadePorUrlEUF"
	var cmd = host + getCidade

	params := "urlCidade=" + url.QueryEscape(cidade) + "&" + "uf=" + url.QueryEscape(uf)
	path := fmt.Sprintf("%s?%s", cmd, params)

	resp, err := http.Get(path)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to Cidade struct
	var cidadeStruct Cidade
	json.Unmarshal(bodyBytes, &cidadeStruct)

	return cidadeStruct
}

func getTotalEscolas(host string, cidade string, uf string) int {
	escolas := getEscolas(host, cidade, uf, 0)
	return escolas.Total
}

func getEscolas(host string, cidade string, uf string, amount int) EscolasLista {
	var listaEscolas = "/api/Instituicao/ListarEscolas"
	var cmd = host + listaEscolas
	params := "cidadeUrl=" + url.QueryEscape(cidade) + "&" +
		"estadoUrl=" + url.QueryEscape(uf) + "&" +
		"fimPaginacao=" + url.QueryEscape(strconv.Itoa(amount)) + "&" +
		"inicioPaginacao=" + url.QueryEscape("0") + "&" +
		"temBolsa=" + url.QueryEscape("false") + "&" +
		"filtroPublica=" + url.QueryEscape("E\\,F\\,M") + "&" +
		"filtroPrivado=" + url.QueryEscape("P") + "&" +
		"temPaginacao=" + url.QueryEscape("false")

	path := fmt.Sprintf("%s?%s", cmd, params)

	resp, err := http.Get(path)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to Cidade struct
	var escolas EscolasLista
	json.Unmarshal(bodyBytes, &escolas)

	return escolas
}

func getEscolaDetails(host string, urlEscola string) []DetalhesEscola {
	var listaEscolas = "/api/Instituicao/ConsultarDadosEscola"
	var cmd = host + listaEscolas
	params := "urlEscola=" + url.QueryEscape(urlEscola)

	path := fmt.Sprintf("%s?%s", cmd, params)

	resp, err := http.Get(path)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to Cidade struct
	var detalhe []DetalhesEscola
	json.Unmarshal(bodyBytes, &detalhe)

	return detalhe
}

type DetalhesEscola struct {
	Cidade   string `json:"CIDADE"`
	Nome     string `json:"NOME"`
	Telefone string `json:"TELEFONE"`
	Email    string `json:"EMAIL"`
	Bairro   string `json:"BAIRRO"`
}

type Cidade struct {
	IdCidade      string `json:"ID_CIDADE"`
	Nome          string `json:"NOME"`
	Uf            string `json:"UF"`
	TipoCidade    string `json:"TIPO_CIDADE"`
	QtdHabitantes string `json:"QTD_HABITANTE"`
	NomeCustom    string `json:"NOME_CUSTOM"`
	CodIbge       string `json:"COD_IBGE"`
}

type EscolaBase struct {
	Nome   string `json:"NOME"`
	Cidade string `json:"CIDADE"`
	Bairro string `json:"BAIRRO"`
	Uf     string `json:"UF"`
	Tipo   string `json:"TIPO"`
	Url    string `json:"URL"`
}

type EscolasLista struct {
	Escolas    []EscolaBase `json:"escolas"`
	Total      int          `json:"escolas_qtd"`
	Federais   int          `json:"escolas_qtd_federal"`
	Municipais int          `json:"escolas_qtd_municipal"`
	Estaduais  int          `json:"escolas_qtd_estadual"`
	Privadas   int          `json:"escolas_qtd_privada"`
}
