package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	initDB()
	defer db.Close()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"ContractNumber":    "12345",
		"OTR":               1000000,
		"AdminFee":          10000,
		"InstallmentAmount": 200000,
		"InterestAmount":    5000,
		"AssetName":         "Motor",
		"CustomerID":        1,
	})

	req, err := http.NewRequest("POST", "/transaction", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTransaction)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := "Transaction created successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateCustomer(t *testing.T) {
	initDB()
	defer db.Close()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"NIK":          "1234567890",
		"FullName":     "Budi Santoso",
		"LegalName":    "Budi Santoso",
		"PlaceOfBirth": "Jakarta",
		"DateOfBirth":  "1990-01-01",
		"Salary":       5000000,
		"KTPPhoto":     "base64-encoded-ktp-photo",
		"SelfiePhoto":  "base64-encoded-selfie-photo",
	})

	req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := "Customer created successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
