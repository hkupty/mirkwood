package command

// Direction represents the four cardinal directions for movement
type Direction uint8

const (
	North Direction = iota
	South
	East
	West
)

type Walk struct {
	Dir Direction
}

type Mark struct {
	// NOTE: At this stage of development mark is intentionally empty.
	// it can bear more data (i.e. color/rune) which might be a resource later,
	// but as of now it is good as it is.
}
