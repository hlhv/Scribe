package scribe

import (
        "log"
)

type MessageType int

const (
        Progress MessageType = iota
        Done
        Info
        Warning
        Error
        Fatal
        Request
        Resolve
        Connect
        Mount
        Disconnect
        Unmount
        Bind
        Unbind
)

type Message struct {
        Type    MessageType
	Level   LogLevel
        Content []interface{}
}

type LogLevel int

const (
	// print any and all info for debugging purposes.
	LogLevelDebug  LogLevel = 0
	// only print some info, such as logging requests and connections, and
	// errors
	LogLevelNormal LogLevel = 1
	// only print errors
	LogLevelError  LogLevel = 2
	// completely disable logging. this is not reccomended.
	LogLevelNone   LogLevel = 3
)

var queue chan Message = make(chan Message, 16)
var levelGate LogLevel = LogLevelNormal

func SetLogLevel (level LogLevel) {
        levelGate = level
}

func ListenOnce () {
        message := <- queue

        content := message.Content
        content = append(content, "")
        copy(content[1:], content)

        switch message.Type {
        case Progress:
                content[0] = "..."
                log.Println(content...)
                break
        case Done:
                content[0] = ".//"
                log.Println(content...)
                break
        case Info:
                content[0] = "(i)"
                log.Println(content...)
                break
        case Warning:
                content[0] = "!!!"
                log.Println(content...)
                break
        case Error:
                content[0] = "ERR"
                log.Println(content...)
                break
        case Fatal:
                content[0] = "XXX"
                log.Fatalln(content...)
                break
        case Request:
                content[0] = "->?"
                log.Println(content...)
                break
        case Resolve:
                content[0] = "->!"
                log.Println(content...)
                break
        case Connect:
                content[0] = "-->"
                log.Println(content...)
                break
        case Mount:
                content[0] = "-=E"
                log.Println(content...)
                break
        case Disconnect:
                content[0] = "<--"
                log.Println(content...)
                break
        case Unmount:
                content[0] = "X=-"
                log.Println(content...)
                break
        case Bind:
                content[0] = "=#="
                log.Println(content...)
                break
        case Unbind:
                content[0] = "=X="
                log.Println(content...)
                break
        }
}

func Print (t MessageType, level LogLevel, content ...interface{}) {
        // if the message level isn't large enough, don't send it
        if levelGate > level { return }
        
        queue <- Message {
                Type:    t,
                Content: content, 
                Level:   level,
        }
}

func PrintProgress   (level LogLevel, content ...interface{}) {
        Print(Progress,   level, content...)
}

func PrintDone       (level LogLevel, content ...interface{}) {
        Print(Done,       level, content...)
}

func PrintInfo       (level LogLevel, content ...interface{}) {
        Print(Info,       level, content...)
}

func PrintWarning    (level LogLevel, content ...interface{}) {
        Print(Warning,    level, content...)
}

func PrintError      (level LogLevel, content ...interface{}) {
        Print(Error,      level, content...)
}

func PrintFatal      (level LogLevel, content ...interface{}) {
        Print(Fatal,      level, content...)
}

func PrintRequest    (level LogLevel, content ...interface{}) {
        Print(Request,    level, content...)
}

func PrintResolve    (level LogLevel, content ...interface{}) {
        Print(Resolve,    level, content...)
}

func PrintConnect    (level LogLevel, content ...interface{}) {
        Print(Connect,    level, content...)
}

func PrintMount      (level LogLevel, content ...interface{}) {
        Print(Mount,      level, content...)
}

func PrintDisconnect (level LogLevel, content ...interface{}) {
        Print(Disconnect, level, content...)
}

func PrintUnmount    (level LogLevel, content ...interface{}) {
        Print(Unmount,    level, content...)
}

