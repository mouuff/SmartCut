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

const SmartCutVersion string = "v1.0.5"

const PubStr string = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAyVY0SJrSdohlXCp8kXhL
89Z5BJxO16uZiwnRaVB2AgVZW9BEsPiVhhPq7PnqDe+0LPmM7ovZZ2h2W6AZk2l6
GC+Ou5W23le+SiBc8BO8TKo7NhWc7UFL35B5V0g4ugZQxaMFo6SNOGnmLUEPi/Ny
rp3DDqNSQOn5jS8SFJICwlpcz/pNOmO3nGZOu/whbGWJ81dDydkubTzN8iMVKGs9
X8YrOR7TK7xvQ1C4nsL3J5RNegYEsxIaUAJSJi3ct1w/dOVfh/7ikHVOIlPsfy6c
7VOPS9gJi0ujsbCAKbrxdL5o7AJaXRyGfOvnGzVUt/fcBSBmJ/0+D5vf/4JZ7+be
Oq10InGeBoTCqBUDekcRokq7xOfa3zCxsFbqEJ1IzMWde+R80u2d8joDZ5btw1/U
MCQ9Rk/Ki8rP0FVSwjJKDY/bfkux3B7CaRywfwfX2TplzpH9Rzsw1QRIL7mcOpzV
9Ugj+XIcSzNqkkQ3qwESpTK2jC9NdET8K4+cv0Cl6EEtVFAxFNn6Vw6IoXtkN7Yi
sqWMeaW6rMFSgzDl0Iamfv2Oz40exd3Z8Efz+Hd79CNXqMeEc0VxeXQcrlv5seC6
hOJuc17JlOSZ+14TdrAqXnhsaeZMP0yy1xnHYzW8B0/H3nZwf3tCSLLxizOz8hb/
VTvb2H4zjs77L8w20V4rDi0CAwEAAQ==
-----END RSA PUBLIC KEY-----`

type SmartCutUpdater struct {
	updater *updater.Updater
}

func GetSmartCutUpdater() *SmartCutUpdater {
	return &SmartCutUpdater{
		updater: &updater.Updater{
			Provider: &provider.Secure{
				BackendProvider: &provider.Github{
					RepositoryURL: "github.com/mouuff/SmartCut",
					ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
				},
				PublicKeyPEM: []byte(PubStr),
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
