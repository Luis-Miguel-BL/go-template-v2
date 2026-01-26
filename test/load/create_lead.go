package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/request"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var token = "your-access-token"

func main() {

	os.MkdirAll("test/load/reports", 0755)

	payload := request.CreateLeadRequest{
		Name:           "luis miguel",
		Email:          "2001lmbl@gmail.com",
		Phone:          "37998153343",
		DocumentNumber: "14064435656",
	}

	body, _ := json.Marshal(payload)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:3000/v1/leads",
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + token},
		},
		Body: body,
	})

	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 2 * time.Minute
	attacker := vegeta.NewAttacker()

	resultsFile, _ := os.Create("test/load/reports/results.bin")
	defer resultsFile.Close()

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "create-lead") {
		metrics.Add(res)
		resultsFile.Write(res.Body)
	}

	metrics.Close()

	reportFile, _ := os.Create("test/load/reports/report.json")
	defer reportFile.Close()

	vegeta.NewJSONReporter(&metrics).Report(reportFile)
}
