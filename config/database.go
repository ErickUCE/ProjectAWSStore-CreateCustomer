package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database   // Variable global para la base de datos
var client *mongo.Client // Cliente de MongoDB

// ✅ Conectar a MongoDB
func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://54.158.252.115:27017/CustomerDB") // ⚠️ Verifica la IP de tu instancia EC2

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Error conectando a MongoDB:", err)
	}

	// Verificar conexión
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ Error haciendo ping a MongoDB:", err)
	}

	fmt.Println("✅ Conexión exitosa a MongoDB")
	DB = client.Database("CustomerDB") // ⚠️ Asegúrate de que el nombre de la base de datos es correcto
}

// ✅ Obtener colección asegurando que la conexión esté inicializada
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. Asegúrate de llamar a ConnectDB() primero.")
	}
	return DB.Collection(collectionName)
}
