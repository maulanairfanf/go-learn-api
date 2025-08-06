package routes

import (
	"myapi/handlers"
	"myapi/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// InitializeRoutes sets up the router with all the routes and middleware
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()


	// Root endpoint for health check or welcome message
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Go Learn API!"))
	}).Methods("GET")

	// Define routes
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	router.Handle("/products", withJWT(handlers.GetProducts)).Methods("GET")
	router.Handle("/products/{id}", withJWT(handlers.GetProduct)).Methods("GET")
	router.Handle("/products", withJWT(handlers.CreateProduct)).Methods("POST")
	router.Handle("/product/{id}", withJWT(handlers.DeleteProduct)).Methods("DELETE")
	router.Handle("/product/{id}", withJWT(handlers.UpdateProduct)).Methods("PUT")

	router.Handle("/category",  withJWT(handlers.CreateCategory)).Methods("POST")
	router.Handle("/category", withJWT(handlers.GetCategories)).Methods("GET")
	router.Handle("/category/{id}", withJWT(handlers.GetCategories)).Methods("GET")

	return router
}

// withJWT is a helper function to wrap handlers with the JWT middleware
func withJWT(handler http.HandlerFunc) http.Handler {
	return middleware.JWTMiddleware(handler)
}
