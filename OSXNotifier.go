package notify

import (
	"fmt"
	"os/exec"
)

type osxNotifier struct {
	path string
}

// Notify sends a notification to the user.
func (n *osxNotifier) Notify(msg *Notification) error {
	cmd := exec.Command(n.path, "-message", msg.Message)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not run command %v", err)
	}
	return nil
}
