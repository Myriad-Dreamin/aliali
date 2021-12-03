package suppress

type ISuppress interface {
	Suppress(err error)
	WarnOnce(err error)
}
