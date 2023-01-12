package main

import (
	"errors"
	"fmt"
	"strings"
)

func verifySimulationValidity(simulation Simulation) error {
	if len(simulation.warehouse.forklifts) == 0 {
		return errors.New("Please provide at least a forklift")
	}

	if len(simulation.warehouse.trucks) == 0 {
		return errors.New("Please provide at least one truck")
	}
	var entities = makeEntitiesArray(simulation)

	var err = checkForOutOfWarehouseBoundEntity(entities, simulation.warehouse)

	if err != nil {
		return err
	}

	err = ensureNoStackedEntities(entities)
	if err != nil {
		return err
	}

	return ensureForDuplicatedEntitiyName(entities)
}

func makeEntitiesArray(simulation Simulation) []Entity {
	var nb_entities = len(simulation.warehouse.parcels) + len(simulation.warehouse.forklifts) + len(simulation.warehouse.trucks)
	var entities = make([]Entity, 0, nb_entities)

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

func checkForOutOfWarehouseBoundEntity(entities []Entity, warehouse Warehouse) error {
	for _, entity := range entities {
		var coord = entity.Coord()

		if !(coord.X < warehouse.width && coord.Y < warehouse.width) {
			return errors.New(fmt.Sprintf("The %s named %s is out of bound", entity.Kind(), entity.Name()))
		}
	}

	return nil
}

func ensureNoStackedEntities(entities []Entity) error {
	var err = hasEntityPropertyDup(entities, Entity.Coord)

	if err != nil {
		return fmt.Errorf("Error found stacked entities:%w", err)
	}
	return nil
}

func ensureForDuplicatedEntitiyName(entities []Entity) error {
	var err = hasEntityPropertyDup(entities, Entity.Name)

	if err != nil {
		return fmt.Errorf("Error found duplicated entities name:%w", err)
	}
	return nil
}

func hasEntityPropertyDup[T interface {
	comparable
	fmt.Stringer
}](entities []Entity, propertyGettor func(Entity) T) error {
	var board = make(map[T][]*Entity, len(entities))

	for i := range entities {
		var entity = entities[i]
		var property = propertyGettor(entity)

		if board[property] == nil {
			board[property] = make([]*Entity, 0, 1)
		}
		board[property] = append(board[property], &entities[i])
	}

	for coord, entities := range board {
		if len(entities) == 1 {
			delete(board, coord)
		}
	}

	if len(board) != 0 {
		var err string
		for property, entities := range board {
			err += fmt.Sprintf("\n  %s ", property)
			var err_entities = make([]string, 0, len(entities))
			for _, entity := range entities {
				var kind = (*entity).Kind()
				var name = (*entity).Name()

				err_entities = append(err_entities, fmt.Sprintf("%s: %s", kind, name))
			}
			err += strings.Join(err_entities, ", ")
		}
		return errors.New(err)
	}
	return nil
}
