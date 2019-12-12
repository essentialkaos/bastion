package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"pkg.re/essentialkaos/ek.v11/knf"
	"pkg.re/essentialkaos/ek.v11/log"

	"github.com/valyala/fasthttp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// startHTTPServer start HTTP server
func startHTTPServer(ip, port string) error {
	addr := ip + ":" + port

	log.Aux("%s %s HTTP server is started on %s", APP, VER, addr)

	server := fasthttp.Server{
		Handler: fastHTTPHandler,
		Name:    knf.GetS(SERVER_NAME, APP+"/"+VER),
	}

	return server.ListenAndServe(addr)
}

// fastHTTPHandler handler for fast http requests
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	defer requestRecover(ctx)

	path := string(ctx.Path())

	writeBasicInfo(ctx)

	if key == "" && !bastionMode {
		if path == "/go" {
			ctx.WriteString(generateSecrets())
		}

		return
	}

	if path == bastionPath && !bastionMode {
		bastionMode = true
		go startBastionMode()
	}
}

// requestRecover recover panic in request
func requestRecover(ctx *fasthttp.RequestCtx) {
	r := recover()

	if r != nil {
		log.Error("Recovered internal error in http request handler: %v", r)
		writeBasicInfo(ctx)
	}
}

// writeBasicInfo add basic info to response
func writeBasicInfo(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html; charset=UTF-8")
	ctx.SetStatusCode(200)
}
