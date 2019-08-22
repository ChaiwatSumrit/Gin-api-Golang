package main

import (
	syslog "github.com/racksec/srslog"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/shinji62/logrus-syslog-ng"
  )
  
  func main() {
	log       := logrus.New()
	hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:8080", syslog.LOG_INFO, "")
  
	if err == nil {
	  log.Hooks.Add(hook)
	}
  }