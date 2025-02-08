package routes

import (
	//"net/http"
	"ProjectAWSStore-CreateCustomer/controllers" // ✅ Asegúrate de que el nombre del módulo es correcto

	"github.com/gorilla/mux"
)

// SetupRoutes configura las rutas del servidor
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Definir rutas
	router.HandleFunc("/customers", controllers.CreateCustomer).Methods("POST")
	router.HandleFunc("/customers", controllers.GetAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", controllers.GetCustomerByID).Methods("GET")
	router.HandleFunc("/customers/{id}", controllers.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", controllers.DeleteCustomer).Methods("DELETE")

	return router
}
