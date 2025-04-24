package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GentleIT/minLogic/minLogic"
	_ "github.com/lib/pq"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type token struct {
	Token int `json:"token"` // Не уверен что здесь нужно добавлять тег. Я ведь отправляю а не читаю
}

type Equipment struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Driver string    `json:"driver"`
	Day    time.Time `json:"day"`
	Gps    int       `json:"gps"`
	Parked bool      `json:"parked"`
}

var db *sql.DB // Эту практику...

func main() {
	var err error // И эту практику подсмотрел из одного примера.
	connStr := "user=postgres password=пароль123 dbname=postgres sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Не удалось подключиться к базе postgres 🔴")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Подключение с базой данных установлено 🟢")
	}

	http.HandleFunc("/getAll", GetEquipment)
	http.HandleFunc("/find", FindEquipment)
	http.HandleFunc("/equipment", AddEquipment)
	http.HandleFunc("/authorize", Authorization)
	http.ListenAndServe("localhost:8080", nil)
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, token")

	var authTemp Auth
	log.Println("Someone logging in")

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&authTemp)
		if err != nil {
			w.Write([]byte("400"))
			log.Println(err)
		}
		if enter := SignUpCheck(authTemp.Login, authTemp.Password); enter { // Почему-то у меня чувство, будто я наговнокодил
			byteToken, _ := json.MarshalIndent(token{Token: 123456789}, "", "") // Не думаю что indent стоит использовать, но всё же
			w.Write([]byte("You've successfully entered to the system\n"))
			w.Write(byteToken)
		}
	}
}

func SignUpCheck(log string, pass string) bool { // Это должно работать через базу данных
	if log == "Admin" && pass == "12345" {
		return true
	}
	return false
}

func AddEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, token")

	var newEquipment Equipment
	log.Println("Someone using AddEquipment")
	if r.Method == http.MethodPost && r.Header.Get("token") == "123456789" {
		w.Write([]byte("Successfully logged into AddEquipment"))

		err := json.NewDecoder(r.Body).Decode(&newEquipment)
		if err != nil {
			log.Println(err)
		}

		db.Exec("INSERT INTO equipment (name, driver, day, gps, parked) values ($1, $2, $3, $4, $5)", newEquipment.Name, newEquipment.Driver, minLogic.TimeFormat(time.Now()), minLogic.GetRandomGPS(), newEquipment.Parked)
	} else {
		w.Write([]byte("400"))
	}
}

func FindEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, token")

	var resp Equipment
	log.Println("Someone using FindEquipment")

	if r.Method == http.MethodGet && r.Header.Get("token") == "123456789" {
		w.Write([]byte("Successfully logged into FindEquipment\n"))

		err := json.NewDecoder(r.Body).Decode(&resp) // Расшифровка
		if err != nil {
			log.Println(err, "FindEquipment's problem with decode")
		}
		fmt.Println(resp)

		dbres, err := db.Query("SELECT * FROM equipment WHERE name ILIKE $1", resp.Name)
		if err != nil {
			log.Println(err, "FindEquipment's problem with db")
		}
		defer dbres.Close()

		equipments := []Equipment{}
		if err = json.NewEncoder(w).Encode(ReadRows(dbres, equipments)); err != nil { // Шифровка
			log.Println("Problem with delievering information 🔴: ", err)
		}
	} else {
		w.Write([]byte("400"))
	}
}

func GetEquipment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, token") // Честно сказать - я не знаю как это решило мои проблемы. CORS? What?

	log.Println("Someone using GetEquipment")

	if r.Method == http.MethodGet && r.Header.Get("token") == "123456789" {
		rows, _ := db.Query("SELECT * FROM equipment")

		equipment := []Equipment{}
		err := json.NewEncoder(w).Encode(ReadRows(rows, equipment))
		if err != nil {
			log.Println(err, "GetEquipment problem 🔴")
		}
	} else {
		json.NewEncoder(w).Encode("400")
	}
}

func ReadRows(rows *sql.Rows, eq []Equipment) []Equipment {
	for rows.Next() {
		p := Equipment{}
		err := rows.Scan(&p.ID, &p.Name, &p.Driver, &p.Day, &p.Gps, &p.Parked)
		if err != nil {
			log.Println(err, rows)
		}
		eq = append(eq, p)
	}
	return eq
}

/*
	Смайлики для отладки:
	🔴	🟡	🟢	🔵

	Изучение:
		- Что такое CORS.
		- Пробовать на Vue что-то поделать.
*/
