package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Target struct {
	ConversionRates map[string]interface{} `json:"conversion_rates"`
}

func Exchange(cCtx *cli.Context) error {
	if cCtx.Args().Len() < 3 {
		return fmt.Errorf("\nerror: not enough arguments\nxchange [amount] [initial currency code] [first final currency code] ... [N final currency code]")
	}
	initialAmount, _ := strconv.ParseFloat(cCtx.Args().Get(0), 64)
	currencyNames := cCtx.Args().Slice()[1:]
	for i, s := range currencyNames {
		currencyNames[i] = strings.ToUpper(s)
	}
	currencyValues := GetAPIResult(currencyNames, initialAmount)
	table := ConvertToTable(currencyNames, currencyValues)
	fmt.Println(table)
	return nil
}
func GetAPIResult(currencies []string, amount float64) []float64 {
	initialCurrency := currencies[0]
	var results []float64

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	var url = "https://v6.exchangerate-api.com/v6/" + cfg.Token + "/latest/" + initialCurrency

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	t := Target{}
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range currencies {
		results = append(results, amount*t.ConversionRates[s].(float64))
	}

	return results
}

func ConvertToTable(currencyNames []string, currencyValues []float64) string {
	cellLength := GetLength(currencyValues)
	tableLine := "+-----+-" + strings.Repeat("-", cellLength) + "-+\n"
	table := tableLine

	for i, value := range currencyValues {
		formattedValue := strconv.FormatFloat(value, 'f', -1, 64)
		padding := strings.Repeat(" ", (cellLength-len(formattedValue))/2)
		formattedValue = padding + formattedValue + padding
		if len(formattedValue) < cellLength {
			formattedValue += " "
		}
		table += fmt.Sprintf("| %s | %*s |\n", currencyNames[i], cellLength, formattedValue)
		table += tableLine
	}

	return table
}

func GetLength(numbers []float64) int {
	maxLength := 0
	for _, n := range numbers {
		length := len(strconv.FormatFloat(n, 'f', -1, 64))
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}
