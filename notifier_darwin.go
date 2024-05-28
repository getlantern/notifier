package notify

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
)

const (
	osascript        = "osascript"
	terminalNotifier = "terminal-notifier"
)

func newNotifier() (Notifier, error) {
	return &darwinNotifier{}, nil
}

type darwinNotifier struct {
	path string
}

// Notify sends a desktop notification
// if terminal-notifier exists, use it. Otherwise, fall back to osascript.
func (n *darwinNotifier) Notify(msg *Notification) error {
	if _, err := exec.LookPath(terminalNotifier); err == nil {
		return n.tnNotify(msg)
	}
	return n.osaNotify(msg)
}

// Notify sends a notification to the user using AppleScript with `osascript` binary
func (n *darwinNotifier) osaNotify(msg *Notification) error {
	osa, err := exec.LookPath(osascript)
	if err != nil {
		return err
	}

	script := fmt.Sprintf("display notification %q with title %q", msg.Message, msg.Title)
	cmd := exec.Command(osa, "-e", script)
	return cmd.Run()
}

func (n *darwinNotifier) tnNotify(msg *Notification) error {
	timeout := msg.AutoDismissAfter
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	args := []string{
		"-message", msg.Message,
		"-title", msg.Title,
		"-timeout", fmt.Sprintf("%d", int(timeout.Seconds())),
	}
	if msg.Sender != "" {
		args = append(args, "-sender", msg.Sender)
		// override IconURL
		msg.IconURL = ""
	}
	if msg.ClickURL != "" {
		label := msg.ClickLabel
		if label == "" {
			label = "Open"
		}
		args = append(args, "-actions", label)
	}
	if msg.IconURL != "" {
		args = append(args, "-appIcon", msg.IconURL)
	}
	log.Debugf("Running command %s %v", n.path, args)
	cmd := exec.Command("terminal-notifier", args...)
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Could not run command %w", err)
		return err
	}
	result := string(res)
	log.Debugf("Received result: %v", result)

	// Note we can't just look for the result being the string "Open" here because
	// it's the label on the button and can be in any language.
	if msg.ClickURL != "" && result != "" && !strings.Contains(result, "@") {
		open.Start(msg.ClickURL)
	}
	return nil
}
