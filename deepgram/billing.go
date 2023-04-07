package deepgram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Balance struct {
	Amount          float64 `json:"amount"`
	BalanceId       string  `json:"balance_id"`
	Units           string  `json:"units"`
	PurchaseOrderId string  `json:"purchase_order_id"`
}

type BalanceList struct {
	Balances []Balance `json:"balances"`
}

func (dg *Client) ListBalances(projectId string) (BalanceList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/balances", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result BalanceList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting balances from project %s: %s\n", projectId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetBalance(projectId string, balanceId string) (Balance, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/balances/%s", dg.Path, projectId, balanceId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Balance
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting balance %s: %s\n", balanceId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}
