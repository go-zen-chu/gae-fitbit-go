package fitbitauth

type Factory interface {
	FileStore() (Store, error)
}
