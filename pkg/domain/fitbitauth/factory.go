package fitbitauth

// Factory : factory for fitbitauth package
type Factory interface {
	FileStore() (Store, error)
}
