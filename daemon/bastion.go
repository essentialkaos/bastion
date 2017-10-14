package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/initsystem"
	"pkg.re/essentialkaos/ek.v9/jsonutil"
	"pkg.re/essentialkaos/ek.v9/knf"
	"pkg.re/essentialkaos/ek.v9/log"
	"pkg.re/essentialkaos/ek.v9/pluralize"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// BASTION_MARKER path to file with info about bastion mode
const BASTION_MARKER = "/root/.bastion"

// ////////////////////////////////////////////////////////////////////////////////// //

type BastionMarker struct {
	Started int64 `json:"started"`
	Until   int64 `json:"until"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	bastionMode   bool
	bastionMarker *BastionMarker
)

// ////////////////////////////////////////////////////////////////////////////////// //

func startBastionMode() {
	duration := knf.GetI64(MAIN_DURATION, 86400)

	enableBastionMode(duration)

	until := time.Now().Unix() + duration
	count := 0

	for {
		time.Sleep(time.Minute)
		count++

		if count == 15 {
			hoursTillExit := int((until - time.Now().Unix()) / 3600)
			log.Info(
				"%s till exit from bastion mode",
				pluralize.Pluralize(hoursTillExit, "hour", "hours"),
			)
			count = 0
		}
	}

	disableBastionMode()
}

// isBastionModeEnabled return true if bastion mode is enabled
func isBastionModeEnabled() bool {
	if bastionMode {
		return true
	}

	return isBastionMarkerExist()
}

// enableBastionMode enable bastion mode on server
func enableBastionMode(duration int64) {
	if knf.HasProp(SCRIPT_BEFORE) {
		runScript(knf.GetS(SCRIPT_BEFORE))
	}

	err := createBastionMarker(duration)

	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Disabling sshd service...")

	err = disableService("sshd")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service disabled")
	}

	log.Info("Stopping sshd service...")

	err = stopService("sshd")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service stopped")
	}

	log.Info("Enabling bastion service...")

	err = enableService("bastion")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("bastion service enabled")
	}

	if knf.HasProp(SCRIPT_IN) {
		runScript(knf.GetS(SCRIPT_IN))
	}
}

// enableBastionMode disable bastion mode on server
func disableBastionMode() {
	if knf.HasProp(SCRIPT_OUT) {
		runScript(knf.GetS(SCRIPT_OUT))
	}

	err := removeBastionMarker()

	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Enabling sshd service...")

	err = enableService("sshd")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service enabled")
	}

	log.Info("Starting sshd service...")

	err = startService("sshd")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service started")
	}

	log.Info("Disabling bastion service...")

	err = disableService("bastion")

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("bastion service disabled")
	}

	if knf.HasProp(SCRIPT_END) {
		runScript(knf.GetS(SCRIPT_END))
	}

	// Shutdown Bastion when bastion mode disabled
	log.Info("Bastion now is shutdown...")

	shutdown(0)
}

// enableService enable service autostart
func enableService(name string) error {
	if initsystem.Systemd() {
		return enableServiceBySystemd(name)
	}

	return enableServiceBySysV(name)
}

// disableService disable service autostart
func disableService(name string) error {
	if initsystem.Systemd() {
		return disableServiceBySystemd(name)
	}

	return disableServiceBySysV(name)
}

// startService start service
func startService(name string) error {
	if initsystem.Systemd() {
		return startServiceBySystemd(name)
	}

	return startServiceBySysV(name)
}

// stopService stop service
func stopService(name string) error {
	if initsystem.Systemd() {
		return stopServiceBySystemd(name)
	}

	return stopServiceBySysV(name)
}

// enableServiceBySysV enable service autostart by chkconfig
func enableServiceBySysV(name string) error {
	err := exec.Command("chkconfig", "--add", name).Start()

	if err != nil {
		return fmt.Errorf("Can't enable %s service through sysv", name)
	}

	enabled, err := initsystem.IsEnabled(name)

	if err != nil || enabled {
		return fmt.Errorf("Can't enable %s service through sysv", name)
	}

	return nil
}

// enableServiceBySystemd enable service autostart by systemctl
func enableServiceBySystemd(name string) error {
	err := exec.Command("systemctl", "enable", name).Start()

	if err != nil {
		return fmt.Errorf("Can't enable %s service through systemd", name)
	}

	enabled, err := initsystem.IsEnabled(name)

	if err != nil || !enabled {
		return fmt.Errorf("Can't enable %s service through systemd", name)
	}

	return nil
}

// disableServiceBySysV disable service autostart by chkconfig
func disableServiceBySysV(name string) error {
	err := exec.Command("chkconfig", "--del", name).Start()

	if err != nil {
		return fmt.Errorf("Can't disable %s service through sysv", name)
	}

	enabled, err := initsystem.IsEnabled(name)

	if err != nil || enabled {
		return fmt.Errorf("Can't disable %s service through sysv", name)
	}

	return nil
}

// disableServiceBySystemd disable service autostart by systemctl
func disableServiceBySystemd(name string) error {
	err := exec.Command("systemctl", "disable", name).Start()

	if err != nil {
		return fmt.Errorf("Can't disable %s service through systemd", name)
	}

	enabled, err := initsystem.IsEnabled(name)

	if err != nil || enabled {
		return fmt.Errorf("Can't disable %s service through systemd", name)
	}

	return nil
}

// startServiceBySysV start service by sysv init script
func startServiceBySysV(name string) error {
	err := exec.Command("service", name, "start").Start()

	if err != nil {
		return fmt.Errorf("Can't stop %s service through sysv", name)
	}

	if isServiceWorks(name) {
		return nil
	}

	return fmt.Errorf("%s service still stopped after 15 sec", name)
}

// startServiceBySystemd start service by systemctl
func startServiceBySystemd(name string) error {
	err := exec.Command("systemctl", "start", name).Start()

	if err != nil {
		return fmt.Errorf("Can't start %s service through systemd", name)
	}

	if isServiceWorks(name) {
		return nil
	}

	return fmt.Errorf("%s service still stopped after 15 sec", name)
}

// stopServiceBySysV stop service by sysv init script
func stopServiceBySysV(name string) error {
	err := exec.Command("service", name, "stop").Start()

	if err != nil {
		return fmt.Errorf("Can't stop %s service through sysv", name)
	}

	if isServiceStopped(name) {
		return nil
	}

	return fmt.Errorf("%s service still works after 15 sec", name)
}

// stopServiceBySystemd stop service by systemctl
func stopServiceBySystemd(name string) error {
	err := exec.Command("systemctl", "stop", name).Start()

	if err != nil {
		return fmt.Errorf("Can't stop %s service through systemd", name)
	}

	if isServiceStopped(name) {
		return nil
	}

	return fmt.Errorf("%s service still works after 15 sec", name)
}

// createBastionMarker create file with info about bastion mode
func createBastionMarker(duration int64) error {
	if isBastionMarkerExist() {
		return nil
	}

	now := time.Now().Unix()
	marker := BastionMarker{now, now + duration}

	err := jsonutil.EncodeToFile(BASTION_MARKER, marker, 0600)

	if err != nil {
		return fmt.Errorf("Can't encode bastion marker: %v", err)
	}

	return nil
}

// removeBastionMarker remove file with info about bastion mode
func removeBastionMarker() error {
	if !isBastionMarkerExist() {
		return nil
	}

	err := os.Remove(BASTION_MARKER)

	if err != nil {
		return fmt.Errorf("Can't remove bastion marker: %v", err)
	}

	return nil
}

// getBastionMarkerInfo read and decode bastion marker
func getBastionMarkerInfo() (*BastionMarker, error) {
	marker := &BastionMarker{}

	err := jsonutil.DecodeFile(BASTION_MARKER, marker)

	if err != nil {
		return nil, err
	}

	return marker, nil
}

// isBastionMarkerExist return true if bastion marker file exist
func isBastionMarkerExist() bool {
	return fsutil.IsExist(BASTION_MARKER)
}

// runScript run script
func runScript(script string) {
	log.Info("Executing sctipt '%s' ...", script)

	err := exec.Command("bash", "script").Start()

	if err != nil {
		log.Error("Script return error")
	} else {
		log.Info("Script successfully executed")
	}
}

// isServiceWorks return true if service works
func isServiceWorks(name string) bool {
	for i := 0; i < 15; i++ {
		works, err := initsystem.IsServiceWorks(name)

		if err == nil && works {
			return true
		}

		time.Sleep(time.Second)
	}

	return false
}

// isServiceStopped return true if service stopped
func isServiceStopped(name string) bool {
	for i := 0; i < 15; i++ {
		works, err := initsystem.IsServiceWorks(name)

		if err == nil && !works {
			return true
		}

		time.Sleep(time.Second)
	}

	return false
}
