package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"runtime"
	"time"
)

// 日志级别
type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	// 在标准库的 log 模块基础上，进行一定的扩展
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{
		newLogger: l,
	}
}

// 复制一个新的
func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// 设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}

	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

// 设置上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	l.ctx = ctx
	return ll
}

// 使用 opentracing 时用，在 fields 中加入 spanId 和 traceId 字段
func (l *Logger) WithTrace() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	// 判断一下是不是 gin.Context，如果不是的话，不一定有 gin 框架提供的方法， 跳过这一步，直接原样返回
	if ok {

		fields := Fields{}

		if traceID, ok := ginCtx.Get("X-Trace-ID"); ok {
			fields["X-Trace-ID"] = traceID
		}

		if spanID, ok := ginCtx.Get("X-Span-ID"); ok {
			fields["X-Span-ID"] = spanID
		}

		if len(fields) > 0 {
			return l.WithFields(fields)
		}
	}
	return l
}

func (l *Logger) WithReqID() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	if ok {
		fields := Fields{}
		if reqID, ok := ginCtx.Get("x-request-id"); ok {
			fields["x-request-id"] = reqID
		}
		if len(fields) > 0 {
			return l.WithFields(fields)
		}
	}
	return l
}

// 设置当前某一层调用栈信息（程序计数器、文件信息和行号）
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	// runtime.Caller 可以获取运行时方法的调用信息
	// skip 表示跳过的栈帧数，0 表示不跳过，也就是 runtime.Caller 的调用者，1 的话就是向上一层，表示调用者的调用者
	// pc 是函数指针，file函数所在文件名路径，line 行号，ok 是否可以获取到当前信息
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s:%d %s", file, line, f.Name())}
	}

	return ll
}

// 设置当前的整个调用栈信息，一般打日志不用这么详细
func (l *Logger) WithCallerFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	var callers []string
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			// 没有更多了，中断循环
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

// 在输出之前，将 Logger 中保存的信息，格式化成一个 json 对象
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
	// 初始化一个 map，在 l.fields 的基础上加四个字段
	data := make(Fields, len(l.fields)+4)

	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data[" message"] = message
	data["callers"] = l.callers

	if len(l.fields) > 0 {
		// 遍历 l.fields ，将非上述四个字段的值都加上去
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Logger) Output(level Level, message string) {
	ll := l.WithTrace().WithReqID().WithCaller(3)
	body, _ := json.Marshal(ll.JSONFormat(level, message))
	content := string(body)

	switch level {
	case LevelDebug:
		fallthrough
	case LevelInfo:
		fallthrough
	case LevelWarn:
		fallthrough
	case LevelError:
		ll.newLogger.Println(content)
	case LevelFatal:
		ll.newLogger.Fatal(content)
	case LevelPanic:
		ll.newLogger.Panic(content)
	}
}

func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelWarn, fmt.Sprint(v...))
}

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelFatal, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).Output(LevelPanic, fmt.Sprint(v...))
}

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).Output(LevelPanic, fmt.Sprintf(format, v...))
}
