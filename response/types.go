package response

import "io"

type ResponseParser interface {
	ParseResponse(io.Reader) error
}
