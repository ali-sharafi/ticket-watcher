package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/erfanmomeniii/jalali"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func init() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./logs/" + time.Now().Format("2006-01-02") + "_app.log"),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = time.RFC3339
	logFormatter.FullTimestamp = true

	Logger = logrus.New()
	Logger.SetFormatter(logFormatter)
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(multiWriter)
}

func GetJalaliDate(gregorianDate string) string {
	t, err := time.Parse("2006-01-02", gregorianDate)
	if err != nil {
		Logger.Error("Error:", err)
		return ""
	}

	// Convert to Jalali
	j := jalali.ConvertGregorianToJalali(t)

	jalaliDate := fmt.Sprintf(`%d-%d-%d`, j.Year(), j.Month(), j.Day())

	return jalaliDate
}
