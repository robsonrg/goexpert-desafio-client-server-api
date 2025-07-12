package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/robsonrg/goexpert-desafio-client-server-api/server/database"
)

type Exchange struct {
	USDBRL struct {
		Code       string `json:"-"`
		Codein     string `json:"-"`
		Name       string `json:"-"`
		High       string `json:"-"`
		Low        string `json:"-"`
		VarBid     string `json:"-"`
		PctChange  string `json:"-"`
		Bid        string `json:"bid"`
		Ask        string `json:"-"`
		Timestamp  string `json:"-"`
		CreateDate string `json:"-"`
	} `json:"USDBRL"`
}

func main() {
	slog.Info("Starting server...")

	db := database.DB.GetDB()
	defer db.Close()

	http.HandleFunc("/cotacao", getExchangeRateHandler)

	slog.Info("Listen requests on 8080...")
	http.ListenAndServe(":8080", nil)
}

func getExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Start handling request")
	defer slog.Info("End handling request")
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		slog.Error("Error consuming exchange rate API")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Timeout when consuming the exchange rate API")
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error receiving response from exchange rate api")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var exchange Exchange
	if err := json.Unmarshal(res, &exchange); err != nil {
		slog.Error("Error processing exchange rate api response")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := saveExchangeRate(exchange.USDBRL.Bid); err != nil {
		slog.Error("Error to save exchange rate")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(exchange.USDBRL)
}

func saveExchangeRate(exchangeRateValue string) error {
	db := database.DB.GetDB()
	stmt, err := db.Prepare(`INSERT INTO exchanges (rate_value) VALUES (?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, err := stmt.ExecContext(ctx, exchangeRateValue); err != nil {
		log.Printf("Error to insert exchange rate into database: %v", err)
		return err
	}
	slog.Info("Exchange rate inserted into database")
	return nil
}
