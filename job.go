package dispatcher

// Job : Common interface that processes raw data
type Job interface {
	Process()
}
