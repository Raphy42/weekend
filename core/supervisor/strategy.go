package supervisor

type RestartStrategy int

const (
	// PermanentRestartStrategy the child process is always restarted.
	PermanentRestartStrategy RestartStrategy = iota
	// TemporaryRestartStrategy the child process is never restarted, regardless of the supervision strategy.
	// Any termination (even abnormal) is considered successful.
	TemporaryRestartStrategy
	// TransientRestartStrategy child process is restarted only if it terminates abnormally
	TransientRestartStrategy
)

type ShutdownStrategy int

const (
	ImmediateShutdownStrategy ShutdownStrategy = iota
	TimeoutShutdownStrategy
)

type SupervisionStrategy int

const (
	// OneForOneSupervisionStrategy if a child process terminates, only that process is restarted.
	OneForOneSupervisionStrategy SupervisionStrategy = iota
	// OneForAllSupervisionStrategy  if a child process terminates, all other child processes are terminated and then
	// all child processes (including the terminated one) are restarted.
	OneForAllSupervisionStrategy
)
