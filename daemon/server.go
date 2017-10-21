package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"pkg.re/essentialkaos/ek.v9/knf"
	"pkg.re/essentialkaos/ek.v9/log"

	"github.com/valyala/fasthttp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var serverName string

// ////////////////////////////////////////////////////////////////////////////////// //

// startHTTPServer start HTTP server
func startHTTPServer(ip, port string) error {
	serverName = knf.GetS(SERVER_NAME, APP+"/"+VER)

	addr := ip + ":" + port

	log.Aux("Bastion %s HTTP server is started on %s", VER, addr)

	return fasthttp.ListenAndServe(addr, fastHTTPHandler)
}

// fastHTTPHandler handler for fast http requests
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	defer requestRecover(ctx)

	path := string(ctx.Path())

	writeBasicInfo(ctx)

	if key == "" {
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
	ctx.Response.Header.Set("Server", serverName)
	ctx.SetStatusCode(200)
}
