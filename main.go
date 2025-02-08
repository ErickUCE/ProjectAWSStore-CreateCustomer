package main

import (
	"fmt"
	"log"
	"net/http"

	"ProjectAWSStore-CreateCustomer/config"
	"ProjectAWSStore-CreateCustomer/routes"
)

func main() {
	fmt.Println("🚀 Iniciando CreateCustomer Service en Golang...")

	// ✅ Conectar a la base de datos
	config.ConnectDB()

	// ✅ Configurar rutas después de la conexión a la base de datos
	router := routes.SetupRoutes()

	// ✅ Iniciar el servidor en el puerto 8081
	fmt.Println("✅ Servidor corriendo en el puerto 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
