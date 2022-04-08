package restapi

import (
	"log"
	"os"
)

var infoLog = log.New(os.Stdout, "I: ", log.LstdFlags)
var errorLog = log.New(os.Stdout, "E: ", log.LstdFlags)

func logInfo(msg string, data ...interface{}) {
	infoLog.Printf(msg+"\n", data...)
}

func logError(msg string, data ...interface{}) {
	errorLog.Printf(msg+"\n", data...)
}

// globally changeable logger styles
var (
	LogInfo  = logInfo
	LogError = logError
)

// Context captures all command line flags values
type Context struct {
	Host                string
	HTTPPort, HTTPSPort int
	TLSRedirect         string
	// Legacy options, TODO: remove in future
	TLSCertificate, TLSKey, TLSca string
}
