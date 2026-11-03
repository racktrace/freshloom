package step

// Step is a single unit of work in freshloom.
// Every step must be able to inspect the current state before changing anything.
type Step interface {
	// Name is a short label used in logs and the summary.
	Name() string

	// Check reports whether the system is already in the desired state.
	Check() (done bool, err error)

	// Apply makes the changes needed to reach the desired state.
	Apply() error
}