package main

import (
	"encoding/json"
	"net/http"
)




func handleValidate(w http.ResponseWriter, req *http.Request) {

  type parameters struct {
    Body string `json:"body"`
  }
  type returnVal struct {
    Valid bool `json:"valid"`
  }

  decoder := json.NewDecoder(req.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    return
  }

  const chripLength = 140

  if len(params.Body) > chripLength {
    respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
    return
  }

  respondWithJson(w, http.StatusOK, returnVal{
    Valid: true,
  })

}
