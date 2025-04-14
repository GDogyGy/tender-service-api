package model

import (
	"errors"
)

var NotFound = errors.New("not found")
var NotFindResponsible = errors.New("not found responsible")
var AlreadyVotedResponsible = errors.New("already voted responsible")
var BadStatus = errors.New("status in object is wrong")
