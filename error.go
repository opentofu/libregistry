package libregistry

import (
	"github.com/opentofu/libregistry/types/module"
)

type ModuleAlreadyExistsError struct {
	Module module.Addr
}

func (m ModuleAlreadyExistsError) Error() string {
	return "Module already exists: " + m.Module.String()
}

type ModuleAddFailedError struct {
	Module module.Addr
	Cause  error
}

func (m ModuleAddFailedError) Error() string {
	return "Adding the module " + m.Module.String() + " failed: " + m.Cause.Error()
}

func (m ModuleAddFailedError) Unwrap() error {
	return m.Cause
}

type ModuleUpdateFailedError struct {
	Module module.Addr
	Cause  error
}

func (m ModuleUpdateFailedError) Error() string {
	return "Updating the module " + m.Module.String() + " failed: " + m.Cause.Error()
}

func (m ModuleUpdateFailedError) Unwrap() error {
	return m.Cause
}
