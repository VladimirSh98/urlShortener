package errors

import "errors"

var ErrParseToken = errors.New("parse token failed")

var ErrNotValidToken = errors.New("not valid token")
