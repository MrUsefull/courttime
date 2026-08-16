package courttime

// Team defines a team.
// A team can be a member of multiple divisions.
type Team struct {
	// Name is the team's name or unique identifier
	Name string
	// GamesScheduled counts the number of scheduled games.
	GamesScheduled int
}
