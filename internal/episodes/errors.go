package episodes

type ErrEpisodeNotFound struct{ error }

func (e ErrEpisodeNotFound) Unwrap() error { return e.error }

func NewErrEpisodeNotFound(err error) ErrEpisodeNotFound {
	return ErrEpisodeNotFound{err}
}
