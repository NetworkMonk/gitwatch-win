package watch

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"time"

	"github.com/NetworkMonk/gitwatch/config"
)

// This package contains everything that is required to perform regular git requests on a folder

// Start will trigger watches on all specified configurations, this function will not return unless there is an error
func Start(layout *config.Layout) error {
	if layout == nil {
		return errors.New("No valid configuration")
	}

	for entryIndex := 0; entryIndex < len(layout.Watch); entryIndex++ {
		go watchEntry(&layout.Watch[entryIndex])
	}

	log.Printf("Watching %d folders", len(layout.Watch))

	// Block forever
	for {
	}
}

// watchEntry will watch a single config entry, this function will not return unless there is an error
func watchEntry(entry *config.Entry) error {
	if entry == nil {
		return errors.New("No valid entry specified")
	}

	// Set the branch on the directory
	cmdBranch := exec.Command("git", "checkout", entry.Branch)
	cmdBranch.Dir = entry.Path
	out, branchErr := cmdBranch.Output()
	if branchErr != nil {
		return errors.New("Failed to select branch in specified directory")
	}
	log.Printf("git checkout %s: %s", entry.Branch, out)

	// Create an infinite loop that performs a pull at the specified interval
	// the loop breaks if there is an error
	var pullErr error
	for {
		pullErr = pull(entry)
		if pullErr != nil {
			break
		}
		time.Sleep(time.Minute * time.Duration(entry.Interval))
	}

	return pullErr
}

// pull will attempt to perform a git pull on the specified folder
// If changes are detected we execute the actions provided
func pull(entry *config.Entry) error {
	if entry == nil {
		return errors.New("No valid entry specified")
	}

	// Execute the pull command
	cmdPull := exec.Command("git", "pull")
	cmdPull.Dir = entry.Path
	out, pullErr := cmdPull.Output()
	if pullErr != nil {
		return errors.New("Failed to perform pull on folder")
	}
	log.Printf("%s", out)

	// Lets check to see if an update has been detected
	match, matchErr := regexp.Match(`^Updating\s`, out)
	if matchErr != nil {
		log.Printf("Failed to check regular expression")
	} else if match == true {
		log.Printf("Update detected")
		// Update has been detected, loop through each action command in sequence
		for _, action := range entry.Action {
			cmdAction := exec.Command(action.Command, action.Args...)
			cmdAction.Dir = entry.Path
			actionErr := cmdAction.Run()
			if actionErr != nil {
				log.Printf("Action Error: %s", actionErr)
				break
			}
		}
	}

	return nil
}
