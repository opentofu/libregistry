package module

type InvalidModuleAddrError struct {
	Addr  Addr
	Cause error
}

func (i InvalidModuleAddrError) Error() string {
	if i.Cause != nil {
		return "Invalid module address: " + i.Addr.String() + " (" + i.Cause.Error() + ")"
	}
	return "Invalid module address: " + i.Addr.String()
}

func (i InvalidModuleAddrError) Unwrap() error {
	return i.Cause
}
