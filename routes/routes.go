package routes

import (
	"ProjectAWSStore-CreateCustomer/controllers"

	"github.com/gorilla/mux"
)

// SetupRoutes configura las rutas del servidor
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Definir rutas
	router.HandleFunc("/customers", controllers.CreateCustomer).Methods("POST")
	router.HandleFunc("/sync-create", controllers.SyncCreateCustomer).Methods("POST")

	return router
}
