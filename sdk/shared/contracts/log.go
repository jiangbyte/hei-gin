package contracts

type LogPersistenceAPI interface {
	SaveLog(entry LogEntry) error
}

type LogEntry struct {
	ID         string
	Category   string
	Name       string
	ExeStatus  string
	ExeMessage string
	OpIP       string
	OpAddress  string
	OpBrowser  string
	OpOS       string
	ReqMethod  string
	ReqURL     string
	ParamJSON  string
	OpTime     string
	TraceID    string
	SignData   string
	OpUser     string
}
