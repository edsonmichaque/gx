package dispatcher

type NoProviderError struct {
	err error
}

func (e NoProviderError) Error() string {
	return e.err.Error()
}

type ConnectionClosedError struct {
	err error
}

func (e ConnectionClosedError) Error() string {
	return e.err.Error()
}

type ParseError struct {
	err error
}

func (e ParseError) Error() string {
	return e.err.Error()
}
