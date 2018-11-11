package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"testing"

	"github.com/krqa/nts/db"
	"github.com/krqa/nts/routing"
)

// buildURL creates stringified URL merged from base and arbitrary amount
// of additional parts, appended in the order of specifying.
func buildURL(t *testing.T, base string, parts ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		t.Fatalf("Error parsing URL: %s", err)
	}
	for _, r := range parts {
		u, err = u.Parse(r)
		if err != nil {
			t.Fatalf("Error parsing URL: %s", err)
		}
	}
	return u.String()
}

// doRequest Performs a HTTP request with given parameters and checks whether
// it completes successfully and with given status, failing the test if not.
func doRequest(t *testing.T, client *http.Client, url string, vals url.Values, status int) *http.Response {
	res, err := client.PostForm(url, vals)
	if err != nil {
		t.Fatalf("Unexpected error return: %s", err)
	}
	if res.StatusCode != status {
		t.Fatalf("Unexpected status code: %d != %d", res.StatusCode, status)
	}
	return res
}

func assertDBState(t *testing.T, dbC *db.MemTicketsDB, remainingExpected int, customersExpected []string) {
	remaining, err := dbC.Remaining()
	if err != nil {
		t.Fatalf("Error getting remaining: %s", err)
	}
	if remaining != remainingExpected {
		t.Fatalf("Invalid remaining count: %d != %d", remaining, remainingExpected)
	}
	customers := []string{}
	for fullname := range dbC.Customers {
		customers = append(customers, fullname)
	}
	sort.Strings(customers)
	sort.Strings(customersExpected)
	if !reflect.DeepEqual(customers, customersExpected) {
		t.Fatalf("Invalid guests list: %v != %v", customers, customersExpected)
	}
}

func TestReserveOne(t *testing.T) {
	memdb := db.NewMemTicketsDB(5)
	ts := httptest.NewServer(routing.New(memdb, ""))
	defer ts.Close()
	client := ts.Client()

	_ = doRequest(
		t,
		client,
		buildURL(t, ts.URL, "/nts/v1/reserve"),
		url.Values{"fullname": {"John Doe"}},
		http.StatusOK,
	)

	assertDBState(t, memdb, 4, []string{"John Doe"})
}

func TestReserveNone(t *testing.T) {
	memdb := db.NewMemTicketsDB(0)
	ts := httptest.NewServer(routing.New(memdb, ""))
	defer ts.Close()
	client := ts.Client()

	_ = doRequest(
		t,
		client,
		buildURL(t, ts.URL, "/nts/v1/reserve"),
		url.Values{"fullname": {"John Doe"}},
		http.StatusInternalServerError,
	)

	assertDBState(t, memdb, 0, []string{})
}

func TestReserveSome(t *testing.T) {
	memdb := db.NewMemTicketsDB(5)
	ts := httptest.NewServer(routing.New(memdb, ""))
	defer ts.Close()
	client := ts.Client()

	customers := []string{"John", "Kenji", "Marcia", "Matt", "Eddie"}

	for _, fullname := range customers {
		_ = doRequest(
			t,
			client,
			buildURL(t, ts.URL, "/nts/v1/reserve"),
			url.Values{"fullname": {fullname}},
			http.StatusOK,
		)
	}

	assertDBState(t, memdb, 0, customers)

	_ = doRequest(
		t,
		client,
		buildURL(t, ts.URL, "/nts/v1/reserve"),
		url.Values{"fullname": {"Josh"}},
		http.StatusInternalServerError,
	)

	assertDBState(t, memdb, 0, customers)
}

func TestReserveTwice(t *testing.T) {
	memdb := db.NewMemTicketsDB(5)
	ts := httptest.NewServer(routing.New(memdb, ""))
	defer ts.Close()
	client := ts.Client()

	for _, fullname := range []string{"Kenji", "Kenji"} {
		_ = doRequest(
			t,
			client,
			buildURL(t, ts.URL, "/nts/v1/reserve"),
			url.Values{"fullname": {fullname}},
			http.StatusOK,
		)
	}

	assertDBState(t, memdb, 4, []string{"Kenji"})
}
