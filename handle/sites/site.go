package sites

type FocusSite interface {
		downLoadSource(url string) (err error)
}