package token

const (
	ConsoleSTSDurationSeconds = "CONSOLE_STS_DURATION_SECONDS" // (deprecated), set value in seconds for sts session, ie: 3600
	ConsoleSTSDuration        = "CONSOLE_STS_DURATION"         // time.Duration format, ie: 3600s, 2h45m, 1h, etc
	ConsolePBKDFPassphrase    = "CONSOLE_PBKDF_PASSPHRASE"
	ConsolePBKDFSalt          = "CONSOLE_PBKDF_SALT"
)
