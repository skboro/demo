package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/skboro/demo-auction/helper"
	"github.com/skboro/demo-auction/models"
)

func NewAuctionController(as *models.AuctionService) *AuctionController {
	return &AuctionController{
		as: as,
	}
}

type AuctionController struct {
	as *models.AuctionService
}

func (ac *AuctionController) Create(w http.ResponseWriter, r *http.Request) {
	if !helper.IsAdmin(r) {
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var form models.Auction
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.Create(&form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "auction created successfully", http.StatusOK)
}

func (ac *AuctionController) Update(w http.ResponseWriter, r *http.Request) {
	if !helper.IsAdmin(r) {
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var form models.Auction
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := ac.as.ByID(form.ID)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.Update(&form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "auction updated successfully", http.StatusOK)
}

func (ac *AuctionController) Delete(w http.ResponseWriter, r *http.Request) {
	if !helper.IsAdmin(r) {
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var form models.Auction
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.Delete(form.ID); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "auction deleted successfully", http.StatusOK)
}

func (ac *AuctionController) GetAll(w http.ResponseWriter, r *http.Request) {
	if !helper.IsAdmin(r) {
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	auction, err := ac.as.GetAll()
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(auction)
}

func (ac *AuctionController) GetLive(w http.ResponseWriter, r *http.Request) {
	auction, err := ac.as.GetLive()
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(auction)
}

func (ac *AuctionController) CreateBid(w http.ResponseWriter, r *http.Request) {
	var form models.Bid
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.CreateBid(&form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "bid created successfully", http.StatusOK)
}

func (ac *AuctionController) UpdateBid(w http.ResponseWriter, r *http.Request) {
	var form models.Bid
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := ac.as.GetBid(form.ID)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.UpdateBid(&form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "bid updated successfully", http.StatusOK)
}

func (ac *AuctionController) DeleteBid(w http.ResponseWriter, r *http.Request) {
	if !helper.IsAdmin(r) {
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var form models.Bid
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ac.as.DeleteBid(form.ID); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "bid deleted successfully", http.StatusBadRequest)
}

func (ac *AuctionController) GetBids(w http.ResponseWriter, r *http.Request) {
	var form models.Bid
	var err error
	if err = helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bids []models.Bid
	if form.AuctionID > 0 && helper.IsAdmin(r) {
		bids, err = ac.as.GetBidsByAuctionId(form.AuctionID)
	} else if form.UserID > 0 {
		bids, err = ac.as.GetBidsByUserId(form.UserID)
	} else if helper.IsAdmin(r) {
		bids, err = ac.as.GetAllBids()
	} else {
		helper.Response(w, "Please provide auction id or user id (non-admin supported) as argument", http.StatusBadRequest)
		return
	}
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(bids)
}

func (ac *AuctionController) GetBid(w http.ResponseWriter, r *http.Request) {
	var form models.Bid
	var err error
	if err = helper.ParseBody(r, &form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	if form.ID == 0 {
		helper.Response(w, "Please provide a valid bid id.", http.StatusBadRequest)
	}
	bid, err := ac.as.GetBid(form.ID)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(bid)
}
