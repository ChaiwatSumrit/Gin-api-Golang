package mgorus

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
)

func main() {
	log := logrus.New()

	// connect to database (ipdatabase, databasename, collection)
	hooker, err := mgorus.NewHooker("mongodb+srv://development:ZGFMzUvDJ745GFDq@clustermaster-zvis2.mongodb.net/test", "logging", "logrus")
	if err == nil {
		log.Hooks.Add(hooker)
	} else {
		fmt.Print(err)
	}

	// set format
	log.SetFormatter(&logrus.JSONFormatter{})

	// output log
	log.WithFields(logrus.Fields{
		"event": "test",
	}).Warning("this is error messages")
}
