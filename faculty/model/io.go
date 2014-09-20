// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package model

import (
	// "fmt"

	"github.com/gocircuit/escher/be"
)

/*
	{ In *, Out *, _ * }
*/
type IO struct{}

func (IO) Spark() {}

func (IO) Cognize_(eye *be.Eye, v interface{}) {
	eye.Show("In", v)
}

func (IO) CognizeIn(eye *be.Eye, v interface{}) {}

func (IO) CognizeOut(eye *be.Eye, v interface{}) {
	eye.Show("_", v)
}
