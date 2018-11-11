package routing

import (
	"html/template"
	"net/http"

	"github.com/krqa/nts/db"
	stripe "github.com/stripe/stripe-go"
	stripeCharge "github.com/stripe/stripe-go/charge"
	stripeCustomer "github.com/stripe/stripe-go/customer"
)

var chargehtml = `
<html>
<head>
</head>
<body>

congratulations {{.}} your payment is complete<br>
enjoy the show!<br>

<a href="/nts/v1/list-guests">guests list</a>
<form action="/nts/v1/debug/reset" method="POST">
	<input type="submit" value="reset the world">
</form>

</body>
</html>
`

var chargetmpl = template.Must(template.New("charge").Parse(chargehtml))

func charge(dbC db.TicketsDB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			handleInternalError(res, "parsing form data", err)
			return
		}

		customerParams := &stripe.CustomerParams{
			Email: stripe.String(req.Form.Get("stripeEmail")),
		}
		customerParams.SetSource(req.Form.Get("stripeToken"))

		customer, err := stripeCustomer.New(customerParams)
		if err != nil {
			handleInternalError(res, "getting stripe customer", err)
			return
		}

		fullname := req.Form.Get("fullname")

		if err := dbC.Charge(fullname); err != nil {
			if err == db.ErrAlreadyCharged {
				tmplExec(res, chargetmpl, fullname)
				return
			}
			handleInternalError(res, "storing charge state", err)
			return
		}

		chargeParams := &stripe.ChargeParams{
			Amount:      stripe.Int64(999),
			Currency:    stripe.String(string(stripe.CurrencyGBP)),
			Description: stripe.String("Show Ticket"),
			Customer:    stripe.String(customer.ID),
		}
		if _, err := stripeCharge.New(chargeParams); err != nil {
			handleInternalError(res, "performing stripe charge", err)
			return
		}

		tmplExec(res, chargetmpl, fullname)
	}
}
