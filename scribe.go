package scribe

import "os"
import "log"
import "time"

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

var logger = log.New (
        os.Stdout, "",
        log.Ldate |
        log.Ltime |
        log.Lmsgprefix)

var loggingToDirectory bool
var logDirectory       string
var previousDay        int

var currentFile *os.File

func SetLogLevel (level LogLevel) {
        levelGate = level
}

func SetLogDirectory (directoryName string) {
        if directoryName[len(directoryName) - 1] != '/' {
                directoryName += "/"
        }
        logDirectory       = directoryName
        loggingToDirectory = true
}

func UnsetLogDirectory () {
        loggingToDirectory = false
        logger.SetOutput(os.Stdout)
}

func ListenOnce () {
        message := <- queue

        switch message.Type {
                case Progress:   logger.SetPrefix("... ")
                case Done:       logger.SetPrefix(".// ")
                case Info:       logger.SetPrefix("(i) ")
                case Warning:    logger.SetPrefix("!!! ")
                case Error:      logger.SetPrefix("ERR ")
                case Fatal:      logger.SetPrefix("XXX ")
                case Request:    logger.SetPrefix("->? ")
                case Resolve:    logger.SetPrefix("->! ")
                case Connect:    logger.SetPrefix("--> ")
                case Mount:      logger.SetPrefix("-=E ")
                case Disconnect: logger.SetPrefix("<-- ")
                case Unmount:    logger.SetPrefix("X=- ")
                case Bind:       logger.SetPrefix("=#= ")
                case Unbind:     logger.SetPrefix("=X= ")
        }

        now := time.Now()
        currentDay := (now.Year() - 1970) * 365 + now.YearDay()

        if loggingToDirectory && currentDay > previousDay {
                if currentFile != nil {
                        currentFile.Close()
                }

                var err error
                currentFile, err = os.OpenFile (
                        logDirectory + now.Format("2006-01-02.log"),
                        os.O_WRONLY |
                        os.O_APPEND |
                        os.O_CREATE,
                        0660)

                if err != nil {
                        PrintError(LogLevelError, "could not open log file")
                } else {
                        logger.SetOutput(currentFile)
                }
        }
        
        logger.Println(message.Content...)
        previousDay = currentDay
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
