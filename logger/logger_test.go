package logger

import (
	"fmt"
	murlog "github.com/Melenium2/Murlog"
	"github.com/RecleverLogger/logger/externallogger"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func config() *murlog.Config {
	c := murlog.NewConfig()
	c.TimePref(time.ANSIC)
	c.Pref(func() interface{} {
		return "INFO"
	})
	c.CallerCustomPref(3)

	return c
}

func TestSuperLogger_Logf_ShouldReturnRightString(t *testing.T) {
	l := New(murlog.NewLogger(config()), nil)
	l.Logf("Custom message %s", "123")
}

func TestSuperLogger_Logs_ShouldReturnRightStringTo(t *testing.T) {
	token := os.Getenv("token")
	if token == "" {
		t.Fatal("For this test you need to provide tg token env var")
	}
	chatId := os.Getenv("chatId")
	if chatId == "" {
		t.Fatal("For this test you need to provide tg chatId env var")
	}
	id, err := strconv.Atoi(chatId)
	assert.NoError(t, err)

	e := externallogger.NewTelegramLogger(token, id)
	l := New(murlog.NewLogger(config()), e)

	l.Logs(fmt.Sprintf("Custom message %s", "123"))
}