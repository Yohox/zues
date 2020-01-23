package logger

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Formatter struct {
}




const (
	color_red = uint8(iota + 91)
	color_green		//	绿
	color_yellow		//	黄
	color_blue			// 	蓝
	color_magenta 		//	洋红
)

const (
	fatalPrefix		=	"[FATAL] "
	errorPrefix		=	"[ERROR] "
	warnPrefix		=	"[WARN] "
	infoPrefix		=	"[INFO] "
	debugPrefix		=	"[DEBUG] "
)


func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}

func getPrefix(entry *log.Entry) string {
	switch entry.Level {
	case log.DebugLevel:
		return blue("[DEBUG]")
	case log.InfoLevel:
		return green("[INFO]")
	case log.ErrorLevel:
		return red("[ERROR]")
	case log.FatalLevel:
		return red("[FATAL]")
	case log.PanicLevel:
		return red("[PANIC]")
	case log.TraceLevel:
		return magenta("[TRACE]")
	case log.WarnLevel:
		return yellow("[WARN]")
	default:
		panic("日志级别有误！")
	}
}

func getDateTime(entry *log.Entry) string {
	time := entry.Time
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%d", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), time.Nanosecond())
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", getDateTime(entry), strings.ToUpper(entry.Level.String()), entry.Message))
	fmt.Println(fmt.Sprintf("%s %s %s", getDateTime(entry), getPrefix(entry), entry.Message))
	return b.Bytes(), nil
}