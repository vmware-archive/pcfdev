package provisioner

type timeoutError struct{}

func (t *timeoutError) Error() string {
	return "timeout error"
}
