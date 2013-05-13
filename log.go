package eyego

import "io"
import "log"

type LogLevel struct {
	level int
	name string
}

var (
	TRACE = LogLevel{level:0, name:"TRACE"}
	DEBUG = LogLevel{level:1, name:"DEBUG"}
	INFO = LogLevel{level:2, name:"INFO"}
	WARN = LogLevel{level:3, name:"WARN"}
	ERROR = LogLevel{level:4, name:"ERROR"}

	levels = []LogLevel{TRACE,DEBUG,INFO,WARN,ERROR}
	level int
)

func LogLevelFromName(name string) LogLevel{
	for i := range levels {
		if levels[i].name == name {
			return levels[i]
		}
	}

	panic("No such level")
}

func Init(w io.Writer, _level LogLevel) {
	log.SetOutput(w)
	level = _level.level
}

func Trace(s string, args ...interface {}) {
	if level <= TRACE.level {
		log.Printf(s, args...)
	}
}

func Debug(s string, args ...interface {}) {
	if level <= DEBUG.level {
		log.Printf(s, args...)
	}
}

func Info(s string, args ...interface {}) {
	if level <= INFO.level {
		log.Printf(s, args...)
	}
}

func Warn(s string, args ...interface {}) {
	if level <= WARN.level {
		log.Printf(s, args...)
	}
}

func LogError(s string, args ...interface {}) {
	if level <= ERROR.level {
		log.Printf(s, args...)
	}
}
