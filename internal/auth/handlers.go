package auth

import (
	"backend/contracts/auth"
	"backend/contracts/db"
	"backend/utils"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := utils.ReadJSON(w, r, &reqPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(reqPayload.Email)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid login credentials"), http.StatusBadRequest)
		return
	}

	valid, err := passwordMatches(user, reqPayload.Password)
	if err != nil || !valid {
		utils.ErrorJSON(w, errors.New("invalid login credentials"), http.StatusBadRequest)
		return
	}

	businessName := r.Header.Get("Business-Name")

	business, err := db.GetBusinessByBusinessName(businessName)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid business"), http.StatusBadRequest)
		return
	}

	isUserInBusiness, err := db.IsUserInBusiness(user.ID, business.ID)
	if err != nil || !isUserInBusiness {
		utils.ErrorJSON(w, errors.New("user not associated with the business"), http.StatusBadRequest)
		return
	}

	// create JWT user
	u := auth.AuthUser{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		BusinessID: business.ID,
	}

	// generate JWT tokens
	tokenPair, err := a.generateSignedTokenPair(&u)
	if err != nil {
		log.Fatalf("Error generating token pair: %v", err)
	}

	accessTokenExpiryInt := int(a.AccessTokenExpiry.Seconds())
	refreshTokenExpiryInt := int(a.RefreshTokenExpiry.Seconds())

	type Token struct {
		Token         string `json:"token"`
		ExpirySeconds int    `json:"expiry_seconds"`
	}

	type LoginResponse struct {
		AccessToken  Token `json:"access_token"`
		RefreshToken Token `json:"refresh_token"`
	}

	loginResponse := LoginResponse{
		AccessToken: Token{
			Token:         tokenPair.AccessToken,
			ExpirySeconds: accessTokenExpiryInt,
		},
		RefreshToken: Token{
			Token:         tokenPair.RefreshToken,
			ExpirySeconds: refreshTokenExpiryInt,
		},
	}

	utils.WriteJSON(w, http.StatusAccepted, loginResponse)
}

func (a *Auth) RefreshToken(w http.ResponseWriter, r *http.Request, db db.DBInterface) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "refresh_token" {
			claims := &auth.JWTClaims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(a.Secret), nil
			})
			if err != nil {
				utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				utils.ErrorJSON(w, errors.New("unknown User"), http.StatusUnauthorized)
				return
			}

			user, err := db.GetUserByID(userID)
			if err != nil {
				utils.ErrorJSON(w, errors.New("unknown User"), http.StatusUnauthorized)
				return
			}

			businessName := r.Header.Get("Business-Name")

			business, err := db.GetBusinessByBusinessName(businessName)
			if err != nil {
				utils.ErrorJSON(w, errors.New("invalid business"), http.StatusBadRequest)
				return
			}

			isUserInBusiness, err := db.IsUserInBusiness(user.ID, business.ID)
			if err != nil || !isUserInBusiness {
				utils.ErrorJSON(w, errors.New("user not associated with the business"), http.StatusBadRequest)
				return
			}

			u := auth.AuthUser{
				ID:         user.ID,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
				BusinessID: business.ID,
			}

			tokenPair, err := a.generateSignedTokenPair(&u)
			if err != nil {
				utils.ErrorJSON(w, errors.New("error generating token pair"), http.StatusInternalServerError)
				return
			}

			host := r.Host

			http.SetCookie(w, a.getAccessCookie(tokenPair.AccessToken, host))
			http.SetCookie(w, a.getRefreshCookie(tokenPair.RefreshToken, host))

			type RefreshResponse struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}

			refreshResponse := RefreshResponse{
				Error:   false,
				Message: "Token refreshed",
			}

			utils.WriteJSON(w, http.StatusAccepted, refreshResponse)
		}
	}
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	http.SetCookie(w, a.getExpiredRefreshCookie(domain))
	w.WriteHeader(http.StatusAccepted)
	utils.WriteJSON(w, http.StatusOK, "Logged out")
}
