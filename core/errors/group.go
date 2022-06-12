package errors

import "github.com/palantir/stacktrace"

type Group []error

func NewGroup() []error {
	return make(Group, 0)
}

func (g *Group) Append(errs ...error) {
	*g = append(*g, errs...)
}

func (g Group) First() error {
	for _, err := range g {
		if err != nil {
			return err
		}
	}
	return nil
}

func (g Group) Coalesce() error {
	rootErr := stacktrace.NewError("one or more errors occurred in the `errors.Group`")
	none := true

	for _, err := range g {
		if err != nil {
			none = false
			rootErr = stacktrace.Propagate(rootErr, stacktrace.RootCause(err).Error())
		}
	}
	if none {
		return nil
	}
	return rootErr
}
