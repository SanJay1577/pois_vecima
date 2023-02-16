package schedule

func noop(channelName string, date string, time string) (string, string) {
	return "", ""
}

// init to register provider with its implememtaion logic signature
func init() {
	registerProvider("noop", noop)
}
