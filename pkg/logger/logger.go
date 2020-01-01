package logger // import github.com/mkaz/hastie/pkg/logger

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

type Logger struct {
	DebugLevel bool
	Verbose    bool
}

func (l Logger) Debug(a ...interface{}) {
	if l.DebugLevel {
		fmt.Println(chalk.Cyan, a, chalk.Reset)
	}
}

func (l Logger) Info(a ...interface{}) {
	if l.Verbose || l.DebugLevel {
		fmt.Println(chalk.Green, a, chalk.Reset)
	}
}

func (l Logger) Warn(a ...interface{}) {
	fmt.Println(chalk.Yellow, a, chalk.Reset)
}

func (l Logger) Fatal(a ...interface{}) {
	fmt.Println(chalk.Red, a, chalk.Reset)
	os.Exit(1)
}
