package errors

import "errors"

// ErrParseToken - parse token failed error
var ErrParseToken = errors.New("parse token failed")

// ErrNotValidToken - not valid token error
var ErrNotValidToken = errors.New("not valid token")
