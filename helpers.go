package dashamail

// Bool returns a pointer to the given bool value.
// Useful for setting optional boolean fields on Message.
//
//	msg.NoTrackOpens = dashamail.Bool(false)
func Bool(v bool) *bool {
	return &v
}
