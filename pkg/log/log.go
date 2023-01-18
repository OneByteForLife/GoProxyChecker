package log

import "github.com/sirupsen/logrus"

func ConfigLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
}
