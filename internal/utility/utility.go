package utility

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

func ParseLogLevelFromString(levelStr string) (slog.Level, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "info", "":
		return slog.LevelInfo, nil
	default:
		return slog.LevelInfo, fmt.Errorf("Couldn't parse '%s' as a valid log level, defaulting to 'info'", levelStr)
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func HandleErr(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}

func CheckFileAccess(path string) error {
	// doesnt't exist -> not read/writeable
	if _, err := os.Stat(path); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// log out pointers for debug purposes.
func LogPtr[T any](ptr *T) {
	slog.Debug(fmt.Sprintf("Pointer debugging: %T %v %p %v", ptr, &ptr, ptr, *ptr))
}
