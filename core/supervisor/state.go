package supervisor

type State struct {
	Retries   int
	LastError error
}
