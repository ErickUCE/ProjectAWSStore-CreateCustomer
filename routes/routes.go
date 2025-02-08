package routes

import (
	"ProjectAWSStore-CreateCustomer/controllers" // ✅ Asegúrate de que el módulo es correcto

	"github.com/gorilla/mux"
)

// SetupRoutes configura las rutas del servidor
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// 📌 Definir rutas principales
	router.HandleFunc("/customers", controllers.CreateCustomer).Methods("POST")
	//router.HandleFunc("/customers", controllers.GetAllCustomers).Methods("GET")
	//	router.HandleFunc("/customers/{id}", controllers.GetCustomerByID).Methods("GET")
	//router.HandleFunc("/customers/{id}", controllers.UpdateCustomer).Methods("PUT")
	//router.HandleFunc("/customers/{id}", controllers.DeleteCustomer).Methods("DELETE")

	// 📌 Ruta para sincronizar clientes desde otro microservicio
	router.HandleFunc("/sync-create", controllers.SyncCreateCustomer).Methods("POST")

	return router
}
