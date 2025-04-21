package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "user=postgres password=пароль123 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", handleProducts)

	fmt.Println("Server running on :8080")

	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	} else {
		fmt.Println("✅ Успешное подключение к базе данных")
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		addProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, description, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", 400)
		log.Println("Ошибка декодирования JSON:", err)
		return
	}
	log.Printf("Добавляем товар: %+v\n", p)
	res, err := db.Exec(
		"INSERT INTO products (name, description, price) VALUES ($1, $2, $3)",
		p.Name, p.Description, p.Price,
	)
	if err != nil {
		log.Println("Ошибка при вставке в базу:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	log.Printf("✅ Вставлено строк: %d\n", rowsAffected)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Product added"}`))
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if (*w).Header().Get("Content-Type") == "" {
		(*w).Header().Set("Content-Type", "application/json")
	}
}
