package routing

import (
	"net/http"

	"github.com/krqa/nts/db"
)

func debugReset(dbC db.TicketsDB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := dbC.Reset(); err != nil {
			handleInternalError(res, "resetting the database", err)
		}

		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}
