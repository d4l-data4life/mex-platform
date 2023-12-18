package grpc_adapt

import (
	"context"
	"fmt"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

// Adapter for the grpc.LoggerV2 interface
type GRPCLogger struct {
	L.Logger
	Level int
}

// Taken from gRPC's loggerv2.go:
const (
	GRPCLogLevelInfo int = iota
	GRPCLogLevelWarn
	GRPCLogLevelError
	GRPCLogLevelFatal
	GRPCLogLevelOff
)

var levelNames = map[string]int{
	"info":    GRPCLogLevelInfo,
	"warn":    GRPCLogLevelWarn,
	"warning": GRPCLogLevelWarn,
	"error":   GRPCLogLevelError,
	"fatal":   GRPCLogLevelFatal,
	"off":     GRPCLogLevelOff,
	"silent":  GRPCLogLevelOff,
}

func ParseLevel(levelString string) int {
	if level, ok := levelNames[levelString]; ok {
		return level
	}

	// Use info as default level
	return GRPCLogLevelInfo
}

func (l GRPCLogger) Info(args ...interface{}) {
	if l.Level <= GRPCLogLevelInfo {
		l.Logger.Info(context.TODO(), L.Message(fmt.Sprint(args...)))
	}
}

func (l GRPCLogger) Infoln(args ...interface{}) {
	if l.Level <= GRPCLogLevelInfo {
		l.Logger.Info(context.TODO(), L.Message(fmt.Sprintln(args...)))
	}
}

func (l GRPCLogger) Infof(format string, args ...interface{}) {
	if l.Level <= GRPCLogLevelInfo {
		l.Logger.Info(context.TODO(), L.Messagef(format, args...))
	}
}

func (l GRPCLogger) Warning(args ...interface{}) {
	if l.Level <= GRPCLogLevelWarn {
		l.Logger.Warn(context.TODO(), L.Message(fmt.Sprint(args...)))
	}
}

func (l GRPCLogger) Warningln(args ...interface{}) {
	if l.Level <= GRPCLogLevelWarn {
		l.Logger.Warn(context.TODO(), L.Message(fmt.Sprintln(args...)))
	}
}

func (l GRPCLogger) Warningf(format string, args ...interface{}) {
	if l.Level <= GRPCLogLevelWarn {
		l.Logger.Warn(context.TODO(), L.Messagef(format, args...))
	}
}

func (l GRPCLogger) Error(args ...interface{}) {
	if l.Level <= GRPCLogLevelError {
		l.Logger.Error(context.TODO(), L.Message(fmt.Sprint(args...)))
	}
}

func (l GRPCLogger) Errorln(args ...interface{}) {
	if l.Level <= GRPCLogLevelError {
		l.Logger.Error(context.TODO(), L.Message(fmt.Sprintln(args...)))
	}
}

func (l GRPCLogger) Errorf(format string, args ...interface{}) {
	if l.Level <= GRPCLogLevelError {
		l.Logger.Error(context.TODO(), L.Messagef(format, args...))
	}
}

func (l GRPCLogger) Fatal(args ...interface{}) {
	if l.Level <= GRPCLogLevelFatal {
		s := fmt.Sprint(args...)
		l.Logger.Error(context.TODO(), L.Message(s), L.Reason("GRPC FATAL"))
		panic(s)
	}
}

func (l GRPCLogger) Fatalln(args ...interface{}) {
	if l.Level <= GRPCLogLevelFatal {
		s := fmt.Sprintln(args...)
		l.Logger.Error(context.TODO(), L.Message(s), L.Reason("GRPC FATAL"))
		panic(s)
	}
}

func (l GRPCLogger) Fatalf(format string, args ...interface{}) {
	if l.Level <= GRPCLogLevelFatal {
		s := fmt.Sprintf(format, args...)
		l.Logger.Error(context.TODO(), L.Message(s), L.Reason("GRPC FATAL"))
		panic(s)
	}
}

func (l GRPCLogger) V(level int) bool {
	return level <= l.Level
}
