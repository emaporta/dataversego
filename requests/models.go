package requests

type Authorization struct {
	token_type     string
	expires_in     string
	ext_expires_in string
	expires_on     string
	not_before     string
	resource       string
	access_token   string
}
