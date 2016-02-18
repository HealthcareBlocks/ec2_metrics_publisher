// Package logger configures default logging behavior
package logger

import log "github.com/Sirupsen/logrus"

// NewLogger sets up default behavior for the logging object
func NewLogger() *log.Entry {
	// Log in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	l := log.WithFields(log.Fields{
		"context": "ec2_metrics_publisher",
	})

	return l
}
