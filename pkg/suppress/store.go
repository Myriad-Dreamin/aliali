package suppress

import (
	"fmt"
	"runtime"
)

type Store struct {
	Err error
}

func (s *Store) Suppress(err error) {
	if err != nil {
		if s.Err == nil {
			s.Err = err
		} else {
			fmt.Println(s.Err, "meets another error:")
			panic(err)
		}
	}
}

func (s *Store) Restore() error {
	return s.Err
}

func (s *Store) WarnOnce(err error) {
	if err != nil {
		// s.Warnings = append(s.Warnings, err)
		var b = make([]byte, 1024)
		b = b[:runtime.Stack(b, false)]
		fmt.Printf("warning occurs: %s:\n%s", err.Error(), string(b))
	}
}
