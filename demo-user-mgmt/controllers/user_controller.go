package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/skboro/demo-user-mgmt/helper"
	"github.com/skboro/demo-user-mgmt/models"
)

func NewUserController(us *models.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

type UserController struct {
	us *models.UserService
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(id uint) (string, error) {
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(os.Getenv("jwt_key")))
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				helper.Response(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			helper.Response(w, err.Error(), http.StatusBadRequest)
			return
		}
		tknStr := c.Value

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwt_key")), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				helper.Response(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			helper.Response(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			helper.Response(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.AddCookie(&http.Cookie{
			Name:  "user_id",
			Value: strconv.FormatUint(uint64(claims.UserID), 10),
		})
		next(w, r)
	})
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var form models.User
	var user *models.User
	var err error
	if err = helper.ParseBody(r, &form); err != nil {
		helper.Response(w, "some error occurred", http.StatusBadRequest)
		return
	}

	// skip login for admin if using secret key
	if form.Email == "admin@sellerapp.com" && form.Password == "admin_secret_key" {
		user = &models.User{}
		user.ID = 0
	} else {
		user, err = u.us.Authenticate(form.Email, form.Password)
		if err != nil {
			helper.Response(w, "invalid username/password", http.StatusBadRequest)
			return
		}
	}

	tokenString, err := CreateToken(user.ID)
	if err != nil {
		helper.Response(w, "some error occurred", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	helper.Response(w, tokenString, http.StatusOK)
}

func (u *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	var form models.User
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, "some error occurred", http.StatusBadRequest)
		return
	}
	if err := u.us.Create(&form); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "account created successfully", http.StatusBadRequest)
}

func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var form models.User
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, "some error occurred", http.StatusBadRequest)
		return
	}
	user, err := u.us.ByID(form.ID)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Name = form.Name
	if err := u.us.Update(user); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "account updated successfully", http.StatusBadRequest)
}

func (u *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	var form models.User
	if err := helper.ParseBody(r, &form); err != nil {
		helper.Response(w, "some error occurred", http.StatusBadRequest)
		return
	}
	if err := u.us.Delete(form.ID); err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	helper.Response(w, "account deleted successfully", http.StatusBadRequest)
}

func (u *UserController) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user_id")
	if err != nil {
		if err == http.ErrNoCookie {
			helper.Response(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.Value != "0" { // admin has account id 0
		helper.Response(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := u.us.GetAllUsers()
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (u *UserController) GetAccount(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("user_id")
	if err != nil {
		if err == http.ErrNoCookie {
			helper.Response(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, _ := strconv.ParseUint(c.Value, 10, 32)
	user, err := u.us.ByID(uint(uid))
	if err != nil {
		helper.Response(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}
