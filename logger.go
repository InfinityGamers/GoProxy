package mcpeproxy

import (
	"time"
	"fmt"
)

const (
	AnsiPre   = "\u001b["
	AnsiReset = AnsiPre + "0m"

	AnsiBold       = AnsiPre + "1m"
	AnsiItalic     = AnsiPre + "3m"
	AnsiUnderlined = AnsiPre + "4m"

	AnsiBlack   = AnsiPre + "30m"
	AnsiRed     = AnsiPre + "31m"
	AnsiGreen   = AnsiPre + "32m"
	AnsiYellow  = AnsiPre + "33m"
	AnsiBlue    = AnsiPre + "34m"
	AnsiMagenta = AnsiPre + "35m"
	AnsiCyan    = AnsiPre + "36m"
	AnsiWhite   = AnsiPre + "37m"
	AnsiGray    = AnsiPre + "30;1m"

	AnsiBrightRed     = AnsiPre + "31;1m"
	AnsiBrightGreen   = AnsiPre + "32;1m"
	AnsiBrightYellow  = AnsiPre + "33;1m"
	AnsiBrightBlue    = AnsiPre + "34;1m"
	AnsiBrightMagenta = AnsiPre + "35;1m"
	AnsiBrightCyan    = AnsiPre + "36;1m"
	AnsiBrightWhite   = AnsiPre + "37;1m"
)

const (
	Pre = "§"

	Black      = Pre + "0"
	Blue       = Pre + "1"
	Green      = Pre + "2"
	Cyan       = Pre + "3"
	Red        = Pre + "4"
	Magenta    = Pre + "5"
	Orange     = Pre + "6"
	BrightGray = Pre + "7"
	Gray       = Pre + "8"
	BrightBlue = Pre + "9"

	BrightGreen   = Pre + "a"
	BrightCyan    = Pre + "b"
	BrightRed     = Pre + "c"
	BrightMagenta = Pre + "d"
	Yellow        = Pre + "e"
	White         = Pre + "f"

	Obfuscated    = Pre + "k"
	Bold          = Pre + "l"
	StrikeThrough = Pre + "m"
	Underlined    = Pre + "n"
	Italic        = Pre + "o"

	Reset = Pre + "r"
)

const (
	InfoPrefix = "INFO"
	AlertPrefix = "ALERT"
	NoticePrefix = "NOTICE"
	DebugPrefix = "DEBUG"
	PanicPrefix = "PANIC"
)

// returns a log prefix with custom prefix
func LogFormat(prefix string) string {
	now := time.Now()
	return "[" + now.Format("2006-01-02 15:04:05") + "][Log/" + prefix + "]: "
}

func Info(log ...interface{})  {
	fmt.Print(AnsiYellow, LogFormat(InfoPrefix), AnsiWhite)
	fmt.Println(log...)
}

func Alert(log ...interface{})  {
	fmt.Print(AnsiRed, LogFormat(AlertPrefix), AnsiYellow)
	fmt.Println(log...)
}

func Notice(log ...interface{})  {
	fmt.Print(AnsiCyan, LogFormat(NoticePrefix), AnsiWhite)
	fmt.Println(log...)
}

func Debug(log ...interface{})  {
	fmt.Print(AnsiWhite, LogFormat(DebugPrefix), AnsiWhite)
	fmt.Println(log...)
}

func Panic(log ...interface{})  {
	fmt.Print(AnsiBold, AnsiRed, LogFormat(PanicPrefix))
	fmt.Println(log...)
}