package suppress

type PanicAll struct {
}

func (PanicAll) Suppress(err error) {
  if err != nil {
    panic(err)
  }
}