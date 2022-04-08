package certs

import "github.com/rjeczalik/notify"

// eventWrite contains the notify events that will cause a write
var eventWrite = []notify.Event{notify.Create, notify.Write}
