package errors

var (
	//EInvalidContext signals that a context is no longer valid, but should have been at the time of invocation
	EInvalidContext = PersistentCode(DSynchro, AInvariant)
)
