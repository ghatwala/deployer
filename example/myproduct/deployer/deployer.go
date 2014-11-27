package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dorzheh/deployer/deployer"
	"github.com/dorzheh/deployer/example/myproduct"
	gui "github.com/dorzheh/deployer/ui"
	"github.com/dorzheh/deployer/ui/dialog_ui"
)

var (
	defaultProductName = ""
	prodType           = ""
	rootDir            = ""
	arch               = ""
)

// Initialize the stuff before the main() is executed.
// Setting default product name , path to the installer ,
// adding internal path to bin
func init() {
	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i686"
	}
	var absPath string
	absPath, _ = filepath.Abs(os.Args[0])
	var pathArr []string = strings.Split(absPath, "/")
	rootDir = strings.Join(pathArr[0:len(pathArr)-4], "/")
	newPath := "$PATH:" + rootDir + "/install/" + arch + "/bin"
	os.Setenv("PATH", (os.ExpandEnv(newPath)))
	buf, err := ioutil.ReadFile(filepath.Join(rootDir, ".product"))
	if err != nil {
		panic(err)
	}

	for _, str := range strings.Split(string(buf), "\n") {
		res := strings.Split(str, "=")
		switch res[0] {
		case "Product":
			prodType = res[1]
		case "Name":
			defaultProductName = res[1]
		default:
		}
	}
}

func main() {
	ui := dialog_ui.NewDialogUi()
	ui.Shadow(false)
	ui.SetCancelLabel("Exit")
	gui.UiValidateUser(ui, 0)
	gui.UiWelcomeMsg(ui, defaultProductName)
	data := &deployer.CommonData{
		RootDir:          rootDir,
		RootfsMp:         filepath.Join(rootDir, "rootfs_mnt"),
		DefaultExportDir: rootDir,
		VaName:           defaultProductName,
		Arch:             arch,
		Ui:               ui,
	}

	gui.UiDeploymentResult(ui, "MyProduct installation completed successfully", myproduct.Deploy(data, prodType))
}