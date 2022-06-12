//go:build grpc

package grpc

func Module(namespace string) dep.Module {
	return dep.Declare(
		dep.Name("wk", "grpc", namespace),
	)
}
