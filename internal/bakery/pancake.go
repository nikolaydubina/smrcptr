package bakery

import (
	"encoding/xml"
	"fmt"
)

type Pancake struct{}

func NewPancake() Pancake { return Pancake{} }

func (s *Pancake) Fry() {}

func (s Pancake) Bake() {}

// Cake has unnamed receivers
type Cake struct{}

func (*Cake) Fry() {}

func (Cake) Bake() {}

// Brownie has constructor mismatch
type Brownie struct{}

func NewBrownie() Brownie { return Brownie{} }

func (*Brownie) Bake() {}

// Cookie has constructor returning error and is ok
type Cookie struct{}

func NewCookie() (*Cookie, error) { return nil, nil }

func (*Cookie) Bake() {}

// BadCookie has constructor returning error and is not ok
type BadCookie struct{}

func NewBadCookie() (*BadCookie, error) { return nil, nil }

func (BadCookie) Bake() {}

// Oven has only constructor
type Oven struct{}

func NewOven() (*Oven, error) { return nil, nil }

// Teacup defines methods from stanard packages that will be skipped
type Teacup struct{}

func (s Teacup) Name() {}

//encoding
func (s *Teacup) UnmarshalJSON([]byte) error                                { return nil }
func (s *Teacup) UnmarshalText([]byte) error                                { return nil }
func (s *Teacup) UnmarshalBinary([]byte) error                              { return nil }
func (s *Teacup) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error { return nil }
func (s *Teacup) UnmarshalXMLAttr(attr xml.Attr) error                      { return nil }

// database/sql
func (s *Teacup) Scan(src any) error { return nil }

// io
func (s *Teacup) Read(p []byte) (n int, err error)

type TeacupTwo struct{}

func (s TeacupTwo) Name() {}

// fmt
func (s *TeacupTwo) Scan(state fmt.ScanState, verb rune) error { return nil }
