package utils

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func PrintSuccess(formatstr string, a ...interface{}) {
	green := color.New(color.FgGreen).Add(color.Bold)
	green.Print("[*] ")
	fmt.Printf(formatstr, a...)
}

func PrintStatus(formatstr string, a ...interface{}) {
	blue := color.New(color.FgBlue).Add(color.Bold)
	blue.Print("[*] ")
	fmt.Printf(formatstr, a...)
}

func PrintWarning(formatstr string, a ...interface{}) {
	yellow := color.New(color.FgYellow).Add(color.Bold)
	yellow.Print("[*] ")
	fmt.Printf(formatstr, a...)
}

func PrintErr(formatstr string, a ...interface{}) {
	red := color.New(color.FgRed).Add(color.Bold)
	white := color.New(color.FgWhite).Add(color.Bold)
	red.Print("[-] ")
	white.Printf(formatstr, a...)
}

func PrintGreen(formatstr string, a ...interface{}) {
	green := color.New(color.FgGreen).Add(color.Bold)
	green.Printf(formatstr, a...)
}

func PrintFatal(err error) {
	if err != nil {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			logrus.Fatalf("%s: %s\n", strings.ToUpper(strings.Split(details.Name(), ".")[1]), err)
		} else {
			logrus.Fatal(err)
		}
	}
}

// randomString - generates random string of given length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, length)
	for i := 0; i < length; i++ {
		random[i] = charset[rand.Intn(len(charset))]
	}
	return string(random)
}

// GetFileSize retrieves the size of the file with given file path
func GetFileSize(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return int(stat.Size()), nil
}
