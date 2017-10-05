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

	err = disableSSHDService()

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service disabled")
	}

	err = stopSSHDService()

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service stopped")
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

	err = enableSSHDService()

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service enabled")
	}

	log.Info("Starting sshd service...")

	err = startSSHDService()

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("sshd service started")
	}

	if knf.HasProp(SCRIPT_END) {
		runScript(knf.GetS(SCRIPT_END))
	}

	// Shutdown Bastion when bastion mode disabled
	log.Info("Bastion now is shutdown...")

	shutdown(0)
}

// stopSSHDService stop sshd daemon
func stopSSHDService() error {
	if initsystem.Systemd() {
		return stopSSHDServiceBySystemd()
	}

	return stopSSHDServiceBySysV()
}

// stopSSHDService stop sshd daemon by systemd
func stopSSHDServiceBySystemd() error {
	err := exec.Command("systemctl", "stop", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't stop sshd service through systemd")
	}

	if isServiceStopped("sshd") {
		return nil
	}

	return fmt.Errorf("sshd service still works after 15 sec")
}

// stopSSHDService stop sshd daemon by sysv
func stopSSHDServiceBySysV() error {
	err := exec.Command("service", "sshd", "stop").Start()

	if err != nil {
		return fmt.Errorf("Can't stop sshd service through sysv")
	}

	if isServiceStopped("sshd") {
		return nil
	}

	return fmt.Errorf("sshd service still works after 15 sec")
}

// disableSSHDService disable autostart for sshd daemon
func disableSSHDService() error {
	if initsystem.Systemd() {
		return disableSSHDServiceBySystemd()
	}

	return disableSSHDServiceBySysV()
}

// disableSSHDService disable autostart for sshd daemon by systemd
func disableSSHDServiceBySystemd() error {
	err := exec.Command("systemctl", "disable", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't disable sshd service through systemd")
	}

	enabled, err := initsystem.IsEnabled("sshd")

	if err != nil || enabled {
		return fmt.Errorf("Can't disable sshd service through systemd")
	}

	return nil
}

// disableSSHDService disable autostart for sshd daemon by sysv
func disableSSHDServiceBySysV() error {
	err := exec.Command("chkconfig", "--del", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't disable sshd service through sysv")
	}

	enabled, err := initsystem.IsEnabled("sshd")

	if err != nil || enabled {
		return fmt.Errorf("Can't disable sshd service through sysv")
	}

	return nil
}

// stopSSHDService start sshd daemon
func startSSHDService() error {
	if initsystem.Systemd() {
		return startSSHDServiceBySystemd()
	}

	return startSSHDServiceBySysV()
}

// stopSSHDService start sshd daemon
func startSSHDServiceBySystemd() error {
	err := exec.Command("systemctl", "start", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't start sshd service through systemd")
	}

	if isServiceWorks("sshd") {
		return nil
	}

	return fmt.Errorf("sshd service still stopped after 15 sec")
}

// stopSSHDService start sshd daemon
func startSSHDServiceBySysV() error {
	err := exec.Command("service", "sshd", "start").Start()

	if err != nil {
		return fmt.Errorf("Can't stop sshd service through sysv")
	}

	if isServiceWorks("sshd") {
		return nil
	}

	return fmt.Errorf("sshd service still stopped after 15 sec")
}

// disableSSHDService enable autostart for sshd daemon
func enableSSHDService() error {
	if initsystem.Systemd() {
		return enableSSHDServiceBySystemd()
	}

	return enableSSHDServiceBySysV()
}

// disableSSHDService enable autostart for sshd daemon
func enableSSHDServiceBySystemd() error {
	err := exec.Command("systemctl", "enable", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't enable sshd service through systemd")
	}

	enabled, err := initsystem.IsEnabled("sshd")

	if err != nil || !enabled {
		return fmt.Errorf("Can't enable sshd service through systemd")
	}

	return nil
}

// disableSSHDService enable autostart for sshd daemon
func enableSSHDServiceBySysV() error {
	err := exec.Command("chkconfig", "--add", "sshd").Start()

	if err != nil {
		return fmt.Errorf("Can't enable sshd service through sysv")
	}

	enabled, err := initsystem.IsEnabled("sshd")

	if err != nil || enabled {
		return fmt.Errorf("Can't enable sshd service through sysv")
	}

	return nil
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
