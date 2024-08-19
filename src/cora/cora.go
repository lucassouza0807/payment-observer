package cora

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	lucasTypes "pagamento/src/headers"
)

func getCoraToken(token *string) error {
	// Faz a requisição para obter o token
	resp, err := http.Get(fmt.Sprintf("%s/cora/token", os.Getenv("SERVER_URL")))
	if err != nil {
		return fmt.Errorf("erro ao chamar a API: %v", err)
	}
	defer resp.Body.Close()

	// Verifica se a requisição foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao chamar a API, Status: %s", resp.Status)
	}

	// Decodifica a resposta JSON para a struct TokenResponse
	var tokenResp lucasTypes.TokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("erro ao decodificar a resposta: %v", err)
	}

	// Atualiza o valor do token com o valor obtido
	*token = tokenResp.Token

	return nil
}

func Handle() {

	fmt.Println("Chamando Cora")

	url := fmt.Sprintf("%s/obter-pagamentos/cora/open", os.Getenv("SERVER_URL"))

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Erro ao chamar a API:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Erro ao chamar api Status: ", resp.Status)
		return
	}

	var apiResponse lucasTypes.ApiResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Println(err.Error())
		return
	}

	var token string

	// Obtém o token usando a função
	if err := getCoraToken(&token); err != nil {
		fmt.Println("Erro ao obter o token:", err)
		return
	}

	for _, paymentDetails := range apiResponse.Data {
		coraURL := fmt.Sprintf("%s/v2/invoices/%s", os.Getenv("CORA_API_ENDPOINT"), paymentDetails.IDPagamento)

		req, err := http.NewRequest(http.MethodGet, coraURL, nil)
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		certFile := "certs/certificate.pem"
		keyFile := "certs/private-key.key"

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Println("Erro ao carregar certificado e chave:", err)
			return
		}

		// Configurar TLS com o certificado
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// Configurar o cliente HTTP com as opções TLS
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}

		// Fazer a requisição usando o cliente configurado
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao chamar a API:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Erro ao chamar a API, Status:", resp.Status)
			return
		}

		var GenericApiResponse lucasTypes.GenericApiResponse

		fmt.Println("API chamada com sucesso!")

		if err := json.NewDecoder(resp.Body).Decode(&GenericApiResponse); err != nil {
			fmt.Println(err.Error())
			return
		}

		var is_paid bool = GenericApiResponse.Status == "PAID"

		is_paid = true

		if is_paid {
			var response lucasTypes.PostGeneric

			url := fmt.Sprintf("%s/mudar-status", os.Getenv("SERVER_URL"))

			data := map[string]string{
				"status":       "pago azideia",
				"id_pagamento": paymentDetails.IDPagamento,
			}

			body, err := json.Marshal(data)

			if err != nil {
				fmt.Println("Erro ao encodar o corpo")

				return
			}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

			if err != nil {
				fmt.Println("Erro ao estruturar request")

				return
			}

			req.Header.Set("Content-Type", "Application/json")

			client := http.Client{}

			resp, err := client.Do(req)

			if err != nil {
				fmt.Println("Erro ao fazer requisição", resp.Status)

				return
			}

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Erro na requisição status:", resp.StatusCode)

				return
			}

			if err != json.NewDecoder(resp.Body).Decode(&response) {
				fmt.Println("Erro ao decodar body")

				return
			}

			fmt.Println(response.Message)

			defer resp.Body.Close()

		}
	}

}
