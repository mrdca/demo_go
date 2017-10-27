package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string){
	conectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", conectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string){
	log.Fatal(http.ListenAndServe(":1234", a.Router))
}


func (a *App) getSmartphone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smartphone ID")
		return
	}

	s := smartphone{ID: id}
	if err := s.getSmartphone(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Smartphone not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func (a *App) getSmartphones(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	smartphones, err := getSmartphones(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, smartphones)
}

func (a *App) createSmartphone(w http.ResponseWriter, r *http.Request) {
	var s smartphone
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.createSmartphone(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) updateSmartphone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smartphone ID")
		return
	}

	var s smartphone
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	s.ID = id

	if err := s.updateSmartphone(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *App) deleteSmartphone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smartphone ID")
		return
	}

	s := smartphone{ID: id}
	if err := s.deleteSmartphone(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/getSmartphones", a.getSmartphones).Methods("GET")
	a.Router.HandleFunc("/createSmartphone", a.createSmartphone).Methods("POST")
	a.Router.HandleFunc("/getSmartphone/{id:[0-9]+}", a.getSmartphone).Methods("GET")
	a.Router.HandleFunc("/updateSmartphone/{id:[0-9]+}", a.updateSmartphone).Methods("PUT")
	a.Router.HandleFunc("/deleteSmartphone/{id:[0-9]+}", a.deleteSmartphone).Methods("DELETE")
}