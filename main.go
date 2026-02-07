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
		"GET	/api/v1/product" : "show all product",
		"POST	/api/v1/product"	: "add product",
		"GET	/api/v1/product/{id}" : "show 1 product",
		"PUT"	/api/v1/product/{id}" : "update product",
		"PATCH	/api/v1/product{id}" : "update field product",
		"DELETE	/api/v1/product/{id}" : "delete 1 product",
		"GET	/api/v1/category" : "show all category",
		"POST	/api/v1/category" : "add kategori",
		"GET	/api/v1/category/{id}" : "show 1 category",
		"PUT"	/api/v1/category/{id}" : "update category",
		"PATCH	/api/v1/category{id}" : "update field category",
		"DELETE	/api/v1/category/{id}" : "delete 1 category",
		"POST	/api/v1/checkout" : "create transaction",
		"GET	/api/v1/report" : "show all transaction",
		"GET	/api/v1/report/today" : "show today's transaction",
		"GET	/api/v1/report?start_date={start_day}&end_date={end_day}" : "show transaction between days",
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
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Handler
	http.Handle("/api/v1/product", productHandler)
	http.Handle("/api/v1/product/", productHandler)
	http.Handle("/api/v1/category", categoryHandler)
	http.Handle("/api/v1/category/", categoryHandler)
	http.Handle("/api/v1/checkout", transactionHandler)
	http.Handle("/api/v1/report", transactionHandler)
	http.Handle("/api/v1/report/today", transactionHandler)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API running OK"))
	})

	fmt.Println("Server running on port: " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
