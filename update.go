package main

import (
	"log"
	"os/exec"
	"os"
	"github.com/deanishe/awgo"
)

// doUpdate checks for a newer version of the workflow.
func doUpdate() error {
	log.Println("Checking for update...")
	return wf.CheckForUpdate()
}

// checkForUpdate runs "./bwg update" in the background if an update check is due.
func checkForUpdate() error {
	if !wf.UpdateCheckDue() || wf.IsRunning("update") {
		return nil
	}
	cmd := exec.Command(os.Args[0], "update")
	return wf.RunInBackground("update", cmd)
}

// showUpdateStatus adds an "update available!" message to Script Filters if an update is available
// and query is empty.
func showUpdateStatus() {
	if len(query) != 0 {
		return
	}
	if wf.UpdateAvailable() {
		wf.Configure(aw.SuppressUIDs(true))
		log.Println("Update available!")
		wf.NewItem("An update is available!").
			Subtitle("⇥ or ↩ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(aw.IconSync)
	}
}


