package logger

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"reflect"
	"sync"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	panicLogger *log.Logger
	logFile     *os.File
	once        sync.Once
)

func init() {
	once.Do(func() {
		var err error
		logFile, err = os.OpenFile("river.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		multiWriter := io.MultiWriter(logFile, os.Stdout)

		infoLogger = log.New(multiWriter, color.HiGreenString("INFO: "), log.Ldate|log.Ltime|log.Lshortfile)
		errorLogger = log.New(multiWriter, color.HiRedString("ERROR: "), log.Ldate|log.Ltime|log.Lshortfile)
		panicLogger = log.New(multiWriter, color.HiRedString("PANIC: "), log.Ldate|log.Ltime|log.Lshortfile)
	})
}

func closeLogFile() {
	if err := logFile.Close(); err != nil {
		log.Fatalf("Failed to close log file: %v", err)
	}
}

func Info(msg ...any) {
	infoLogger.Output(2, formatMessage(msg...))
}

func Error(err error, msg ...interface{}) error {
	if err != nil {
		formattedMsg := fmt.Sprintf("[ERROR] %s: %v", fmt.Sprint(msg...), err)
		errorLogger.Output(2, formattedMsg)
	} else {
		formattedMsg := fmt.Sprint(msg...)
		errorLogger.Output(2, formattedMsg)
	}
	return fmt.Errorf(fmt.Sprint(msg...), err)
}

func Panic(err error, msg ...any) {
	closeLogFile()
	if err != nil {
		panicLogger.Output(2, formatMessage(append([]any{err}, msg...)...))
	} else {
		panicLogger.Output(2, formatMessage(msg...))
	}
	panic(fmt.Sprintln(formatMessage(msg...)))
}

func formatMessage(msg ...any) string {
	formattedMsg := ""
	for _, m := range msg {
		switch v := m.(type) {
		case string:
			formattedMsg += v + " "
		case fmt.Stringer:
			formattedMsg += v.String() + " "
		default:
			if jsonData, err := jsonMarshal(v); err == nil {
				formattedMsg += string(jsonData) + " "
			} else {
				formattedMsg += fmt.Sprint(v) + " "
			}
		}
	}
	return formattedMsg
}

func jsonMarshal(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct || rv.Kind() == reflect.Map {
		return json.MarshalIndent(v, "", "  ")
	}
	return nil, fmt.Errorf("not a struct or map")
}
