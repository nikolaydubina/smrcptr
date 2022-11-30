package bakery

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
