// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-05-06 05:43:22 3C180D                   go-experiments/[tls_demo.go]
// -----------------------------------------------------------------------------

package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func tlsServerDemo() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add(
			"Strict-Transport-Security", "max-age=63072000; includeSubDomains",
		)
		w.Write([]byte("This is an example server.\n"))
	})
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: cfg,
		TLSNextProto: make(map[string]func(
			*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
} //                                                               tlsServerDemo

/*
import (
    "fmt"
    "log"
    "net/http"
)

func serverDemo() {
    fmt.Println("running serverDemo()")
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    redirect := func() {
        const path = "/main.html"
        fmt.Println("REDIRECT:", r.URL.Path, "TO:", path)
        http.Redirect(w, r, path, http.StatusSeeOther) // StatusSeeOther (303)
        // StatusFound (302) also works, but code 303 is more appropriate

        // StatusMultipleChoices  = 300
        // StatusMovedPermanently = 301
        // StatusFound            = 302
        // StatusSeeOther         = 303
        // StatusNotModified      = 304
        // StatusUseProxy         = 305
    }
    fmt.Println("request:", r.URL.Path)
    switch r.URL.Path {
    case "/", "/r":
        redirect()
    case "/index.html":
        http.ServeFile(w, r, "webpages/index.html")
    case "/main.html":
        http.ServeFile(w, r, "webpages/main.html")
    case "/page_one/page_one.html":
        http.ServeFile(w, r, "webpages/page_one/page_one.html")
    case "/page_one/image.png":
        http.ServeFile(w, r, "webpages/page_one/image.png")
    case "/page_two/page_two.html":
        http.ServeFile(w, r, "webpages/page_two/page_two.html")
    case "/page_two/image.png":
        http.ServeFile(w, r, "webpages/page_two/image.png")
    case "/favicon.ico":
        // ignore
    default:
        fmt.Println("not handled:", r.URL.Path)
    }
}
*/

// end
