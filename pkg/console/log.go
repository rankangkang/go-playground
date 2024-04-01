package console

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Log(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func Warning(format string, args ...interface{}) {
	fields := []interface{}{IconWarning}
	fields = append(fields, fmt.Sprintf(format, args...))
	fmt.Println(fields...)
}

func Error(format string, args ...interface{}) {
	fields := []interface{}{color.RedString(IconFailure)}
	fields = append(fields, fmt.Sprintf(format, args...))
	fmt.Println(fields...)
}

func Success(format string, args ...interface{}) {
	fields := []interface{}{color.GreenString(IconSuccess)}
	fields = append(fields, fmt.Sprintf(format, args...))
	fmt.Println(fields...)
}

func Info(format string, args ...interface{}) {
	fields := []interface{}{(color.BlueString(IconInfo))}
	fields = append(fields, fmt.Sprintf(format, args...))
	fmt.Println(fields...)
}

func Fatal(err error) {
	Error(err.Error())
	os.Exit(1)
}
