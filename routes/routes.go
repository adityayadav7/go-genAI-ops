package routes

import (
	// "net/http"

	"crudapp/authentication"
	"crudapp/generativeai"
	"crudapp/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/items", handlers.GetItems).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.GetItem).Methods("GET")
	r.HandleFunc("/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", handlers.DeleteItem).Methods("DELETE")

	r.HandleFunc("/sendOtp", authentication.SendOtpHandler).Methods("POST")
	r.HandleFunc("/verifyOtp", authentication.VerifyOtpHandler).Methods("POST")
	r.HandleFunc("/generateAccessToken", authentication.GenerateAccessTokenHandler).Methods("POST")
	r.HandleFunc("/generate", generativeai.GenerateResponse).Methods("POST")
	return r
}
