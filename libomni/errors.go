package libomni

type MissingProviderError struct {
	err error
}

func (e MissingProviderError) Error() string {
	return e.err.Error()
}

type ParseError struct {
	err error
}

func (e ParseError) Error() string {
	return e.err.Error()
}
