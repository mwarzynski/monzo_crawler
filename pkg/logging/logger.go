package logging

import "github.com/sirupsen/logrus"

type Logger = logrus.FieldLogger

const (
	FieldKeyPackage   = "package"
	FieldKeyComponent = "component"
)

func WithFields(l Logger, pack, component string) Logger {
	return l.WithField(FieldKeyComponent, component).
		WithField(FieldKeyPackage, pack)
}
