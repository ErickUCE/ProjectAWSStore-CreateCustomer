# Usar la imagen oficial de Go como base
FROM golang:1.23.3-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar go mod y go sum para asegurar que las dependencias se instalen correctamente
COPY go.mod go.sum ./
RUN go mod tidy

# Copiar el resto del código fuente de la aplicación
COPY . .

# Establecer las variables de entorno necesarias para la conexión a MongoDB
ENV MONGO_URI=mongodb://44.207.106.151:27017/CreateCustomerDB
ENV MONGO_DB_NAME=CreateCustomerDB
ENV PORT=8081

# Exponer el puerto de la aplicación
EXPOSE 8081

# Compilar la aplicación
RUN go build -o main .

# Comando para ejecutar la aplicación
CMD ["./main"]
