package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/mouuff/SmartCut/pkg/brain"
	"github.com/mouuff/SmartCut/pkg/controller"
	"github.com/mouuff/SmartCut/pkg/reader"
	"github.com/mouuff/SmartCut/pkg/types"
	"github.com/mouuff/SmartCut/pkg/utils"
	"github.com/mouuff/SmartCut/pkg/view"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"golang.design/x/clipboard"
)

const SmartCutVersion string = "v1.0.2"

// verifyInstallation verifies if the executable is installed correctly
// we are going to run the newly installed program by running it with -version
// if it outputs the good version then we assume the installation is good
func verifyInstallation(u *updater.Updater) error {
	latestVersion, err := u.GetLatestVersion()
	if err != nil {
		return err
	}
	executable, err := u.GetExecutable()
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
		return errors.New("Version not found in program output")
	}
	return nil
}

func selfUpdate() error {
	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/SmartCut",
			ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("smartcut_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        SmartCutVersion,
	}

	updateStatus, err := u.Update()

	if err != nil {
		return err
	}

	if updateStatus == updater.Updated {
		if err := verifyInstallation(u); err != nil {
			return u.Rollback()
		}
	}

	if updateStatus == updater.UpToDate {
		if err := u.CleanUp(); err != nil {
			return err
		}
	}

	return nil
}

type SmartCutCmd struct {
	flagSet *flag.FlagSet

	config      string
	versionFlag bool
}

// Init initializes the command
func (cmd *SmartCutCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet("smartcut", flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file override")
	cmd.flagSet.BoolVar(&cmd.versionFlag, "version", false, "prints the version and exit")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *SmartCutCmd) Run(a fyne.App, w fyne.Window) error {
	config, err := utils.GetOrCreateConfiguration(cmd.config)
	if err != nil {
		return err
	}

	b, err := brain.NewOllamaBrain(config.HostUrl, config.Model)
	if err != nil {
		return err
	}

	err = clipboard.Init()
	if err != nil {
		return err
	}

	// ir := reader.NewClipboardReader(ctx.Background())
	ir := reader.NewShortcutReader()
	ir.Start()

	// MVC setup
	m := types.NewSmartCutModel(config)
	v := view.NewSmartCutView(w, m)
	c := controller.NewSmartCutController(context.Background(), b, m, config)

	// Setup View / Controller hooks
	ir.OnInput = c.GenerateForInput
	v.OnRequestGenerate = c.GenerateForInput
	c.OnRequestFocus = v.RequestFocus

	w.SetContent(v.Layout())
	return nil
}

func main() {
	cmd := &SmartCutCmd{}
	err := cmd.Init(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if cmd.versionFlag {
		fmt.Println(SmartCutVersion)
		return
	}

	a := app.NewWithID("com.mouuff.smartcut")
	w := a.NewWindow("SmartCut - " + SmartCutVersion)
	w.Resize(fyne.NewSize(800, 400))
	err = cmd.Run(a, w)

	if err != nil {
		dialog.ShowError(err, w)
	}

	w.ShowAndRun()

	err = selfUpdate()
	if err != nil {
		log.Println(err)
	}
}
