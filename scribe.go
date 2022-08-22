package scribe

import "os"
import "log"
import "time"

/* Different types of message. These determine what kind of symbol is displayed
 * on a message.
 */
type MessageType int

/* All message types.
 */
const (
        // ...
        Progress MessageType = iota
        // .//
        Done
        // (i)
        Info
        // !!!
        Warning
        // ERR
        Error
        // XXX
        Fatal
        // ->?
        Request
        // ->!
        Resolve
        // -->
        Connect
        // -=E
        Mount
        // <--
        Disconnect
        // X=-
        Unmount
        // =#=
        Bind
        // =X=
        Unbind
)

/* Message structs are sent down Scribe's message channel. They contain
 * information about the message to be logged.
 */
type Message struct {
        // The type of the message
        Type    MessageType
        // The importance of a message
	Level   LogLevel
	// The message content (what gets printed)
        Content []interface{}
}

/* LogLevel specifies what messages are to be logged, and what messages are to
 * be ignored. */
type LogLevel int

/* All log levels.
 */
const (
	// Print any and all info for debugging purposes.
	LogLevelDebug  LogLevel = 0
	// Only print some info, such as logging requests and connections, and
	// errors
	LogLevelNormal LogLevel = 1
	// Only print errors
	LogLevelError  LogLevel = 2
	// Completely disable logging. You DO NOT want to do this!
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

func init () {
	go func () {
		for {
			listenOnce()
		}	
	} ()
}

/* SetLogLevel sets the log level. Only messages with the specified log level
 * or higher will be logged.
 */
func SetLogLevel (level LogLevel) {
        levelGate = level
}

/* SetLogDirectory sets the directory logs are to be written to. In this
 * directory, log files will be created with a name formatted like this:
 *
 * YYYY-MM-DD.log
 *
 * When a message is logged, and a day has passed since the last message, a new
 * file is created.
 */
func SetLogDirectory (directoryName string) {
        if directoryName[len(directoryName) - 1] != '/' {
                directoryName += "/"
        }
        logDirectory       = directoryName
        loggingToDirectory = true
}

/* UnsetLogDirectory causes the logging system to stop writing to files, and
 * sets all output to stdout.
 */
func UnsetLogDirectory () {
        loggingToDirectory = false
        logger.SetOutput(os.Stdout)
}

/* listenOnce listens for one message. This function is blocking and should be
 * run in a loop.
 */
func listenOnce () {
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

        updateCurrentFile()
        logger.Println(message.Content...)
}

func updateCurrentFile () {
        now := time.Now()
        currentDay := (now.Year() - 1970) * 365 + now.YearDay()

        if loggingToDirectory && currentDay > previousDay {
                if currentFile != nil {
                        // TODO: compress the previous file
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
        
        previousDay = currentDay
}

/* Print queues a single message to be logged/
 */
func Print (t MessageType, level LogLevel, content ...interface{}) {
        // if the message level isn't large enough, don't send it
        if levelGate > level { return }
        
        queue <- Message {
                Type:    t,
                Content: content, 
                Level:   level,
        }
}

/* PrintProgress logs a message of type progress. This should be used to log
 * steps of an overall task.
 */
func PrintProgress (level LogLevel, content ...interface{}) {
        Print(Progress, level, content...)
}

/* PrintDone logs a message of type done. This should be used to log the
 * successful completion of a task.
 */
func PrintDone (level LogLevel, content ...interface{}) {
        Print(Done, level, content...)
}

/* PrintDone logs a message of type info. This should be used to log
 * informational messages and general observations the program makes that could
 * be useful to whoever is reading the log.
 */
func PrintInfo (level LogLevel, content ...interface{}) {
        Print(Info, level, content...)
}

/* PrintWarning logs a message of type warning. This should be used to log a
 * problem which did not cause an error. Useful for logging suspicious behavior
 * or insecure configuration parameters.
 */
func PrintWarning (level LogLevel, content ...interface{}) {
        Print(Warning, level, content...)
}

/* PrintError logs a message of type error. This should be used when an
 * operation encounters a problem. This typically causes the operation to halt.
 * This should not log an event that causes a goroutine or the entire program
 * to stop.
 */
func PrintError (level LogLevel, content ...interface{}) {
        Print(Error, level, content...)
}

/* PrintFatal logs a message of type fatal. This should be used to log errors
 * that caused a goroutine, internal communication channel, or the entire
 * program to terminate.
 */
func PrintFatal (level LogLevel, content ...interface{}) {
        Print(Fatal, level, content...)
}

/* PrintRequest logs a message of type request. This should be used to log an
 * external request made by a client.
 */
func PrintRequest (level LogLevel, content ...interface{}) {
        Print(Request, level, content...)
}

/* PrintResolve logs a message of type resolve. This should be used when a
 * request made by a client was resolved successfully.
 */
func PrintResolve (level LogLevel, content ...interface{}) {
        Print(Resolve, level, content...)
}

/* PrintConnect logs a message of type connect. This should be used to log a
 * client connecting.
 */
func PrintConnect (level LogLevel, content ...interface{}) {
        Print(Connect, level, content...)
}

/* PrintMount logs a message of type mount. This should be used to log a client
 * mounting successfully.
 */
func PrintMount (level LogLevel, content ...interface{}) {
        Print(Mount, level, content...)
}

/* PrintDisconnect logs a message of type disconnect. This should be used to log
 * a client disconnecting.
 */
func PrintDisconnect (level LogLevel, content ...interface{}) {
        Print(Disconnect, level, content...)
}

/* PrintUnmount logs a message of type unmount This should be used to log a
 * client unmounting successfully.
 */
func PrintUnmount (level LogLevel, content ...interface{}) {
        Print(Unmount, level, content...)
}
