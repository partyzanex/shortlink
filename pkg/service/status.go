package service

type Status uint8

const (
	StatusInitialized Status = iota + 1
	StatusStarted
	StatusFinished

	enumStatusInitialized = "initialized"
	enumStatusStarted     = "started"
	enumStatusFinished    = "finished"
	enumStatusUnknown     = "unknown"
)

func (s Status) String() string {
	switch s {
	case StatusInitialized:
		return enumStatusInitialized
	case StatusStarted:
		return enumStatusStarted
	case StatusFinished:
		return enumStatusFinished
	default:
		return enumStatusUnknown
	}
}
