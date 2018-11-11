package routing

import (
	"html/template"
	"net/http"

	"github.com/krqa/nts/db"
)

var roothtml = `
<html>
<head>
</head>
<body>

don't miss your ticket for the best ever show<br>
only {{.}} left
<form action="/nts/v1/reserve" method="POST">
	your name:
	<input type="text" name="fullname">
	<input type="submit" value="reserve ticket for 5 minutes">
</form>

<a href="/nts/v1/list-guests">guests list</a>
<form action="/nts/v1/debug/reset" method="post">
	<input type="submit" value="reset the world">
</form>

</body>
</html>
`

var roottmpl = template.Must(template.New("root").Parse(roothtml))

func root(dbC db.TicketsDB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		remaining, err := dbC.Remaining()
		if err != nil {
			handleInternalError(res, "getting remaining number from the database", err)
			return
		}

		tmplExec(res, roottmpl, remaining)
	}
}
