package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/skboro/demo-user-mgmt/controllers"
	"github.com/skboro/demo-user-mgmt/models"
)

var config *Config

func initDB(config *Config) *gorm.DB {
	db, err := gorm.Open("mysql", config.Database.ConnectionInfo())
	if err != nil {
		panic(err)
	}
	return db
}

func RequestForwarder(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	url := fmt.Sprintf("http://%s%s", config.AuctionURL, req.RequestURI)
	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	proxyReq.Header = req.Header
	httpClient := http.Client{}
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func main() {
	config = LoadConfig()
	db := initDB(config)
	defer db.Close()
	us := models.NewUserService(db)
	uc := controllers.NewUserController(us)
	router := mux.NewRouter()

	router.HandleFunc("/signin", uc.Login).Methods("POST")
	router.HandleFunc("/signup", uc.Signup).Methods("POST")
	router.HandleFunc("/update", controllers.Authenticate(uc.Update)).Methods("POST")
	router.HandleFunc("/delete", controllers.Authenticate(uc.Delete)).Methods("POST")
	router.HandleFunc("/getAllUsers", controllers.Authenticate(uc.GetAllAccounts)).Methods("GET")
	router.HandleFunc("/account", controllers.Authenticate(uc.GetAccount)).Methods("GET")

	router.PathPrefix("/auction").HandlerFunc(controllers.Authenticate(RequestForwarder))

	http.ListenAndServe(":"+strconv.FormatInt(config.Port, 10), router)
}
