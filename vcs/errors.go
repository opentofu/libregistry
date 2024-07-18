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
