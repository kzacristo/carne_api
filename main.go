package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Parcelas struct {
	Numero         int       `json:"numero"`
	Valor          float64   `json:"valor"`
	DataVencimento time.time `json:"data_vencimento"`
	Entrada        bool      `json:"entrada,omitempty"`
}

type Carne struct {
	ValorTotal             float64    `json:"valor_total"`
	QtdParcelas            int        `json:"qtd_parcelas"`
	DataPrimeiroVencimento string     `json:"data_primeiro_vencimento"`
	Periodicidade          string     `json:"periodicidade"`
	ValorEntrada           float64    `json:"valor_entrada,omitempty"`
	Parcelas               []Parcelas `jdon:"parcelas"`
}


func calcularParcelas(carne *Carne){
      valorRestante := carne.ValorTotal - carne.valor_entrada
	  valorParcelas  := valorRestante / float64(carne.QtdParcelas)
	  dataVencimento, _ := time.Parse("2006-01-02", carne.DataPrimeiroVencimento)

	  for i := 1, i <= carne.QtdParcelas; i++{
		carne.Parcelas = append(carne.Parcelas, Parcela{
			Numero: 1,
			Valor: valorParcelas,
			DataVencimento: dataVencimento
		})

		if carne.Periodicidade == "mensal"{
			dataVencimento = dataVencimento.AddDate(0,1,0)
		}
		else if carne.Periodicidade == "semanal" {
			dataVencimento = dataVencimento.AddDate(0,0,7)
		}
	  }
	if carne.ValorEntrada > 0 {
		carne.Parcelas = append([]Parcela{{
			Numero: 1,
			Valor: carne.ValorEntrada,
			DataVencimento: time.Now(),
			Entrada: true,
		}}, carne.Parcelas...)
	}
}

func criarCarne(w http.ResponseWriter, r *http.Request){
	var carne Carne
	_= json.NewDecoder(r.body).Decode(&carne)

	calcularParcelas(&carne)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carne)
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/carne", criarCarne).Methods("POST
	fmt.Println("API rodando em :8080")
	http.ListenAndServe(":8080", r)
}