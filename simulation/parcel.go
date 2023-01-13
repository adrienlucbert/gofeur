package simulation

import (
	"github.com/adrienlucbert/gofeur/parsing"
	"github.com/adrienlucbert/gofeur/pkg"
)

type parcel struct {
	name    string
	pos     pkg.Vector
	weight  uint
	carried bool
}

func newParcelFromParsing(from *parsing.Parcel) parcel {
	return parcel{
		name: from.Name,
		pos:  pkg.Vector{X: int(from.X), Y: int(from.Y)},
		weight: map[parsing.Color]uint{
			parsing.Yellow: 100,
			parsing.Green:  200,
			parsing.Blue:   500,
		}[from.Color],
	}
}
