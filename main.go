package main

import (
	"fmt"
	"net/http"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shadowshot-x/micro-product-go/authservice"
	"github.com/shadowshot-x/micro-product-go/authservice/middleware"
	"github.com/shadowshot-x/micro-product-go/clientclaims"
)

func main() {
	mainRouter := mux.NewRouter()
	// We will create a Subrouter for Authentication service
	// route for sign up and signin. The Function will come from auth-service package
	// checks if given header params already exists. If not,it adds the user
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authservice.SignupHandler).Methods("POST")

	// The Signin will send the JWT Token back as we are making microservices.
	// JWT token will make sure that other services are protected.
	// So, ultimately, we would need a middleware
	authRouter.HandleFunc("/signin", authservice.SigninHandler).Methods("GET")

	// File Upload SubRouter
	claimsRouter := mainRouter.PathPrefix("/claims").Subrouter()
	claimsRouter.HandleFunc("/upload", clientclaims.UploadFile)
	claimsRouter.Use(middleware.TokenValidationMiddleware)

	// CORS Header
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
	// Add the Middleware to different subrouter
	// HTTP Server
	// Add Time outs
	server := &http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: ch(mainRouter),
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error Booting the Server")
	}
}
