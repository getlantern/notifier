package notify

func newNotifier() (Notifier, error) {
	return &emptyNotifier{}, nil
}
