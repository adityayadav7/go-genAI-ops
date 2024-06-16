package handlers

import (
	// "database/sql"
	"fmt"
	"encoding/json"
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
	"crudapp/models"
	"crudapp/db"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)

	// Insert newItem into the database
	insertSQL := "INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id"
	err := db.DB.QueryRow(insertSQL, newItem.Name, newItem.Description).Scan(&newItem.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
// func createUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u Item
// 		json.NewDecoder(r.Body).Decode(&u)

// 		err := db.QueryRow("INSERT INTO  (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		json.NewEncoder(w).Encode(u)
// 	}
// }
func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// fmt.Println(rows.Next())
	var items []models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(item)
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item models.Item

	err := db.DB.QueryRow("SELECT * FROM items WHERE id=$1",id).Scan(&item.ID, &item.Name, &item.Description)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var updatedItem models.Item
	_ = json.NewDecoder(r.Body).Decode(&updatedItem)

	vars := mux.Vars(r)
	id := vars["id"]

	updateSQL := "UPDATE items SET name=$1,description=$2 WHERE id=$3"

	_, err := db.DB.Exec(updateSQL,updatedItem.Name, updatedItem.Description,id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id := vars["id"]

	_,err := db.DB.Exec("DELETE from items WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
