package store

type Session struct {
	ID             string
	UserID         string
	KeyFingerprint string
	RemoteAddr     string
	StartedAt      string
	Status         string
}

type VM struct {
	ID        string
	SessionID string
	StateDir  string
	FCPid     int
	StartedAt string
}
