package vcs

type InvalidRepositoryAddrError struct {
	RepositoryAddr string
	Cause          error
}

func (r InvalidRepositoryAddrError) Error() string {
	if r.Cause != nil {
		return "Failed to parse repository address: " + r.RepositoryAddr + " (" + r.Cause.Error() + ")"
	}
	return "Failed to parse repository address: " + r.RepositoryAddr
}

func (r InvalidRepositoryAddrError) Unwrap() error {
	return r.Cause
}

type RequestFailedError struct {
	Cause error
}

func (r RequestFailedError) Error() string {
	return "VCS request failed: " + r.Cause.Error()
}

func (r RequestFailedError) Unwrap() error {
	return r.Cause
}
