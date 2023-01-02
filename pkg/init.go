package pkg

// Temporary (Wainting for the parser to be done)
func InitStorageBuilding() *StorageBuilding {
	t1 := Truck{"f", 3, 10, 12, 199}
	truck := []Truck{t1}
	trans := []Transpals{{"a", 3, 10}, {"d", 3, 10}}
	parcels := []Parcel{
		{"ch", 2, 2, GREEN},
		{"ca", 3, 5, YELLOW},
	}
	sb := StorageBuilding{trans, parcels, truck}
	return &sb
}
