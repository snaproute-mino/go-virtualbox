package virtualbox

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	VBOXM   string // Path to VBoxManage utility.
	Verbose bool   // Verbose mode.
)

func init() {
	VBOXM = "VBoxManage"
	if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" && runtime.GOOS == "windows" {
		VBOXM = filepath.Join(p, "VBoxManage.exe")
	}
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrVBMNotFound     = errors.New("VBoxManage not found")
)

func VBM(args ...string) error {
	cmd := exec.Command(VBOXM, args...)
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", VBOXM, strings.Join(args, " "))
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrVBMNotFound
		}
		return err
	}
	return nil
}

func VBMOut(args ...string) (string, error) {
	cmd := exec.Command(VBOXM, args...)
	if Verbose {
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", VBOXM, strings.Join(args, " "))
	}

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return string(b), err
}

func VBMOutErr(args ...string) (string, string, error) {
	cmd := exec.Command(VBOXM, args...)
	if Verbose {
		log.Printf("executing: %v %v", VBOXM, strings.Join(args, " "))
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}
