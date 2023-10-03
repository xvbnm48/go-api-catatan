package log

const (
	//LogTime is log key for timestamp
	LogTime = "ts"
	//LogCaller is log key for source file name
	LogCaller = "caller"
	//LogMethod is log key for method name
	LogMethod = "method"
	//LogUser is log key for user
	LogUser = "user"
	//LogEmail is log key for email
	LogEmail = "email"
	//LogMobile is log key for mobile no
	LogMobile = "mobile"
	//LogRole is log key for role
	LogRole = "role"
	//LogTook is log key for call duration
	LogTook = "took"
	//LogInfo is log key for info
	LogInfo = "[INFO]"
	//LogDebug is log key for debug
	LogDebug = "[DEBUG]"
	//LogCritical is log key for critical
	LogCritical = "[CRITICAL]"
	//LogError is log key for error
	LogError = "[ERROR]"
	//LogBasic is log key for basic log
	LogBasic = "[BASIC]"
	//LogWarning is log key for warning log
	LogWarning = "[WARNING]"
	//LogReq is log key for request log
	LogReq = "[REQUEST]"
	//LogResp is log key for response log
	LogResp = "[RESPONSE]"
	//LogData is log key for data log
	LogData = "[DATA]"
	//LogService is log key for service name
	LogService = "service"
	//LogToken is log key for token
	LogToken = "token"
	//LogExit is log key for exit
	LogExit = "exit"
	//default file logger
	logFile = "service.log"
)
