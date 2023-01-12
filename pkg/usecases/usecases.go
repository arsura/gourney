package usecase

type UseCase struct {
	Post    PostUseCaseProvider
	PostLog PostLogUseCaseProvider
}
