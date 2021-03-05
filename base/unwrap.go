package base

func UnwrapOnce(err error) error {
	switch e := err.(type) {
	case interface{ Cause() error }:
		return e.Cause()
	case interface{ Unwrap() error }:
		return e.Unwrap()
	}

	return nil
}

func UnwrapAll(err error) error {
	for {
		cause := UnwrapOnce(err)

		if cause == nil {
			break
		}

		err = cause
	}

	return err
}
