package suppress

type PanicAll struct {
}

func (PanicAll) Suppress(err error) {
	if err != nil {
		panic(err)
	}
}

func (PanicAll) Restore() error {
	return nil
}

func (PanicAll) WarnOnce(err error) {
	if err != nil {
		panic(err)
	}
}
