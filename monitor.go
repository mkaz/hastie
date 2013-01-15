package hastie

import (
	"fmt"
	"time"
)

// Monitor is an interface whose methods are called during hastie.Compile.
type Monitor interface {
	Start()
	Walked()
	ParsingSource(file string)
	ParsedSources()
	Listed()
	ParsedTemplates()
	ParsingTemplate(file string)
	WritingTemplate(file string)
	GeneratedTemplates()
	Filtered()
	End()
}

type discardMonitor struct{}

// DiscardMonitor implements the Monitor interface and does nothing with the events.
var DiscardMonitor Monitor = discardMonitor{}

func (discardMonitor) Start() {
}

func (discardMonitor) Walked() {
}

func (discardMonitor) ParsingSource(file string) {
}

func (discardMonitor) ParsedSources() {
}

func (discardMonitor) Listed() {
}

func (discardMonitor) ParsedTemplates() {
}

func (discardMonitor) ParsingTemplate(file string) {
}

func (discardMonitor) WritingTemplate(file string) {
}

func (discardMonitor) GeneratedTemplates() {
}

func (discardMonitor) Filtered() {
}

func (discardMonitor) End() {
}

type logMonitor struct {
	log     func(msg string)
	timing  bool
	verbose bool
	start   time.Time
	last    time.Time
}

// NewLogMonitor creates a Monitor implementation which simply logs the event output to the provided log method.
func NewLogMonitor(log func(msg string), timing bool, verbose bool) Monitor {
	return &logMonitor{log: log, timing: timing, verbose: verbose}
}

func (l *logMonitor) Start() {
	now := time.Now()
	l.start = now
	l.last = now
	l.log("Starting build")
}

func (l *logMonitor) Walked() {
	l.elapsedTimer("File Walker")
}

func (l *logMonitor) ParsingSource(file string) {
	if l.verbose {
		l.log(fmt.Sprint(" File:", file))
	}
}

func (l *logMonitor) ParsedSources() {
	l.elapsedTimer("Loop and Parse")
}

func (l *logMonitor) Listed() {
	l.elapsedTimer("Recent and Category Lists")
}

func (l *logMonitor) ParsedTemplates() {
	l.elapsedTimer("Parsed Templates")
}

func (l *logMonitor) ParsingTemplate(file string) {
	if l.verbose {
		l.log(fmt.Sprint(" Generating Template:", file))
	}
}

func (l *logMonitor) WritingTemplate(file string) {
	if l.verbose {
		l.log(fmt.Sprint(" Writing File:", file))
	}
}

func (l *logMonitor) GeneratedTemplates() {
	l.elapsedTimer("Generate Templates")
}

func (l *logMonitor) Filtered() {
	l.elapsedTimer("Process Filters")
}

func (l *logMonitor) elapsedTimer(str string) {
	if !l.timing {
		return
	}
	l.log(fmt.Sprintf("Event: %-25s -- %9v  (%9v)", str, time.Since(l.last), time.Since(l.start)))
	l.last = time.Now()
}

func (l *logMonitor) End() {
	l.log("Finished build")
}
