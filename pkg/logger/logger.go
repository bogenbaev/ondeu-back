package logger

import "github.com/sirupsen/logrus"

func GetLogger(channel string) *logrus.Entry {
	return logrus.WithFields(
		logrus.Fields{
			"channel": channel,
		},
	)
}
