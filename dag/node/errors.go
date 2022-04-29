package node

import "golang.org/x/xerrors"

var (
	ErrNone                = xerrors.New("mutcask: error none")
	ErrValueFormat         = xerrors.New("mutcask: invalid value format")
	ErrDataRotted          = xerrors.New("mutcask: data may be rotted")
	ErrKeySizeTooLong      = xerrors.New("mutcask: key size is too long")
	ErrHintFormat          = xerrors.New("mutcask: invalid hint format")
	ErrPathUndefined       = xerrors.New("mutcask: should define path within config")
	ErrPath                = xerrors.New("mutcask: path should be directory not file")
	ErrHintLogBroken       = xerrors.New("mutcask: hint log broken")
	ErrReadHintBeyondRange = xerrors.New("mutcask: read hint out of file range")
	ErrRepoLocked          = xerrors.New("mutcask: repo has been locked")
)
