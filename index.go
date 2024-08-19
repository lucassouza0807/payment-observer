package main

import (
	cora "pagamento/src/cora"
	mercadopago "pagamento/src/mercado_pago"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	interval := time.Second * 10

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {

		for range ticker.C {
			go cora.Handle()
			go mercadopago.Handle()

		}

	}()

	select {}
}
