package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Parcela struct {
	Numero         int       `json:"numero"`
	Valor          float64   `json:"valor"`
	DataVencimento time.Time `json:"data_vencimento"`
	Entrada        bool      `json:"entrada,omitempty"`
	Somatoria      float64   `json:"somatoria"`
}

type Carne struct {
	ID                     int       `json:"id"`
	ValorTotal             float64   `json:"valor_total"`
	QtdParcelas            int       `json:"qtd_parcelas"`
	DataPrimeiroVencimento string    `json:"data_primeiro_vencimento"`
	Periodicidade          string    `json:"periodicidade"`
	ValorEntrada           float64   `json:"valor_entrada,omitempty"`
	Parcelas               []Parcela `json:"parcelas"`
}

var (
	carneStore = make(map[int]Carne)
	currentID  = 1
	storeMutex sync.Mutex
)

func calcularParcelas(carne *Carne) {
	valorRestante := carne.ValorTotal - carne.ValorEntrada
	valorParcela := valorRestante / float64(carne.QtdParcelas)
	dataVencimento, _ := time.Parse("2006-01-02", carne.DataPrimeiroVencimento)

	somatoria := carne.ValorEntrada // Inicializa com o valor de entrada, se houver

	if carne.ValorEntrada > 0 {
		carne.Parcelas = append(carne.Parcelas, Parcela{
			Numero:         1,
			Valor:          carne.ValorEntrada,
			DataVencimento: time.Now(),
			Entrada:        true,
			Somatoria:      somatoria,
		})
	}

	for i := 1; i <= carne.QtdParcelas; i++ {
		somatoria += valorParcela
		carne.Parcelas = append(carne.Parcelas, Parcela{
			Numero:         i,
			Valor:          valorParcela,
			DataVencimento: dataVencimento,
			Somatoria:      somatoria,
		})
		if carne.Periodicidade == "mensal" {
			dataVencimento = dataVencimento.AddDate(0, 1, 0)
		} else if carne.Periodicidade == "semanal" {
			dataVencimento = dataVencimento.AddDate(0, 0, 7)
		}
	}
}

func criarCarne(w http.ResponseWriter, r *http.Request) {
	var carne Carne
	_ = json.NewDecoder(r.Body).Decode(&carne)

	storeMutex.Lock()
	carne.ID = currentID
	currentID++
	storeMutex.Unlock()

	calcularParcelas(&carne)

	storeMutex.Lock()
	carneStore[carne.ID] = carne
	storeMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carne)
}

func recuperarParcelas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	storeMutex.Lock()
	carne, ok := carneStore[id]
	storeMutex.Unlock()

	if !ok {
		http.Error(w, "Carnê não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carne.Parcelas)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/carne", criarCarne).Methods("POST")
	r.HandleFunc("/carne/{id}/parcelas", recuperarParcelas).Methods("GET")

	fmt.Println("API rodando em :8080")
	http.ListenAndServe(":8080", r)
}
