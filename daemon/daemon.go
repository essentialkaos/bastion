package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"runtime"

	"pkg.re/essentialkaos/ek.v11/fmtc"
	"pkg.re/essentialkaos/ek.v11/knf"
	"pkg.re/essentialkaos/ek.v11/log"
	"pkg.re/essentialkaos/ek.v11/options"
	"pkg.re/essentialkaos/ek.v11/pid"
	"pkg.re/essentialkaos/ek.v11/signal"
	"pkg.re/essentialkaos/ek.v11/usage"

	knfv "pkg.re/essentialkaos/ek.v11/knf/validators"
	knff "pkg.re/essentialkaos/ek.v11/knf/validators/fs"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "Bastion"
	VER  = "0.0.2"
	DESC = "Utility for temporary access limitation to server"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Daemon info
const (
	MAIN_DURATION = "main:duration"
	MAIN_URL      = "main:url"
	MAIN_PATH     = "main:path"
	SERVER_IP     = "server:ip"
	SERVER_PORT   = "server:port"
	SERVER_NAME   = "server:name"
	LOG_DIR       = "log:dir"
	LOG_FILE      = "log:file"
	LOG_PERMS     = "log:perms"
	LOG_LEVEL     = "log:level"
	SCRIPT_BEFORE = "script:before"
	SCRIPT_IN     = "script:in"
	SCRIPT_OUT    = "script:out"
	SCRIPT_END    = "script:complete"
)

// Options
const (
	OPT_CONFIG   = "c:config"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// Pid info
const PID_FILE = "bastion.pid"

// ////////////////////////////////////////////////////////////////////////////////// //

// Options map
var optMap = options.Map{
	OPT_CONFIG:   {Value: "/etc/bastion.knf"},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

var key string
var bastionPath string

// ////////////////////////////////////////////////////////////////////////////////// //

// Main function
func Init() {
	runtime.GOMAXPROCS(4)

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) {
		showUsage()
		return
	}

	loadConfig()
	validateConfig()
	registerSignalHandlers()
	setupLogger()
	createPidFile()

	if isBastionModeEnabled() {
		restoreBastionMode()
	} else {
		startHTTPServer(
			knf.GetS(SERVER_IP),
			knf.GetS(SERVER_PORT),
		)
	}

	shutdown(0)
}

// loadConfig read and parse configuration file
func loadConfig() {
	err := knf.Global(options.GetS(OPT_CONFIG))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// validateConfig validate configuration file values
func validateConfig() {
	errs := knf.Validate([]*knf.Validator{
		{SERVER_PORT, knfv.Empty, nil},
		{LOG_DIR, knfv.Empty, nil},
		{LOG_FILE, knfv.Empty, nil},

		{LOG_DIR, knff.Perms, "DW"},
		{LOG_DIR, knff.Perms, "DX"},
		{LOG_LEVEL, knfv.NotContains, []string{"debug", "info", "warn", "error", "crit"}},

		{MAIN_DURATION, knfv.Less, 3600},
		{MAIN_DURATION, knfv.Greater, 604800},

		{SCRIPT_BEFORE, knff.Perms, "FS"},
		{SCRIPT_BEFORE, knff.Perms, "FX"},
		{SCRIPT_IN, knff.Perms, "FS"},
		{SCRIPT_IN, knff.Perms, "FX"},
		{SCRIPT_OUT, knff.Perms, "FS"},
		{SCRIPT_OUT, knff.Perms, "FX"},
		{SCRIPT_END, knff.Perms, "FS"},
		{SCRIPT_END, knff.Perms, "FX"},
	})

	if len(errs) != 0 {
		printError("Error while configuration file validation:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}
}

// registerSignalHandlers register signal handlers
func registerSignalHandlers() {
	signal.Handlers{
		signal.TERM: termSignalHandler,
		signal.INT:  intSignalHandler,
		signal.HUP:  hupSignalHandler,
	}.TrackAsync()
}

// setupLogger setup logger
func setupLogger() {
	err := log.Set(knf.GetS(LOG_FILE), knf.GetM(LOG_PERMS, 644))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	err = log.MinLevel(knf.GetS(LOG_LEVEL))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// createPidFile create PID file
func createPidFile() {
	err := pid.Create(PID_FILE)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// INT signal handler
func intSignalHandler() {
	log.Aux("Received INT signal, shutdown...")
	shutdown(0)
}

// TERM signal handler
func termSignalHandler() {
	log.Aux("Received TERM signal, shutdown...")
	shutdown(0)
}

// HUP signal handler
func hupSignalHandler() {
	log.Info("Received HUP signal, log will be reopened...")
	log.Reopen()
	log.Info("Log reopened by HUP signal")
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// shutdown stop deamon
func shutdown(code int) {
	pid.Remove(PID_FILE)
	os.Exit(code)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo()

	info.AddOption(OPT_CONFIG, "Path to config file", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	about.Render()
}
