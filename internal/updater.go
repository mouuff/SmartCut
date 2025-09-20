package internal

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

const SmartCutVersion string = "v1.0.3"

type SmartCutUpdater struct {
	updater *updater.Updater
}

func GetSmartCutUpdater() *SmartCutUpdater {
	return &SmartCutUpdater{
		updater: &updater.Updater{
			Provider: &provider.Github{
				RepositoryURL: "github.com/mouuff/SmartCut",
				ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
			},
			ExecutableName: fmt.Sprintf("smartcut_%s_%s", runtime.GOOS, runtime.GOARCH),
			Version:        SmartCutVersion,
		},
	}
}

// verifyInstallation verifies if the executable is installed correctly
// we are going to run the newly installed program by running it with -version
// if it outputs the good version then we assume the installation is good
func (u *SmartCutUpdater) verifyInstallation() error {
	latestVersion, err := u.updater.GetLatestVersion()
	if err != nil {
		return err
	}
	executable, err := u.updater.GetExecutable()
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path: executable,
		Args: []string{executable, "-version"},
	}
	// Should be replaced with Output() as soon as test project is updated
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	strOutput := string(output)

	if !strings.Contains(strOutput, latestVersion) {
		return errors.New("version not found in program output")
	}

	return nil
}

func (u *SmartCutUpdater) SelfUpdate() error {
	updateStatus, err := u.updater.Update()

	if err != nil {
		return err
	}

	if updateStatus == updater.Updated {
		if err := u.verifyInstallation(); err != nil {
			return u.updater.Rollback()
		}
	}

	if updateStatus == updater.UpToDate {
		if err := u.updater.CleanUp(); err != nil {
			return err
		}
	}

	return nil
}
