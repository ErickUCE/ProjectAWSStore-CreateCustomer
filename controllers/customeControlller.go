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
)

// ✅ Conectar a la colección de clientes
var customerCollection = config.GetCollection("customers")

// 📌 **Crear un nuevo cliente desde la API principal**
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// ✅ Generar un nuevo ID único de MongoDB
	customer.ID = primitive.NewObjectID().Hex()

	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	// ✅ Notificar a los demás microservicios
	notifyMicroservices(customer)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// ✅ **Sincronizar la creación desde otros microservicios**
func SyncCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// ✅ Verificar si ya existe antes de insertar
	existingCustomer := customerCollection.FindOne(context.TODO(), bson.M{"_id": customer.ID})
	if existingCustomer.Err() == nil {
		fmt.Println("⚠️ Cliente ya existe, no se inserta nuevamente")
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "Failed to sync customer", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado exitosamente en CreateCustomerDB")
	w.WriteHeader(http.StatusOK)
}

// 📌 **Notificar a otros microservicios sobre la creación de un cliente**
func notifyMicroservices(customer models.Customer) {
	instances := []string{
		"http://54.158.252.116:8082/sync-create", // ReadCustomerDB
		"http://54.158.252.117:8083/sync-create", // UpdateCustomerDB
		"http://54.158.252.118:8084/sync-create", // DeleteCustomerDB
	}

	for _, instance := range instances {
		go func(instance string) {
			jsonData, _ := json.Marshal(customer)
			resp, err := http.Post(instance, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("❌ Error notificando a:", instance, err)
				return
			}
			defer resp.Body.Close()
			fmt.Println("✅ Cliente sincronizado en", instance)
		}(instance)
	}
}
