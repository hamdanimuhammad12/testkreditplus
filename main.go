package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

func initDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/xyz_multifinance"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/transaction", createTransaction).Methods("POST")
	router.HandleFunc("/customer", createCustomer).Methods("POST")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8000", router))
}

var mutex sync.Mutex

type Transaction struct {
	ContractNumber    string `json:"ContractNumber"`
	OTR               int    `json:"OTR"`
	AdminFee          int    `json:"AdminFee"`
	InstallmentAmount int    `json:"InstallmentAmount"`
	InterestAmount    int    `json:"InterestAmount"`
	AssetName         string `json:"AssetName"`
	CustomerID        int    `json:"CustomerID"`
}

type Customer struct {
	NIK          string `json:"NIK"`
	FullName     string `json:"FullName"`
	LegalName    string `json:"LegalName"`
	PlaceOfBirth string `json:"PlaceOfBirth"`
	DateOfBirth  string `json:"DateOfBirth"`
	Salary       int    `json:"Salary"`
	KTPPhoto     string `json:"KTPPhoto"`
	SelfiePhoto  string `json:"SelfiePhoto"`
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var input Transaction
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO transactions(contract_number, otr, admin_fee, installment_amount, interest_amount, asset_name, customer_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(input.ContractNumber, input.OTR, input.AdminFee, input.InstallmentAmount, input.InterestAmount, input.AssetName, input.CustomerID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Transaction created successfully")
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var input Customer
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO customers(nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp_photo, selfie_photo) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(input.NIK, input.FullName, input.LegalName, input.PlaceOfBirth, input.DateOfBirth, input.Salary, input.KTPPhoto, input.SelfiePhoto)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Customer created successfully")
}
