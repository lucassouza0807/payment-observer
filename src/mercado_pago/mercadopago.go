package mercadopago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	lucasTypes "pagamento/src/headers"
	"pagamento/src/telegram"
)

func Handle() {
	fmt.Println("Chamando Mercado Pago\n\r")

	var url string = fmt.Sprintf("%s/obter-pagamentos/mercado-pago/pending", os.Getenv("SERVER_URL"))

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Erro ao chamar a API:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		fmt.Println("Erro na resposta da API. Status:", resp.Status)
		return
	}

	var apiResponse lucasTypes.ApiResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}

	// Processando cada item em "data"
	for _, paymentDetail := range apiResponse.Data {

		var url string = fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentDetail.IDPagamento)
		var token string = fmt.Sprintf("Bearer %s", os.Getenv("MP_ACCESS_TOKEN"))
		var GenericApiResponse lucasTypes.GenericApiResponse

		req, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			fmt.Println("Err", err.Error())

			return
		}

		req.Header.Set("Authorization", token)

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Erro ao chamar a API", err.Error())

			return
		}

		if err := json.NewDecoder(resp.Body).Decode(&GenericApiResponse); err != nil {
			fmt.Println("Erro ao chamar API", err.Error())

			return
		}

		defer resp.Body.Close()

		var is_paid bool = GenericApiResponse.Status == "APPROVED"

		if !is_paid {
			fmt.Println("Este pedido ainda não foi pago.\n\r")
			return
		}

		message := fmt.Sprintf(
			"Pagamento do pedido %s aprovado \nNome do cliente: %s",
			paymentDetail.IDPagamento,
			paymentDetail.NomeCliente,
		)

		telegram.SendMessage(message)

	}
}
