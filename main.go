package main

import (
	"fmt"
	"log"
	"net/http"

	"ProjectAWSStore-CreateCustomer/config"
	"ProjectAWSStore-CreateCustomer/routes"
)

func main() {
	fmt.Println("🚀 Iniciando CustomerService en Golang...")

	// ✅ Conectar a la base de datos MongoDB
	config.ConnectDB()

	// ✅ Configurar rutas
	router := routes.SetupRoutes()

	// ✅ Cambiar el puerto si el 8080 ya está en uso
	PORT := "8081" // Cambia el puerto aquí
	fmt.Printf("✅ Servidor corriendo en el puerto %s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
