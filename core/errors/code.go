package errors

import (
	"github.com/palantir/stacktrace"

	"github.com/Raphy42/weekend/pkg/bitmask"
)

const (
	//KTransient is an error kind that may happen outside the application invariants (network failure, service down)
	// transient errors must only be associated with retryable idempotent domain
	KTransient = 0x0000
	//KPersistent is an error kind that is considered inside the application invariants (bad input, wrong initialisation order)
	// persistent errors are generally unrecoverable due to bad domain invariants, or programming errors
	KPersistent = 0x7000

	//DResource is an error domain for resources (database model, configuration, etc...)
	DResource = 0x0100
	//DService is an error domain for services (clients, database drivers, etc...)
	DService = 0x0200
	//DTemporal is an error domain for everything `time` related (timeout, tickers, etc...)
	DTemporal = 0x0300
	//DSize is an error domain for sizeable items (array, buffer, etc...)
	DSize = 0x0400
	//DLogic is an error domain for logic (loops, control-flow, etc...)
	DLogic = 0x0500
	//DIO is an error domain for IO (file, http-transport, etc...)
	DIO = 0x0600
	//DUser is an error domain related to user interactions (form submission, user inputs, etc...)
	DUser = 0x0700
	//DEncoding is an error domain related to serialization (JSON, text, etc...)
	DEncoding = 0x0800
	//DDependency is an error domain related to dependencies, mostly used by the di system
	DDependency = 0x0900
	//DSynchro is an error domain related to synchronisation semantics (context, mutex, etc...)
	DSynchro = 0x0a00
	//DType is an error domain related to the type system (any, reflection, etc...)
	DType = 0x0b00
	//DValue is an error domain related to values
	DValue = 0x0c00

	//ANotFound is an error axiom for missing or not found items
	ANotFound = 0x0010
	//AMalformed is an error axiom for malformed or corrupted items
	AMalformed = 0x0020
	//AInvariant is an error axiom for invariants (eg: nil check)
	AInvariant = 0x0030
	//ATooBig is an error axiom mostly for sizeable items (mostly used in conjunction with `DSize`)
	ATooBig = 0x0040
	//ATooSmall is an error axiom mostly for sizeable items (mostly used in conjunction with `DSize`)
	ATooSmall = 0x0050
	//AUnexpected is an error axiom used to mark an interaction as unwanted
	AUnexpected = 0x0060
	//AUnreachable is an error axiom used to define unreachable areas of code or logic
	AUnreachable = 0x0070
	//AUnimplemented is an error axiom for missing domain
	AUnimplemented = 0x0080
	//ANil is an error axiom for nil values
	ANil = 0x0090
	//ATimeout is an error relating to time outs
	ATimeout = 0x00a0
)

func init() {
	globalRegistry.
		RegisterDomain(DResource, "resource").
		RegisterDomain(DService, "service").
		RegisterDomain(DTemporal, "temporal").
		RegisterDomain(DSize, "size").
		RegisterDomain(DLogic, "logic").
		RegisterDomain(DIO, "io").
		RegisterDomain(DUser, "user").
		RegisterDomain(DEncoding, "encoding").
		RegisterDomain(DDependency, "dependency").
		RegisterDomain(DSynchro, "sync").
		RegisterDomain(DType, "type").
		RegisterDomain(DValue, "value")
	globalRegistry.
		RegisterAxiom(ANotFound, "not_found").
		RegisterAxiom(AMalformed, "malformed").
		RegisterAxiom(AInvariant, "invariant").
		RegisterAxiom(ATooBig, "too_big").
		RegisterAxiom(ATooSmall, "too_small").
		RegisterAxiom(AUnexpected, "unexpected").
		RegisterAxiom(AUnreachable, "unreachable").
		RegisterAxiom(AUnimplemented, "unimplemented").
		RegisterAxiom(ANil, "nil").
		RegisterAxiom(ATimeout, "timeout")
}

// Error code bit structure
// Maximum:
// 0 x 7	f	f	f
// Bit masks:
// 0 x 0	0	0	0
//     |	|	|	|
//	   |	|	|	available for user
//     |	|	axiom code (not_found, timed_out, unreachable, etc...)
//     |    domain code (resource, service, io, etc...)
//	   kind code (transient, persistent)

func ResourceNotFound(bit uint16) uint16 {
	return bitmask.Set(bit, ANotFound)
}

func Transient(bit uint16) uint16 {
	return bitmask.Clear(bit, KTransient)
}

func Persistent(bit uint16) uint16 {
	return bitmask.Set(bit, KPersistent)
}

func IsPersistentCode(bit uint16) bool {
	return bitmask.Has(bit, KPersistent)
}

func IsTransientCode(bit uint16) bool {
	return !IsPersistentCode(bit)
}

func Code(kind, ns, axiom uint16) stacktrace.ErrorCode {
	return stacktrace.ErrorCode(bitmask.Set(kind, bitmask.Set(ns, axiom)))
}

func PersistentCode(ns, axiom uint16) stacktrace.ErrorCode {
	return stacktrace.ErrorCode(bitmask.Set(KPersistent, bitmask.Set(ns, axiom)))
}

func TransientCode(ns, axiom uint16) stacktrace.ErrorCode {
	return stacktrace.ErrorCode(bitmask.Set(KTransient, bitmask.Set(ns, axiom)))
}

func IsNotFoundCode(bit uint16) bool {
	return bitmask.Has(bit, ANotFound)
}
