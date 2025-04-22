package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Auth struct {
	Login    string `json:"login"` // Не смей забывать о больших буквах для полей.
	Password string `json:"password"`
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	var authTemp Auth

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&authTemp)
		if err != nil {
			log.Println(err)
		}
		// Я аж бл Console.WriteLine(); вспомнил пока всё это писал...
		fmt.Println(authTemp.Login, authTemp.Password) // Робит через b, err := io.ReadAll(r.Body) и Println(string(b))
	}
}

func main() {
	http.HandleFunc("/authorize", Authorize)
	http.ListenAndServe("localhost:8080", nil)
}

/*
	- Нужно чтоб через постмэн отправился логин и пароль на /login
	- Сервер проверил данные
	- Ответил Входом или не входом сообщением ок не ок.
	- После этого загорается флаг isAuthinticated и можно использовать что-то далее.
	= ВСЁ ЭТА КОММУНИКАЦИЯ ПРОХОДИТ ЧЕРЕЗ JSON.
*/
