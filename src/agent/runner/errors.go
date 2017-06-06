package runner

type TimeoutError struct{}

func (t *TimeoutError) Error() string {
	return "timeout error"
}
