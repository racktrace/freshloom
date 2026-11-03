package step

import (
	"fmt"
	"strings"

	"github.com/racktrace/freshloom/internal/system"
)

// Apt installs Debian/Ubuntu packages using apt-get.
type Apt struct {
	Packages []string
	Exec     system.Executor
}

func (a *Apt) Name() string {
	return fmt.Sprintf("apt: %s", strings.Join(a.Packages, ", "))
}

func (a *Apt) Check() (bool, error) {
	missing, err := a.missing()
	if err != nil {
		return false, err
	}
	return len(missing) == 0, nil
}

func (a *Apt) Apply() error {
	missing, err := a.missing()
	if err != nil {
		return err
	}
	if len(missing) == 0 {
		return nil
	}

	if out, err := a.Exec.Run("sudo", "DEBIAN_FRONTEND=noninteractive", "apt-get", "update"); err != nil {
		return fmt.Errorf("update package index: %w\n%s", err, out)
	}

	args := append([]string{"DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y"}, missing...)
	if out, err := a.Exec.Run("sudo", args...); err != nil {
		return fmt.Errorf("install %s: %w\n%s", strings.Join(missing, " "), err, out)
	}
	return nil
}

// missing returns the packages that are not installed yet.
func (a *Apt) missing() ([]string, error) {
	var missing []string
	for _, pkg := range a.Packages {
		out, err := a.Exec.Run("dpkg-query", "-W", "-f=${Status}", pkg)
		if err != nil || out != "install ok installed" {
			// dpkg-query exits non-zero when the package is unknown,
			// which also means it is not installed.
			missing = append(missing, pkg)
		}
	}
	return missing, nil
}