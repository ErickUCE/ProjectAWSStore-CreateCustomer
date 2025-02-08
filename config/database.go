package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database   // 🔥 Variable global para la base de datos
var client *mongo.Client // 🔥 Cliente de MongoDB

// ✅ Conectar a MongoDB
func ConnectDB() {
	fmt.Println("📌 Ejecutando ConnectDB...")

	// 🔥 Intentar cargar .env primero
	err := godotenv.Load()

	if err != nil {
		fmt.Println("⚠️ No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	// 📌 Leer valores de .env o establecer valores predeterminados
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if mongoURI == "" || mongoDBName == "" {
		log.Fatal("❌ ERROR: Las variables de entorno MONGO_URI o MONGO_DB_NAME están vacías.")
	}

	fmt.Println("🔗 URI cargada:", mongoURI)
	fmt.Println("🛢  Base de datos:", mongoDBName)

	// ✅ Configurar cliente de MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)

	var connectErr error
	client, connectErr = mongo.Connect(context.TODO(), clientOptions)
	if connectErr != nil {
		log.Fatal("❌ Error conectando a MongoDB:", connectErr)
	}

	// ✅ Verificar conexión
	pingErr := client.Ping(context.TODO(), nil)
	if pingErr != nil {
		log.Fatal("❌ Error haciendo ping a MongoDB:", pingErr)
	}

	// ✅ Asignar la base de datos
	DB = client.Database(mongoDBName)
	fmt.Println("✅ Conexión exitosa a MongoDB en", mongoDBName)
}

// GetDB devuelve la instancia de la base de datos ya conectada
func GetDB() *mongo.Database {
	if DB == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. Asegúrate de llamar a ConnectDB() primero.")
	}
	return DB
}
