package parsing

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errAtLeastOneForklift = errors.New("Please provide at least a forklift")
	errAtLeastOneTruck    = errors.New("Please provide at least one truck")
)

// VerifySimulationValidity return an error if one of the following condition is meet:
//
//   - there is no forklift in `simulation`
//   - there is no truck in `simulation`
//   - two entities are on the same grid cell
//   - two entities bears the same name
func VerifySimulationValidity(simulation Simulation) error {
	if len(simulation.Warehouse.Forklifts) == 0 {
		return errAtLeastOneForklift
	}

	if len(simulation.Warehouse.Trucks) == 0 {
		return errAtLeastOneTruck
	}

	err := checkWarehouseSize(simulation.Warehouse)
	if err != nil {
		return err
	}

	err = ensureTrucksAreOnAWarehouseSide(simulation.Warehouse.Trucks, simulation.Warehouse)
	if err != nil {
		return err
	}

	entities := makeEntitiesArray(simulation)

	err = checkForOutOfWarehouseBoundEntity(entities, simulation.Warehouse)
	if err != nil {
		return err
	}

	err = ensureNoStackedEntities(entities)
	if err != nil {
		return err
	}

	return ensureForDuplicatedEntitiyName(entities)
}

type tooSmallWarehouseError struct {
	size gridUnit
}

func (err tooSmallWarehouseError) Error() string {
	return fmt.Sprintf("too small warehouse (%d)", err.size)
}

func checkWarehouseSize(warehouse Warehouse) error {
	size := warehouse.Length * warehouse.Width
	var minimumWarehouseSize gridUnit = 2

	if size < minimumWarehouseSize {
		return tooSmallWarehouseError{size: size}
	}
	return nil
}

type notOnSideTrucksError struct {
	trucks []Truck
}

func (err notOnSideTrucksError) Error() string {
	output := fmt.Sprintf("Error found %d truck(s) not on a side of the warehouse:\n", len(err.trucks))

	trucks := make([]string, 0, len(err.trucks))
	for _, truck := range err.trucks {
		trucks = append(trucks, fmt.Sprintf("  %s: %s", truck.coordinate, truck.Name))
	}

	output += strings.Join(trucks, "\n")
	return output
}

func ensureTrucksAreOnAWarehouseSide(trucks []Truck, warehouse Warehouse) error {
	min := coordinate{}
	max := coordinate{X: warehouse.Width - 1, Y: warehouse.Length - 1}

	errTrucks := make([]Truck, 0)
	for _, truck := range trucks {
		isOnLeftOrRightSize := truck.X == min.X || truck.X == max.X
		isOnTopOrBottomSize := truck.Y == min.Y || truck.Y == max.Y
		isOnASide := isOnLeftOrRightSize || isOnTopOrBottomSize

		if !isOnASide {
			errTrucks = append(errTrucks, truck)
		}
	}

	if len(errTrucks) > 0 {
		return notOnSideTrucksError{trucks: errTrucks}
	}
	return nil
}

func makeEntitiesArray(simulation Simulation) []entity {
	nbEntities := len(simulation.Warehouse.Parcels) + len(simulation.Warehouse.Forklifts) + len(simulation.Warehouse.Trucks)
	entities := make([]entity, 0, nbEntities)

	for _, parcel := range simulation.Warehouse.Parcels {
		entities = append(entities, parcel)
	}

	for _, forklift := range simulation.Warehouse.Forklifts {
		entities = append(entities, forklift)
	}

	for _, truck := range simulation.Warehouse.Trucks {
		entities = append(entities, truck)
	}

	return entities
}

type outOfBoundError struct {
	entity entity
}

func (err outOfBoundError) Error() string {
	return fmt.Sprintf("The %s named %s is out of bound", err.entity.kind(), err.entity.stringerName())
}

func checkForOutOfWarehouseBoundEntity(entities []entity, warehouse Warehouse) error {
	for _, entity := range entities {
		coord := entity.coord()

		if !(coord.X < warehouse.Width && coord.Y < warehouse.Width) {
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
	err := hasEntityPropertyDup(entities, entity.coord)
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
	err := hasEntityPropertyDup(entities, entity.stringerName)
	if err != nil {
		return dupEntityNameError{err: err}
	}
	return nil
}

type dupEntityError[T interface {
	fmt.Stringer
	comparable
}] struct {
	propertiesEntities map[T][]*entity
}

func (err dupEntityError[T]) Error() string {
	output := make([]string, 0, len(err.propertiesEntities))

	for property, entities := range err.propertiesEntities {
		errEntities := make([]string, 0, len(entities))
		for _, entity := range entities {
			kind := (*entity).kind()
			name := (*entity).stringerName()

			errEntities = append(errEntities, fmt.Sprintf("%s: %s", kind, name))
		}
		output = append(output, fmt.Sprintf("%s: %s", property, strings.Join(errEntities, ", ")))
	}

	return strings.Join(output, "\n")
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
		return dupEntityError[T]{propertiesEntities: board}
	}
	return nil
}
