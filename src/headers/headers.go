package lucasTypes

type TokenResponse struct {
	Token string `json:"token"`
}

type PaymentDetail struct {
	ID                int     `json:"id"`
	IDPagamento       string  `json:"id_pagamento"`
	DetalhesPagamento *string `json:"detalhes_pagamento"`
	StatusPagamento   string  `json:"status_pagamento"`
	NomeCliente       string  `json:"nome_cliente"`
	EmailCliente      *string `json:"email_cliente"`
	Contato           *string `json:"contato"`
	Beneficiario      string  `json:"beneficiario"`
	MetodoPagamento   string  `json:"metodo_pagamento"`
	MeioPagamento     string  `json:"meio_pagamento"`
}

type PostGeneric struct {
	Message string `json:"message"`
}

type GenericApiResponse struct {
	Status string `json:"status"`
}

type ApiResponse struct {
	Data []PaymentDetail `json:"data"`
}
type SendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}
