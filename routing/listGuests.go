package routing

import (
	"html/template"
	"net/http"

	"github.com/krqa/nts/db"
)

var listGuestshtml = `
<html>
<head>
</head>
<body>

{{range .}}
{{.}}<br>
{{end}}

</body>
</html>
`

var listGueststmpl = template.Must(template.New("listGuests").Parse(listGuestshtml))

func listGuests(dbC db.TicketsDB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		guests, err := dbC.Guests()
		if err != nil {
			handleInternalError(res, "getting guests list from the database", err)
			return
		}

		tmplExec(res, listGueststmpl, guests)
	}
}
