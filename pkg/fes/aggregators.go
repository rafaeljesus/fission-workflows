package fes

import (
	"errors"
	"reflect"
)

// AggregatorMixin is a helper to implement most of the Aggregator interface.
//
// Structs using this struct will only need to implement the following methods:
// - ApplyEvent(event)
type AggregatorMixin struct {
	aggregate Aggregate
	parent    Aggregator
}

func (am *AggregatorMixin) Aggregate() Aggregate {
	return am.aggregate
}

// For improved performance override this method with a aggregate-specific one
func (am *AggregatorMixin) UpdateState(newState Aggregator) error {
	if newState.Aggregate() != am.Aggregate() {
		return errors.New("invalid newState")
	}

	n := reflect.Indirect(reflect.ValueOf(newState))
	old := reflect.Indirect(reflect.ValueOf(am.parent))

	for i := 0; i < old.NumField(); i++ {
		updatedField := n.Field(i)

		field := old.Field(i)
		if field.IsValid() {
			if field.CanSet() {
				field.Set(updatedField)
			}
		}
	}
	return nil
}

func NewAggregatorMixin(thiz Aggregator, aggregate Aggregate) *AggregatorMixin {
	return &AggregatorMixin{
		aggregate: aggregate,
		parent:    thiz,
	}
}
