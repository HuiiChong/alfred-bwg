//
// Copyright (c) 2016 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2016-11-08
//

package aw

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/deanishe/awgo/util"
)

// ErrJobExists is the error returned by RunInBackground if a job with
// the given name is already running.
type ErrJobExists struct {
	Name string // Name of the job
	Pid  int    // PID of the running job
}

// Error implements error interface.
func (a ErrJobExists) Error() string {
	return fmt.Sprintf("Job '%s' already running with PID %d", a.Name, a.Pid)
}

// IsJobExists returns true if error is of type ErrJobExists.
func IsJobExists(err error) bool {
	_, ok := err.(ErrJobExists)
	return ok
}

// RunInBackground executes cmd in the background. It returns an
// ErrJobExists error if a job of the same name is already running.
func (wf *Workflow) RunInBackground(jobName string, cmd *exec.Cmd) error {
	if wf.IsRunning(jobName) {
		pid, _ := wf.getPid(jobName)
		return ErrJobExists{jobName, pid}
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Prevent process from being killed when parent is
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		return err
	}

	return wf.savePid(jobName, cmd.Process.Pid)
}

// Kill stops a background job.
func (wf *Workflow) Kill(jobName string) error {
	pid, err := wf.getPid(jobName)
	if err != nil {
		return err
	}
	p := wf.pidFile(jobName)
	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		// Delete stale PID file
		os.Remove(p)
		return err
	}
	os.Remove(p)
	return nil
}

// IsRunning returns true if a job with name jobName is currently running.
func (wf *Workflow) IsRunning(jobName string) bool {
	pid, err := wf.getPid(jobName)
	if err != nil {
		return false
	}
	if err = syscall.Kill(pid, 0); err != nil {
		// Delete stale PID file
		os.Remove(wf.pidFile(jobName))
		return false
	}
	return true
}

// Save PID a job-specific file.
func (wf *Workflow) savePid(jobName string, pid int) error {
	p := wf.pidFile(jobName)
	return ioutil.WriteFile(p, []byte(strconv.Itoa(pid)), 0600)
}

// Return PID for job.
func (wf *Workflow) getPid(jobName string) (int, error) {
	p := wf.pidFile(jobName)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// Path to PID file for job.
func (wf *Workflow) pidFile(jobName string) string {
	dir := util.MustExist(filepath.Join(wf.awCacheDir(), "jobs"))
	return filepath.Join(dir, jobName+".pid")
}
