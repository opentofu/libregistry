# OCI client for OpenTofu

> [!WARNING]
> This file is not an end-user documentation, it is intended for developers. Please follow the user documentation on the OpenTofu website unless you want to work on the OCI implementation.

This package contains the OCI registry client code for OpenTofu. It consists of two client implementations, the `OCIClient` (high level) and the `RawOCIClient` (low level) clients. For almost all use cases you should use the `OCIClient` (high level) interface.
