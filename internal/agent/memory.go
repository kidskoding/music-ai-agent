package agent

type SessionMemory struct {
	LastTracks     []string
	SkipHistory    map[string]bool
	EnergyHistory  []float64
	CurrentMode    string
}