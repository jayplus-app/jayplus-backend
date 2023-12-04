package auth

import (
	"backend/contracts/db"
	"backend/utils"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *Auth) AuthRequired(db db.DBInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := a.getHeaderFromTokenAndVerify(w, r)
			if err != nil {
				utils.ErrorJSON(w, err, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)

			tokenBusinessID := claims.BusinessID

			requestBusinessName := r.Header.Get("Business-Name")
			business, err := db.GetBusinessByBusinessName(requestBusinessName)
			if err != nil {
				utils.ErrorJSON(w, err, http.StatusBadRequest)
				return
			}

			if tokenBusinessID != business.ID {
				utils.ErrorJSON(w, errors.New("invalid business"), http.StatusBadRequest)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				utils.ErrorJSON(w, err, http.StatusBadRequest)
				return
			}

			isUserInBusiness, err := db.IsUserInBusiness(userID, business.ID)
			if err != nil || !isUserInBusiness {
				utils.ErrorJSON(w, errors.New("user not associated with the business"), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
