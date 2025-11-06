package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	const charLimit = 140

	resBody := returnVals{}
	w.Header().Set("Content-Type", "application/json")

	if len(params.Body) > charLimit {
		w.WriteHeader(http.StatusBadRequest)
		resBody.Error = "Chirp is too long"
	} else {
		w.WriteHeader(http.StatusOK)
		resBody.Valid = true
	}

	dat, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dat)

}
