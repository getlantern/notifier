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

type darwinNotifier struct{}

// Notify sends a desktop notification
// Note: terminal-notifier will be used if it is installed; otherwise, fall back to osascript.
func (n *darwinNotifier) Notify(msg *Notification) error {
	if _, err := exec.LookPath(terminalNotifier); err == nil {
		return tnNotify(msg)
	}
	return osaNotify(msg)
}

// osaNotify sends a notification using AppleScript with `osascript` binary
func osaNotify(msg *Notification) error {
	log.Debug("Sending notification with osascript")
	osa, err := exec.LookPath(osascript)
	if err != nil {
		return err
	}

	script := fmt.Sprintf("display notification %q with title %q", msg.Message, msg.Title)
	cmd := exec.Command(osa, "-e", script)
	return cmd.Run()
}

// tnNotify sends a notification using terminal-notifier
func tnNotify(msg *Notification) error {
	log.Debug("Sending notification with terminal-notifier")
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
	cmd := exec.Command(terminalNotifier, args...)
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
