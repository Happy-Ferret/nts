package routing

import (
	"html/template"
	"net/http"

	"github.com/krqa/nts/db"
)

var reservehtml = `
<html>
<head>
</head>
<body>

thank you {{.Fullname}}<br>
your ticket has been reserved for ~5 minutes<br>
please use the button below to finalize the purchase
<form action="/nts/v1/charge" method="POST">
	<input type="hidden" name="fullname" value="{{.Fullname}}">
	<script
		src="https://checkout.stripe.com/checkout.js" class="stripe-button"
		data-key="{{.PublicKey}}"
		data-amount="999"
		data-currency="GBP"
		data-name="NTS"
		data-description="Show Ticket"
		data-image="https://stripe.com/img/documentation/checkout/marketplace.png"
		data-locale="auto"
		data-zip-code="true">
	</script>
</form>

</body>
</html>
`

var reservetmpl = template.Must(template.New("reserve").Parse(reservehtml))

func reserve(dbC db.TicketsDB, stripePublicKey string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			handleInternalError(res, "parsing form data", err)
			return
		}
		fullname := req.Form.Get("fullname")

		if err := dbC.Reserve(fullname); err != nil {
			handleInternalError(res, "storing the reservation", err)
			return
		}

		tmplExec(res, reservetmpl, map[string]string{
			"Fullname":  fullname,
			"PublicKey": stripePublicKey,
		})
	}
}
