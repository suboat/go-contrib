package contrib

// Logger 日志输出
type Logger interface {
	//
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	//
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	//
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	//
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	//
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	//
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	//
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	//
	SetLevel(level int)
	GetLevel() int
	Output(calldepth int, s string) error
}
