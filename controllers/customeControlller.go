package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ProjectAWSStore-CreateCustomer/config"
	"ProjectAWSStore-CreateCustomer/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Obtiene la colección de clientes desde la base de datos
func getCustomerCollection() *mongo.Collection {
	if config.DB == nil {
		fmt.Println("❌ Error: La base de datos aún no está inicializada.")
		return nil
	}
	return config.DB.Collection("customers")
}

// 📌 **Crear un nuevo cliente y sincronizarlo con `ReadCustomer` y `UpdateCustomer`**
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

	// ✅ Generar un nuevo `ObjectID`
	customer.ID = primitive.NewObjectID()

	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente creado:", customer.Email)

	// 🔄 **Sincronizar con ReadCustomer y UpdateCustomer**
	go syncWithMicroservices(customer)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// 📌 **Función para enviar los datos a `ReadCustomer` y `UpdateCustomer`**
func syncWithMicroservices(customer models.Customer) {
	services := []string{
		"http://54.87.55.114:8082/sync-create", // ReadCustomer

		"http://l44.217.27.149:8083/sync-create", // UpdateCustomer
		"http://52.205.123.191:8084/sync-create", // UpdateCustomer
	}

	customerJSON, _ := json.Marshal(customer)

	for _, service := range services {
		resp, err := http.Post(service, "application/json", bytes.NewBuffer(customerJSON))
		if err != nil {
			fmt.Println("❌ Error notificando a", service, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
			fmt.Println("✅ Cliente sincronizado con:", service)
		} else {
			fmt.Println("⚠️ No se pudo sincronizar cliente con", service, "Código:", resp.StatusCode)
		}
	}
}

// 📌 **Sincronizar creación de clientes desde otro microservicio**
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

// 📌 **Sincronizar actualización de clientes desde `UpdateCustomer`**
func SyncUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updatedCustomer models.Customer
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización para:", updatedCustomer.Email)

	customerCollection := getCustomerCollection()
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Validar si `ID` está vacío
	if updatedCustomer.ID == primitive.NilObjectID {
		fmt.Println("⚠️ Error: `ID` vacío en la sincronización de actualización.")
		http.Error(w, "⚠️ Error: `ID` vacío en la sincronización", http.StatusBadRequest)
		return
	}

	// 📌 Crear el filtro para buscar por `_id`
	filter := bson.M{"_id": updatedCustomer.ID}
	update := bson.M{"$set": updatedCustomer}

	// 📌 Intentar actualizar el cliente en la base de datos
	result, err := customerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("❌ Error al actualizar cliente en MongoDB:", err)
		http.Error(w, "❌ Error al sincronizar actualización", http.StatusInternalServerError)
		return
	}

	// 📌 Verificar si realmente se encontró y actualizó el cliente
	if result.MatchedCount == 0 {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización.")
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización.", http.StatusNotFound)
		return
	}

	// ✅ Cliente actualizado correctamente
	fmt.Println("✅ Cliente sincronizado correctamente en CreateCustomer/ReadCustomer:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
}

func SyncDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// 🔄 Intentar convertir el ID de string a ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Error: ID inválido en la sincronización de eliminación:", id)
		http.Error(w, "❌ ID inválido en la sincronización", http.StatusBadRequest)
		return
	}

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 🔍 Verificar si el cliente existe antes de eliminarlo
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingCustomer)
	if err != nil {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización:", id)
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización", http.StatusNotFound)
		return
	}

	// 🗑️ Eliminar el cliente
	_, err = customerCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println("❌ Error al eliminar cliente en sincronización:", err)
		http.Error(w, "❌ Error al eliminar cliente en sincronización", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente eliminado en sincronización:", existingCustomer.Email)
	w.WriteHeader(http.StatusOK)
}
