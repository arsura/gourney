package repository

type Repository struct {
	Post    PostRepositoryProvider
	PostLog PostLogRepositoryProvider
}
