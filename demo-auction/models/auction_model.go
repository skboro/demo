package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Auction struct {
	gorm.Model
	StartTime    time.Time `json:"start_time" gorm:"not null;index:auction_start_time"`
	EndTime      time.Time `json:"end_time" gorm:"not null;index:auction_end_time"`
	StartPrice   int64     `json:"start_price" gorm:"not null"`
	ItemName     string    `json:"item_name" gorm:"not null"`
	WinnerUserID uint      `json:"winner_user_id"`
}

type Bid struct {
	gorm.Model
	Price     int64 `json:"price" gorm:"not null"`
	UserID    uint  `json:"user_id" gorm:"not null;index:bid_user_id"`
	AuctionID uint  `json:"auction_id" gorm:"not null;index:bid_auction_id"`
}

type AuctionService struct {
	db *gorm.DB
}

func NewAuctionService(db *gorm.DB) *AuctionService {
	as := AuctionService{
		db: db,
	}
	if !as.db.HasTable(&Auction{}) {
		as.db.CreateTable(&Auction{})
	}
	if !as.db.HasTable(&Bid{}) {
		as.db.CreateTable(&Bid{})
	}
	return &as
}

func (as *AuctionService) ByID(id uint) (*Auction, error) {
	var auction Auction
	db := as.db.Where("id = ?", id)
	err := db.First(&auction).Error
	return &auction, err
}

func (as *AuctionService) Create(auction *Auction) error {
	return as.db.Create(auction).Error
}

func (as *AuctionService) Update(auction *Auction) error {
	return as.db.Save(auction).Error
}

func (as *AuctionService) Delete(id uint) error {
	auction := Auction{Model: gorm.Model{ID: id}}
	return as.db.Delete(&auction).Error
}

func (as *AuctionService) GetLive() ([]Auction, error) {
	var auctions []Auction
	t := time.Now()
	err := as.db.Where("start_time < ? AND end_time > ?", t, t).Find(&auctions).Error
	return auctions, err
}

func (as *AuctionService) GetAll() ([]Auction, error) {
	var auctions []Auction
	err := as.db.Find(&auctions).Error
	return auctions, err
}

func (as *AuctionService) GetFinishedWithoutWinner() ([]Auction, error) {
	var auctions []Auction
	err := as.db.Where("end_time < ? AND winner_user_id = 0", time.Now()).Find(&auctions).Error
	return auctions, err
}

func (as *AuctionService) GetAllBids() ([]Bid, error) {
	var bids []Bid
	err := as.db.Find(&bids).Error
	return bids, err
}

func (as *AuctionService) GetBid(id uint) (*Bid, error) {
	var bid Bid
	err := as.db.Where("id = ?", id).First(&bid).Error
	return &bid, err
}

func (as *AuctionService) GetBidsByUserId(user_id uint) ([]Bid, error) {
	var bids []Bid
	db := as.db.Where("user_id = ?", user_id)
	err := db.Find(&bids).Error
	return bids, err
}

func (as *AuctionService) GetBidsByAuctionId(auction_id uint) ([]Bid, error) {
	var bids []Bid
	db := as.db.Where("auction_id = ?", auction_id).Order("price desc")
	err := db.Find(&bids).Error
	return bids, err
}

func (as *AuctionService) GetWinningBid(auction_id uint) (*Bid, error) {
	var bid Bid
	db := as.db.Where("auction_id = ?", auction_id).Order("price desc")
	err := db.First(&bid).Error
	return &bid, err
}

func (as *AuctionService) CreateBid(bid *Bid) error {
	auction, err := as.ByID(bid.AuctionID)
	if err != nil {
		return err
	}
	if auction.EndTime.Unix() < time.Now().Unix() {
		return errors.New("auction is over")
	}
	if auction.StartPrice > bid.Price {
		return errors.New("minimum bid is not met")
	}
	return as.db.Create(bid).Error
}

func (as *AuctionService) UpdateBid(bid *Bid) error {
	auction, err := as.ByID(bid.AuctionID)
	if err != nil {
		return err
	}
	if auction.EndTime.Unix() < time.Now().Unix() {
		return errors.New("auction is over")
	}
	if auction.StartPrice > bid.Price {
		return errors.New("minimum bid is not met")
	}
	return as.db.Save(bid).Error
}

func (as *AuctionService) DeleteBid(id uint) error {
	bid := Bid{Model: gorm.Model{ID: id}}
	return as.db.Delete(&bid).Error
}
