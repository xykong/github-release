package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Fields map[string]interface{}

func Error(format string, a ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, color.RedString(format+"\n", a...))
}

func Essential(format string, a ...interface{}) {
	_, _ = fmt.Fprint(os.Stdout, color.CyanString(format+"\n", a...))
}

func Info(format string, a ...interface{}) {
	if viper.GetBool("essential") {
		return
	}
	_, _ = fmt.Fprint(os.Stdout, color.GreenString(format+"\n", a...))
}

func Infof(fields Fields, format string, a ...interface{}) {
	if viper.GetBool("essential") {
		return
	}

	var result []string
	for k, v := range fields {
		result = append(result, fmt.Sprintf("%s: %v", k, v))
	}

	value := fmt.Sprintf(format, a...)
	value = fmt.Sprintf("%s %s\n", value, strings.Join(result, " "))

	_, _ = fmt.Fprint(os.Stdout, color.GreenString(value))
}

func Verbose(format string, a ...interface{}) {
	if viper.GetBool("essential") {
		return
	}

	if viper.GetBool("verbose") {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func Debug(format string, a ...interface{}) {
	if viper.GetBool("essential") {
		return
	}

	if viper.GetBool("debug") {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}
