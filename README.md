# nts

[aka. Naive Ticketing System]

Simple web service for handling ticket sales.

Meant as a simulation on how to do it reliably, i.e. without overselling, charging anyone more than once, etc.

Also meant primarily as a backend implementation, so frontend part is just pure HTML, no style, no "dynamisms".

## usage

Commands below will start a server instance at http://localhost:9000, you can go there with your browser to try it out.

### stripe

To be able to perform actual charge, Stripe keys need to be provided. NTS reads them from `NTS_STRIPE_PUBLISHABLE_KEY` and `NTS_STRIPE_SECRET_KEY` environment variables. Note that the same variables are also picked up by `docker-compose`.

### using docker-compose

```bash
$ docker-compose up
```

### manually

*Requires*: Go

*Locally:*

```bash
$ go build
$ ./nts
```

*Globally:*

```bash
$ go install
$ nts
```

## "ceveats"

[aka. "self code review" :-).]

*In the order of appearing in my mind.*

* Database is in memory, so restarts loose everything. An interface is defined to give possibility of implementing alternative database layers.
* State atomicity is in memory as well, therefore if there was more than one instance, race conditions could occur. Can be redeemed by using DB specific atomic transactions solution (if available) or perhaps something like Etcd.
* Charge state is stored in our database before performing the charge. This can potentially cause a situation, where we think we've got the charge, but it actually failed. We could 1) try to revoke the charge state in case of charge failure, 2) store failure state, so we can check it later, if user retries the payment. In worst case, "please contact support, but we did not charge you twice" :-].
* Obviously, number of remaining tickets might lie at times and user can be left in the cold after trying to reserve his. Guess it can be worked-around with better UI?
* Some implementation "deficiences":
	* We're contacting Stripe every time on `/charge`, regardless of whether we actually perform charge or not. It is to get customer ID from them, but maybe there's a way to avoid it?
	* Things like amount to charge, currency, etc. are stored separately in `/charge` code and in form displayed to the user.
