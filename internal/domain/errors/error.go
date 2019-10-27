package errors

type EventError string

func (ee EventError) Error() string {
	return string(ee)
}

var (
	ErrNotFound         = EventError("event not found")
	ErrOverlaping       = EventError("another event exists for this date")
	ErrIncorrectEndDate = EventError("end-date is incorrect")
)
