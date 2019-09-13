package def

type Executor interface {
	Output() string
	HasRun() bool
	Run() error
}
