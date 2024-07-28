package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	t.Run("check level by string", func(t *testing.T) {
		results := []struct {
			level    string
			excepted zapcore.Level
		}{
			{"fatal", zapcore.FatalLevel},
			{"error", zapcore.ErrorLevel},
			{"warning", zapcore.WarnLevel},
			{"debug", zapcore.DebugLevel},
			{"info", zapcore.InfoLevel},
			{"FaTAL", zapcore.FatalLevel},
			{"", zapcore.InfoLevel},
		}

		for _, v := range results {
			require.Equal(t, v.excepted, levelByString(v.level))
		}
	})
}
