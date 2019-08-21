package util

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"
)

// Tracing for debugging
func Trace1() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}

// Tracing for debugging
func Trace2() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
}

// Tracing for debugging
func Debug(format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	info := fmt.Sprintf(format, a...)

	log.Printf("%s:%d %v", file, line, info)
}

// string of function name
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// ex. fmt.Printf("%s\n", whereami.WhereAmI())
func WhereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("File: %s  Function: %s Line: %d", file, runtime.FuncForPC(function).Name(), line)
}

func DebugPrintf(fmt_ string, args ...interface{}) {
	programCounter, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	prefix := fmt.Sprintf("[%s,%d: %s] %s", path.Base(file), line, fn.Name(), fmt_)
	fmt.Printf(prefix, args...)
	fmt.Println()
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		host := r.Host
		// Slice off any port information.
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		return host
	}
	return r.URL.Host
}
