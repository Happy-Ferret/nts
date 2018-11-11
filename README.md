# nts

[aka. Naive Ticketing System]

Simple web service for handling ticket sales.

Meant as a simulation on how to do it reliably, i.e. without overselling, charging anyone more than once, etc.

Also meant primarily as a backend implementation, so frontend part is just pure HTML, no style, no "dynamisms".

## usage

Commands below will start a server instance at http://localhost:9000, you can go there with your browser to try it out.

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
