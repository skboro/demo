package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/skboro/demo-auction/controllers"
	"github.com/skboro/demo-auction/models"
)

func initDB(config *Config) *gorm.DB {
	db, err := gorm.Open("mysql", config.Database.ConnectionInfo())
	if err != nil {
		panic(err)
	}
	return db
}

func worker(as *models.AuctionService) {
	for {
		auctions, err := as.GetFinishedWithoutWinner()
		if err == nil {
			for _, auction := range auctions {
				bid, err := as.GetWinningBid(auction.ID)
				if err == nil {
					auction.WinnerUserID = bid.UserID
					as.Update(&auction)
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func main() {
	config := LoadConfig()
	db := initDB(config)
	defer db.Close()
	as := models.NewAuctionService(db)
	ac := controllers.NewAuctionController(as)
	router := mux.NewRouter()

	router.HandleFunc("/auction/create", ac.Create).Methods("POST")
	router.HandleFunc("/auction/update", ac.Update).Methods("POST")
	router.HandleFunc("/auction/delete", ac.Delete).Methods("POST")
	router.HandleFunc("/auction/getAll", ac.GetAll).Methods("GET")
	router.HandleFunc("/auction/getLive", ac.GetLive).Methods("GET")

	router.HandleFunc("/auction/bid/create", ac.CreateBid).Methods("POST")
	router.HandleFunc("/auction/bid/update", ac.UpdateBid).Methods("POST")
	router.HandleFunc("/auction/bid/delete", ac.DeleteBid).Methods("POST")
	router.HandleFunc("/auction/bid/get", ac.GetBid).Methods("POST")
	router.HandleFunc("/auction/bid/getBids", ac.GetBids).Methods("POST")

	go worker(as)
	http.ListenAndServe(":"+strconv.FormatInt(config.Port, 10), router)
}
