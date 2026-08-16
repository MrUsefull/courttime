package courttime

type Schedule struct {
	// PrimaryDivision maps a team to a single division that is their primary.
	// This is the division they should usually play in.
	PrimaryDivision map[*Team]*Division
	// SecondaryDivisions maps a team to a collection secondary divisions. These
	// are divisions a team can play "in" but are not the main division.
	// For example, a C division team is primarily a C team. But a C team can play
	// vs a D team or a B team.
	SecondaryDivision map[*Team][]*Division

	Courts []Court
}
