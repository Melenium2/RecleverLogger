package externallogger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTelegramLogger_Sendlog_ShouldSaveLogToTelegramBot(t *testing.T) {
	bot := NewTelegramLogger("1293039613:AAER81Qqklo9JZQa3kt2iHrKBA9ptPpJ8IY", 708015155)
	err := bot.Sendlog(0, "Just for test")
	assert.NoError(t, err)
}
