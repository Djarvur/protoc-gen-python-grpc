package flags

func must[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(err)
	}

	return v
}
