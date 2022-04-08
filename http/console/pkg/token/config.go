package token

import (
	"github.com/filedag-project/filedag-storage/http/console/env"
	"github.com/filedag-project/filedag-storage/http/console/pkg/utils"
	"time"
)

// GetConsoleSTSDuration returns the default session duration for the STS requested tokens (defaults to 1h)
func GetConsoleSTSDuration() time.Duration {
	durationSeconds := env.Get(ConsoleSTSDurationSeconds, "")
	if durationSeconds != "" {
		duration, err := time.ParseDuration(durationSeconds + "s")
		if err != nil {
			duration = 1 * time.Hour
		}
		return duration
	}
	duration, err := time.ParseDuration(env.Get(ConsoleSTSDuration, "1h"))
	if err != nil {
		duration = 1 * time.Hour
	}
	return duration
}

var defaultPBKDFPassphrase = utils.RandomCharString(64)

// GetPBKDFPassphrase returns passphrase for the pbkdf2 function used to encrypt JWT payload
func GetPBKDFPassphrase() string {
	return env.Get(ConsolePBKDFPassphrase, defaultPBKDFPassphrase)
}

var defaultPBKDFSalt = utils.RandomCharString(64)

// GetPBKDFSalt returns salt for the pbkdf2 function used to encrypt JWT payload
func GetPBKDFSalt() string {
	return env.Get(ConsolePBKDFSalt, defaultPBKDFSalt)
}
