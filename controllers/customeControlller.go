package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ProjectAWSStore-CreateCustomer/config"
	"ProjectAWSStore-CreateCustomer/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCustomerCollection() *mongo.Collection {
	if config.DB == nil {
		fmt.Println("❌ Error: La base de datos aún no está inicializada.")
		return nil
	}
	return config.DB.Collection("customers")
}

// 📌 Crear un nuevo cliente y sincronizarlo con `ReadCustomer`
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	customerCollection := getCustomerCollection()
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Convertir ObjectID a string antes de asignarlo
	customer.ID = primitive.NewObjectID().Hex()

	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente creado:", customer.Email)

	// 🔄 **Sincronizar con ReadCustomer**
	go syncWithReadCustomer(customer)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// 📌 Función para enviar los datos a `ReadCustomer`
func syncWithReadCustomer(customer models.Customer) {
	readCustomerURL := "http://localhost:8082/sync-create" // ⚠️ URL del microservicio ReadCustomer

	// Serializar el cliente en JSON
	jsonData, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("❌ Error serializando cliente:", err)
		return
	}

	// Enviar petición HTTP
	resp, err := http.Post(readCustomerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Error notificando a ReadCustomer:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Cliente sincronizado con ReadCustomer:", customer.Email)
	} else {
		fmt.Println("⚠️ No se pudo sincronizar cliente con ReadCustomer. Código:", resp.StatusCode)
	}
}

// 📌 Sincronizar creación de clientes desde otro microservicio
func SyncCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	customerCollection := getCustomerCollection()
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Verificar si el cliente ya existe
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"email": customer.Email}).Decode(&existingCustomer)
	if err == nil {
		fmt.Println("⚠️ Cliente ya existe, no se duplica:", customer.Email)
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Insertar nuevo cliente
	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "Failed to sync customer", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente:", customer.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
