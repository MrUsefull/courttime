package courttime

// Division defines a division.
// A division is a collection of teams that can play each other
// Teams can be in multiple divisions.
type Division struct {
	// Name is the unique name of this division
	Name string
	// Teams are all of the teams in this division
	Teams []*Team
}
