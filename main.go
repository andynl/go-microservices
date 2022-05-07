package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/andynl/go-microservices/data"
	"github.com/andynl/go-microservices/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	ph := handlers.NewProducts(l, v)

	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)

	// opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
