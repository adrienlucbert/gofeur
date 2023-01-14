<p align="center">
  <a>
    <img alt="GOFEUR Logo" src="./gofeur_logo.png" style="width:500px;"/>
  </a>
  <p align="center">
    Projet de gestion d'entrepôt du module de GO d'Epitech.<br>
  </p>
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
bo build # Compile
./gofeur input_file # See Input file section
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
