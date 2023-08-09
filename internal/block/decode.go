package block

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

var ErrDecodeJSON = errors.New("failed to decode json of your rules")
var ErrDecodeMapStructure = errors.New("failed to decode mapstructure")
var ErrDecodeGuard = errors.New("failed to decode guard")

type GuardDecoder interface {
	Decode(interface{}) (DecodedGuard, error)
}

type InterfaceGuardDecoder struct{}

func (decoder *InterfaceGuardDecoder) Decode(input interface{}) (DecodedGuard, error) {
	var headerGuard HeaderGuard
	err := mapstructure.Decode(input, &headerGuard)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeMapStructure, err)
	}

	if headerGuard.IsValid() {
		return &headerGuard, nil
	}

	var queryParamGuard QueryParamGuard
	err = mapstructure.Decode(input, &queryParamGuard)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeMapStructure, err)
	}

	if queryParamGuard.IsValid() {
		return &queryParamGuard, nil
	}

	var methodGuard MethodGuard
	err = mapstructure.Decode(input, &methodGuard)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeMapStructure, err)
	}

	if methodGuard.IsValid() {
		return &methodGuard, nil
	}

	var pathGuard PathGuard
	err = mapstructure.Decode(input, &pathGuard)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeMapStructure, err)
	}

	if pathGuard.IsValid() {
		return &pathGuard, nil
	}

	return nil, fmt.Errorf("%w : %#v", ErrDecodeGuard, input)
}

func GuardsFromInterface(jsonData [][]interface{}, decoder GuardDecoder) (Guard, error) {
	var collections []Guard
	for _, rulestToJoin := range jsonData {
		joinedRules := []Guard{}

		for _, rule := range rulestToJoin {
			guard, err := decoder.Decode(rule)
			if err != nil {
				return nil, err
			}

			joinedRules = append(joinedRules, guard)
		}

		collections = append(collections, NewGuardsJoiner(joinedRules))
	}

	return NewGuardsCollection(collections), nil
}
