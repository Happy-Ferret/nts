package routing

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/krqa/nts/db"
)

func handleInternalError(res http.ResponseWriter, op string, err error) {
	http.Error(
		res,
		fmt.Sprintf("Error while %s: `%s`", op, err),
		http.StatusInternalServerError,
	)
}

func tmplExec(res http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if err := tmpl.Execute(res, data); err != nil {
		log.Printf("template execute error: %s", err)
	}
}

// New creates new Handler that provides routing, initialized to use given
// database implementation and Stripe Public/Publishable Key.
func New(dbC db.TicketsDB, stripePublicKey string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", root(dbC))
	r.Route("/nts/v1/", func(r chi.Router) {
		r.Post("/reserve", reserve(dbC, stripePublicKey))
		r.Post("/charge", charge(dbC))
		r.Get("/list-guests", listGuests(dbC))
		r.Post("/debug/reset", debugReset(dbC))
	})
	return r
}
