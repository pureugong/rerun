package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type processManager struct {
	conf  *config
	oscmd *exec.Cmd
}

func (pm *processManager) formatBuildTime(duration time.Duration) string {
	return fmt.Sprintf("%.2f(s)", duration.Seconds())
}

func (pm *processManager) run() {
	logger.Debug("building application...")

	start := time.Now()

	os.Remove(pm.conf.build)

	// Join build options with spaces
	buildOpts := strings.Join(pm.conf.buildOpts, " ")

	// Print build command and options for debugging
	fmt.Println("Build command:", pm.conf.build)
	fmt.Println("Build options:", pm.conf.buildOpts)

	// Create a slice with command arguments
	args := append([]string{"build"}, strings.Fields(buildOpts)...)
	args = append(args, "-o", pm.conf.build)

	// Execute the build command with options
	out, err := exec.Command("go", args...).CombinedOutput()

	if err != nil {
		logger.Errorf("build failed! %s", err.Error())
		fmt.Printf("%s", out)
		return
	}

	// build success, display build time
	logger.Infof("build took %s", pm.formatBuildTime(time.Since(start)))

	if pm.conf.Test {
		testOut, testErr := exec.Command("go", "test").CombinedOutput()
		if testErr != nil {
			logger.Error("Tests failed!")
			fmt.Printf("==========\n%s==========\n", testOut)
		} else {
			logger.Info("Tests OK!")
		}
	}

	pm.oscmd = exec.Command(pm.conf.build, pm.conf.Args...)
	pm.oscmd.Stdout = os.Stdout
	pm.oscmd.Stdin = os.Stdin
	pm.oscmd.Stderr = os.Stderr

	logger.Debugf("starting application with arguments: %v", pm.conf.Args)
	err = pm.oscmd.Start()
	if err != nil {
		logger.Errorf("error while starting application! %s", err.Error())
	}
}

func (pm *processManager) stop() {
	logger.Debug("stopping application")

	if pm.oscmd == nil {
		return
	}

	err := pm.oscmd.Process.Kill()
	if err != nil {
		logger.Errorf("error while stopping application! %s", err.Error())
	}
}
