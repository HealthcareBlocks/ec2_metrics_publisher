package logger

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerReturnsLogEntry(t *testing.T) {
	assert.IsType(t, log.Entry{}, *NewLogger())
}

func TestNewLoggerSetsContextField(t *testing.T) {
	assert.Equal(t, "ec2_metrics_publisher", NewLogger().Data["context"])
}

func TestNewLoggerSetsJSONFormatter(t *testing.T) {
	assert.Equal(t, &log.JSONFormatter{}, NewLogger().Logger.Formatter)
}
