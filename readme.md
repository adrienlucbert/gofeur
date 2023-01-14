<p align="center">
  <a>
    <img alt="GOFEUR Logo" src="./gofeur_logo.png" style="width:500px;"/>
  </a>
</p>

Yet another Epitech project but in [Go](https://never-again.go).

Gofeur is a simulation program which goal is to optimize forklifts inside a
warehouse in order to maximize trucks loading.

The warehouse is a two dimensionals grid in which parcel, forklift, and truck
entities live.

## How to run

### Prerequisites

- Go v.19

### Run the project
```bash
go get # Fetch the dependencies
go build # Compile
./gofeur -filename ./input_file # Run gofeur (See Input file section for the file format)
```

### Launch tests
```bash
go test
```

### Input file 

Gofeur runs a simulation from a text input file. It has 4 sections:
- Warehouse:

  A line with three unsigned integer representing respectivelly:
  - The warehouse's width: a positive integer
  - The warehouse's length: a positive integer
  - Simulation's cycle number: an integer in the following range [10, 100_000]
- Parcel
  
  The N following lines might be parcels if they fit the parcel format. The
  parcel format is composed of 4 tokens separated by a space. In order the
  tokens are:
    - Parcel's name: a string without space character
    - Parcel's x coordonate: An unsigned integer
    - Parcel's y coordonate: An unsigned integer
    - Parcel's weight: One of: `yellow`, `green` or `blue`. Weight possible
      values can be in lower or in upper case.
      Each color reprensents a specific weight as so:
        - `yellow`: 100
        - `green`: 200
        - `blue`: 500

- Forklift:

  The N following lines should be forklifts. The forklift format is composed of
  3 tokens separated by a space. In order the tokens are:
    - Forklift's name: a string without space character
    - Forklift's x coordonate: An unsigned integer
    - Forklift's y coordonate: An unsigned integer

- Truck:

  The N following lines should be trucks. The truck format is composed of
  5 tokens separated by a space. In order the tokens are:
    - Truck's name: a string without space character
    - Truck's x coordonate: An unsigned integer
    - Truck's y coordonate: An unsigned integer
    - Truck's maximum weight: An unsigned integer
    - Truck's delivery cycle: An unsigned integer describing the number of
      cycles it takes for the truck once it lefts for delivery to come back.

For instanve a valid input file could be:
```
10 10 15
parcel_a 0 0 yellow
forklift_a 1 0
truck_a 1 9
```

Section are separated by nothing, empty lines aren't allowed and there is no
way to put comment. Good luck with that!

## Code overview

The project is composed of multiple packages, each serving a
distinct purpose.

**Packages**:
- `board`:

  The `board` package is used to represent the warehouse grid with
  the parcels, forklifts, and trucks. This representation is then
  used by the `pathfinding` package.

- `config`

  The `config` package provides a way to store and retrieve the application
  configuration.

- `logger`

  The `logger` provides a set of functions to print messages at different
  levels.

- `optional`

  The `optional` package provides a generic wrapper around a pointer
  to avoid raw pointer manipulation on optional/nullable types. 
 

- `parsing`

  The `parsing` package provides functions to parse an input file into a
  simulation and to check that simulation validity.

- `pathfinding`
  The `pathfinding` package provides functions to find the path between two
  points on a board.
  
- `pkg`
  The `pkg` package provides the common type `Vector`.

- `simulation`
  The `simulation` package is the heart of the project and is responsible of
  running the simulation.
  
- `ui`
  The `ui` package contains utilities to display a TUI interface for the Gofeur
  application.
