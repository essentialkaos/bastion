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
	"time"

	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/jsonutil"
	"pkg.re/essentialkaos/ek.v9/log"
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

var bastionMode bool

// ////////////////////////////////////////////////////////////////////////////////// //

// isBastionModeEnabled return true if bastion mode is enabled
func isBastionModeEnabled() bool {
	if bastionMode {
		return true
	}

	return isBastionMarkerExist()
}

// enableBastionMode enable bastion mode on server
func enableBastionMode(duration int64) {
	return
}

// enableBastionMode disable bastion mode on server
func disableBastionMode() {

	// Shutdown Bastion when bastion mode disabled
	log.Info("Bastion now is shutdown...")
	shutdown(0)
}

// stopSSHDService stop sshd daemon
func stopSSHDService() {
	return
}

// disableSSHDService disable autostart for sshd daemon
func disableSSHDService() {
	return
}

// stopSSHDService start sshd daemon
func startSSHDService() {

}

// disableSSHDService enable autostart for sshd daemon
func enableSSHDService() {
	return
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
