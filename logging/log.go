package logging

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const LogPath = "log"

func Setup() error {
	logf, err := rotatelogs.New(
		path.Join(LogPath, "%Y%m%d"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logf))
	return nil
}