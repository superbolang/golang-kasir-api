package main

import (
	"fmt"
	"gokasir-api/database"
	"gokasir-api/handler"
	"gokasir-api/repository"
	"gokasir-api/service"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var message = `{
	"endpoint" : {
		"GET	/api/v1/product" : "tampilkan semua product",
		"POST	/api/v1/product"	: "tambah product",
		"GET	/api/v1/product/{id}" : "tampilkan 1 product",
		"PUT"	/api/v1/product/{id}" : "update seluruh field",
		"PATCH	/api/v1/product{id}" : "update sebagian field",
		"DELETE	/api/v1/product/{id}" : "menghapus 1 product",
		"GET	/api/v1/kategori" : "tampilkan semua kategori",
		"POST	/api/v1/kategori" : "tambah kategori",
		"GET	/api/v1/kategori/{id}" : "tampilkan 1 kategori",
		"PUT"	/api/v1/kategori/{id}" : "update seluruh field",
		"PATCH	/api/v1/kategori{id}" : "update sebagian field",
		"DELETE	/api/v1/kategori/{id}" : "menghapus 1 kategori",
	},
	"environtment" : "production",
	"message" : "simple API",
	"version" : "1.0.0"
}`

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Init DB
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(message))
	})

	// Layer
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductServiceImpl(productRepository)
	productHandler := handler.NewProductHandler(productService)

	// Handler
	http.Handle("/api/v1/product", productHandler)
	http.Handle("/api/v1/product/", productHandler)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API running OK"))
	})

	fmt.Println("Server running on port: " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
