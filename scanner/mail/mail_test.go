package mail_test

import (
	"testing"

	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/logger"
	"github.com/back2basic/siaalert/scanner/mail"
	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	// Mock the configuration
	config.LoadConfig("../config.yaml")

	log := logger.GetLogger()
	defer logger.Sync()
	to := "test@example.com"
	host := "example-host"
	status := "down"

	// Call the SendMail function
	assert.NotPanics(t, func() {
		mail.SendAlert(to, host, status, log)
	}, "The SendMail function should not panic")

	// Add additional assertions if applicable (e.g., check the creation of email content or mock function calls)
	t.Log("SendMail tested successfully with no panics.")
}
