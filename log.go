package logger

import (
	"strings"
	"time"

	"github.com/cihub/seelog"
)

const (
	FLUSH_INTERVAL = 2 //定时刷新的时间间隔
)

var logger seelog.LoggerInterface
var isClose bool = false

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func Criticalf(format string, args ...interface{}) {
	logger.Criticalf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Warning(args ...interface{}) {
	logger.Warn(args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Critical(args ...interface{}) {
	logger.Critical(args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

type LogConfig struct {
	Root     string
	Type     string
	Maxdays  string
	Maxfiles string
	Maxsize  string
	Level    string
}
// debug,info,warn,error,critical
func Init(cfg LogConfig) {

	var useConfig string
	if strings.EqualFold(cfg.Type, "date") {
		useConfig = `
<seelog minlevel="debug" maxlevel="critical">
    <outputs formatid="main">
        <filter levels="critical"> 
            <console formatid="colored-critical" />
        </filter>
 		 <filter levels="error"> 
            <console formatid="colored-error" />
        </filter>
 		<filter levels="warn"> 
            <console formatid="colored-warn"/> 
        </filter>
         <filter levels="debug"> 
            <console formatid="colored-debug"/> 
        </filter>
		<filter levels="info"> 
            <console formatid="colored-info"/>   
        </filter>
		<rollingfile formatid="main" type="date" filename="` + cfg.Root + `/log.txt" datepattern="2006.01.02" fullname="true" maxrolls="` + cfg.Maxdays + `"/>  //the unit of size is byte 
    </outputs>
    <formats>
		<format id="main"  format="%EscM(31)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
        <format id="colored-critical"  format="%EscM(35)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
		<format id="colored-error"  format="%EscM(31)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
        <format id="colored-warn"  format="%EscM(33)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
		<format id="colored-debug"  format="%EscM(32)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
        <format id="colored-info"  format="%EscM(37)%Date %Time [%LEV] [%File:%Line] [%Func]->%EscM(0) [%Msg] %n"/>
    </formats>
</seelog> 
		`
	} else if strings.EqualFold(cfg.Type, "size") {
		useConfig = `
<seelog minlevel="debug" maxlevel="critical">
    <outputs formatid="main">
        <filter levels="` + cfg.Level + `"> 
            <console/>
        </filter>
		<filter levels="` + cfg.Level + `"> 
            <rollingfile type="size" filename="` + cfg.Root + `/log.txt" maxsize="` + cfg.Maxsize + `"maxrolls="` + cfg.Maxfiles + `"/> //the unit of size is byte
        </filter>
    </outputs>
    <formats>
       									%Date %Time [%LEV] [%File:%Line] [%Func]->%Msg%n
		<format id="warning"  format="%EscM(31) %Date %Time [%LEV] [%File:%Line] [%Func]->%Msg%n EscM(0)"/> 
    </formats>
</seelog>
		`
	}
	var err error
	logger, err = seelog.LoggerFromConfigAsBytes([]byte(useConfig))
	if err != nil {
		panic(err)
	}
	err = seelog.ReplaceLogger(logger)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(FLUSH_INTERVAL))
			seelog.Flush()
			if isClose {
				break
			}
		}
	}()
}

func DeInit() {
	isClose = true
	seelog.Flush()
}
