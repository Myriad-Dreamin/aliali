package suppress

type ISuppress interface {
  Suppress(err error)
}


