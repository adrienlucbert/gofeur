package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errAtLeastOneForklift = errors.New("Please provide at least a forklift")
	errAtLeastOneTruck    = errors.New("Please provide at least one truck")
)

type dupPropertyError struct {
	dupPropertiesEntities []string
}

func (err dupPropertyError) Error() string {
	return strings.Join(err.dupPropertiesEntities, "\n")
}

func verifySimulationValidity(simulation simulation) error {
	if len(simulation.warehouse.forklifts) == 0 {
		return errAtLeastOneForklift
	}

	if len(simulation.warehouse.trucks) == 0 {
		return errAtLeastOneTruck
	}
	entities := makeEntitiesArray(simulation)

	err := checkForOutOfWarehouseBoundEntity(entities, simulation.warehouse)
	if err != nil {
		return err
	}

	err = ensureNoStackedEntities(entities)
	if err != nil {
		return err
	}

	return ensureForDuplicatedEntitiyName(entities)
}

func makeEntitiesArray(simulation simulation) []entity {
	nbEntities := len(simulation.warehouse.parcels) + len(simulation.warehouse.forklifts) + len(simulation.warehouse.trucks)
	entities := make([]entity, 0, nbEntities)

	for _, parcel := range simulation.warehouse.parcels {
		entities = append(entities, parcel)
	}

	for _, forklift := range simulation.warehouse.forklifts {
		entities = append(entities, forklift)
	}

	for _, truck := range simulation.warehouse.trucks {
		entities = append(entities, truck)
	}

	return entities
}

type outOfBoundError struct {
	entity entity
}

func (err outOfBoundError) Error() string {
	return fmt.Sprintf("The %s named %s is out of bound", err.entity.Kind(), err.entity.Name())
}

func checkForOutOfWarehouseBoundEntity(entities []entity, warehouse warehouse) error {
	for _, entity := range entities {
		coord := entity.Coord()

		if !(coord.X < warehouse.width && coord.Y < warehouse.width) {
			return outOfBoundError{entity: entity}
		}
	}

	return nil
}

type stackedEntitiesError struct {
	err error
}

func (err stackedEntitiesError) Error() string {
	return fmt.Sprintf("Error found stacked entities:\n%s", err.err.Error())
}

func ensureNoStackedEntities(entities []entity) error {
	err := hasEntityPropertyDup(entities, entity.Coord)
	if err != nil {
		return stackedEntitiesError{err: err}
	}
	return nil
}

type dupEntityNameError struct {
	err error
}

func (err dupEntityNameError) Error() string {
	return fmt.Sprintf("Error found duplicated entities name:\n%s", err.err.Error())
}

func ensureForDuplicatedEntitiyName(entities []entity) error {
	err := hasEntityPropertyDup(entities, entity.Name)
	if err != nil {
		return dupEntityNameError{err: err}
	}
	return nil
}

func hasEntityPropertyDup[T interface {
	comparable
	fmt.Stringer
}](entities []entity, propertyGettor func(entity) T,
) error {
	board := make(map[T][]*entity, len(entities))

	for i := range entities {
		property := propertyGettor(entities[i])

		if board[property] == nil {
			board[property] = make([]*entity, 0, 1)
		}
		board[property] = append(board[property], &entities[i])
	}

	for coord, entities := range board {
		if len(entities) == 1 {
			delete(board, coord)
		}
	}

	if len(board) != 0 {
		err := make([]string, 0, len(board))
		for property, entities := range board {
			errEntities := make([]string, 0, len(entities))
			for _, entity := range entities {
				kind := (*entity).Kind()
				name := (*entity).Name()

				errEntities = append(errEntities, fmt.Sprintf("%s: %s", kind, name))
			}
			err = append(err, fmt.Sprintf("%s: %s", property, strings.Join(errEntities, ", ")))
		}
		return dupPropertyError{dupPropertiesEntities: err}
	}
	return nil
}
