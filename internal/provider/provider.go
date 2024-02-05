package provider

type Provider interface {
	Links() ([]string, error)
}
