package domain

import (
	"time"

	"happyday/abstract"
)

type (
	AggregateRoot interface {
		abstract.AggregateRoot[*State]
		Create(name string, price float64, priority int64, products []Product, checkIfExists func(Product) (bool, error)) error
		ChangeName(name string) error
		ChangePrice(price float64) error
		Enable() error
		Disable() error
		ChangePriority(priority int64) error
		ChangeProducts(products []Product, checkIfExists func(Product) (bool, error)) error
	}

	aggregateRoot struct {
		abstract.DefaultAggregateRoot[*State]
	}
)

var _ AggregateRoot = (*aggregateRoot)(nil)

func NewAggregateRoot(state *State, version int64) *aggregateRoot {
	return &aggregateRoot{
		abstract.NewAggregateRoot(state, version),
	}
}

func (root *aggregateRoot) Create(name string, price float64, priority int64, products []Product, checkIfExists func(Product) (bool, error)) error {
	if err := checkName(name); err != nil {
		return err
	}

	if err := checkPrice(price); err != nil {
		return err
	}

	if err := checkProducts(products, checkIfExists); err != nil {
		return err
	}

	root.On(Created{
		name:       name,
		price:      price,
		isEnable:   false,
		priority:   priority,
		products:   products,
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func checkName(name string) error {
	if len(name) == 0 {
		return NameIsEmpty
	}

	if len(name) > 100 {
		return NameIsTooLarge
	}

	return nil
}

func (root *aggregateRoot) ChangeName(name string) error {
	if err := checkName(name); err != nil {
		return err
	}

	root.On(NameChanged{
		name:       name,
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func checkPrice(price float64) error {
	if price < 0 {
		return PriceIsInvalid
	}

	return nil
}

func (root *aggregateRoot) ChangePrice(price float64) error {
	if err := checkPrice(price); err != nil {
		return err
	}

	root.On(PriceChanged{
		price:      price,
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func (root *aggregateRoot) Enable() error {
	root.On(Enabled{
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func (root *aggregateRoot) Disable() error {
	root.On(Disabled{
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func (root *aggregateRoot) ChangePriority(priority int64) error {
	root.On(PriorityChanged{priority: priority, occurredOn: time.Now().UTC()})
	return nil
}

func checkProducts(products []Product, checkIfExists func(Product) (bool, error)) error {
	for _, product := range products {
		exist, err := checkIfExists(product)
		if err != nil {
			return err
		}

		if !exist {
			return ProductNotExist
		}
	}

	return nil
}
func (root *aggregateRoot) ChangeProducts(products []Product, checkIfExists func(Product) (bool, error)) error {
	if err := checkProducts(products, checkIfExists); err != nil {
		return err
	}

	root.On(ProductsChanged{products: products, occurredOn: time.Now().UTC()})
	return nil
}
