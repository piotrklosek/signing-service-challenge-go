package handlers

import "errors"

var (
	// general errors
	ErrInvalidJson = errors.New("invalid json")
	ErrBadRequest  = errors.New("bad request")

	// signature errors

	// device errors
	ErrDeviceNotFounc = errors.New("device not founc")
	ErrGenerateKeys   = errors.New("error generating keys")
	ErrCreatingDevice = errors.New("error creating device")
	ErrListDevices    = errors.New("error list devices")

	// user errors
	ErrUserNotFounc = errors.New("user not found")
)
