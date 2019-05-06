// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-05-06 06:19:33 DA0ECF                go-experiments/[server_demo.go]
// -----------------------------------------------------------------------------

package main

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

	/*
		301 StatusMovedPermanently:
		A permanent redirect. Clients should make future requests using the
		new URL and the old URL is obsolete. Clients should only redirect
		automatically for GET requests. This redirect is cacheable.

		302 StatusFound:
		A temporary redirect. Clients should
		check the original URL in the future.

		303 StatusSeeOther:
	*/

	/*
		/// SORT OUT THESE NOTES:
		----------------------------------------------------------------------

		302: Redirect for undefined reason. Clients making subsequent
		requests for this resource should not use the new URI. Clients
		should not follow the redirect automatically for POST/PUT/DELETE
		requests.

		----------------------------------------------------------------------

		A 303 redirect is the same as a 302 except that the follow-up
		request is now explicitly changed to a GET request and no
		confirmation is required. 303 is meant to redirect a POST request
		to a GET resource (otherwise, the client assumes that the request
		method for the new location is the same as for the original
		resource).

		303: Redirect for undefined reason. Typically, 'Operation has
		completed, continue elsewhere.' Clients making subsequent requests
		for this resource should not use the new URI. Clients should
		follow the redirect for POST/ PUT/DELETE requests, but use GET
		for the follow-up request.

		303: See other. The request is received correctly, but the results
		should be retrieved using a GET on the redirect url.

		303: The request is received correctly. Any PUT requests are
		processed. The resulting document can be retrieved from the
		redirect url. Future request should still go to the original url.

		----------------------------------------------------------------------

		307: Temporary redirect. Resource may return to this location at
		a later point. Clients making subsequent requests for this
		resource should use the old URI. Clients should not follow the
		redirect automatically for POST/PUT/DELETE requests.

		307: Temporary redirect. The entire request should be redirected
		to the new url. Any post data should be re-posted.

		----------------------------------------------------------------------

		If you're redirecting a client as part of your web application but
		expect them to always start at the web application (for example,
		a URL shortener), a 302 redirect seems to make sense. A 303
		redirect is for use when you are receiving POST data from a
		client (e.g., a form submission) and you want to redirect them to
		a new web page to be retrieved using GET instead of POST (e.g.,
		a standard page request).

		But see this note from the status code definitions -- most clients
		will do the same thing for either a 302 or 303:

		Note: RFC 1945 and RFC 2068 specify that the client is not allowed
		to change the method on the redirected request.  However, most
		existing user agent implementations treat 302 as if it were a 303
		response, performing a GET on the Location field-value regardless
		of the original request method. The status codes 303 and 307 have
		been added for servers that wish to make unambiguously clear
		which kind of reaction is expected of the client.

		I personally recommend avoiding 302 if you have the choice. Many
		clients do not follow the spec when they encounter a 302. For
		temporary redirects, you should use either 303 or 307, depending
		on what type of behavior you want on non-GET requests. Prefer 307
		to 303 unless you need the alternate behavior on POST/PUT/DELETE.

		Note that 302 was intended to have the behavior of 307, but most
		browsers implemented it as the behavior of 303 (both of which
		didn't exist back then). Therefore, those two new codes were
		introduced to replace 302.

		The difference between 301 and 303:

		Note: Be careful with this code. Browsers and proxies tend to
		apply really aggressive caching on it, so if you reply with a 301
		it might take a long while for someone to revisit that url.
	*/

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

//end

/* ///
    http.HandleFunc("/redirectTo", targetHandler)
    http.HandleFunc("/", homeHandler)
    srv := http.Server{Addr: httpAddr}
    log.Fatal(srv.ListenAndServe())
}
*/
