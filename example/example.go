package main

import (
	"time"

	"github.com/getlantern/notifier"
)

func main() {
	n := notify.NewNotifications()

	msg := &notify.Notification{
		Title:    "Your Lantern time is up",
		Message:  "You have reached your data cap limit",
		ClickURL: "https://www.getlantern.org",
		//IconURL:  "https://www.getlantern.org",
	}

	n.Notify(msg)
	time.Sleep(3 * time.Second)
}
