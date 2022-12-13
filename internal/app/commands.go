package app

type SetTimerCommand struct {
	Hours   int
	Minutes int
	Seconds int
	URLRaw  string
}

type GetTimer struct {
	ID string
}
