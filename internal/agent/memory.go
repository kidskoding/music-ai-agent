package agent

type SessionMemory struct {
	LastTracks     []*Track
	SkipHistory    map[string]bool
	EnergyHistory  []float64
	CurrentMode    string
}