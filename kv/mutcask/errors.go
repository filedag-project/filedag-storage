package mutcask

import "golang.org/x/xerrors"

var (
	ErrValueFormat    = xerrors.New("mutcask: invalid value format")
	ErrDataRotted     = xerrors.New("mutcask: data may be rotted")
	ErrKeySizeTooLong = xerrors.New("mutcask: key size is too long")
	ErrHintFormat     = xerrors.New("mutcask: invalid hint format")
)
