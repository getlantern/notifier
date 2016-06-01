package notify

import (
	"fmt"
	"os/exec"
)

type windowsNotifier struct {
	path string
}

// Notify sends a notification to the user.
func (n *windowsNotifier) Notify(msg *Notification) error {
	cmd := exec.Command(n.path, "/m", msg.Message, "/p", msg.Title)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not run command %v", err)
	}
	return nil
}
