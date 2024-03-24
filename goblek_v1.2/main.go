package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

// Buat koneksi ke Redis
var rdb = redis.NewClient(&redis.Options{
	Addr:     "redis:6379", // Alamat Redis dalam Docker Compose
	Password: "",           // Password Redis Anda, jika ada
	DB:       0,            // Nomor database Redis yang akan digunakan
})

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Ambil data dari form
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusInternalServerError)
				return
			}

			// Simpan data ke Redis
			data := r.FormValue("data")
			err = rdb.Set(context.Background(), "input_data", data, 0).Err()
			if err != nil {
				http.Error(w, "Failed to save data to Redis", http.StatusInternalServerError)
				return
			}
		}

		// Ambil data dari Redis
		data, err := rdb.Get(context.Background(), "input_data").Result()
		if err != nil && err != redis.Nil {
			http.Error(w, "Failed to get data from Redis", http.StatusInternalServerError)
			return
		}

		// Tampilkan form dan data
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Input Data</title>
			</head>
			<body>
				<form method="post">
					<label for="data">Input Data:</label><br>
					<input type="text" id="data" name="data"><br>
					<input type="submit" value="Submit">
				</form>
				<br>
				<h2>Data:</h2>
				<p>%s</p>
			</body>
			</html>
		`, data)
	})

	// Menjalankan server HTTP di port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
