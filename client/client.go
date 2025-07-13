package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Exchange struct {
	Rate string `json:"bid"`
}

func main() {
	slog.Info("Starting client...")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		slog.Error("Error consuming server API")
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("Request timeout to server API")
			panic(err)
		}

		slog.Error(err.Error())
		panic(err)
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error receiving response from server API")
		panic(err)
	}
	var exchange Exchange
	if err := json.Unmarshal(res, &exchange); err != nil {
		slog.Error("Error processing exchange rate api response")
		panic(err)
	}
	slog.Info("Exchange rate received successfully", "rate", exchange.Rate)
	if err := saveExchangeRate(exchange.Rate); err != nil {
		slog.Error("Error to save exchange rate")
		panic(err)
	}
	slog.Info("Exchange rate saved successfully")
}

func saveExchangeRate(exchangeRateValue string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		slog.Error("Error to create file to save exchange rate")
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("Dólar: %s\n", exchangeRateValue))
	if err != nil {
		fmt.Printf("Erro gravando o arquivo: %v", err)
		panic(err)
	}
	return nil
}
