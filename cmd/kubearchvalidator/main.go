package main

import (
	"kubearchvalidator/pkg/webhook"
)

func main() {
	ac := webhook.NewAdmissionController(":8080", webhook.Admit)
	if err := ac.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}
