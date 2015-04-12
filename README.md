# common_actors
leverage the OMDb API to find common actors in two movies or a list of actors in a particular film or tv series.

# Searching
To search for a single movie or tv series:

`go run main.go -s "Enter the Dragon"`

To specify a year:

`go run main.go -s "Daredevil" -y 2015`

# Find common actors
`go run main.go -c "Enter the Dragon, Game of Death"`
