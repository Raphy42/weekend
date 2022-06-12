package redis

var (
	ModuleName = dep.Name("wk", "redis")
)

func Module() dep.Module {
	return dep.Declare(
		ModuleName,
	)
}
