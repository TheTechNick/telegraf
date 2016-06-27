package systemd_units

import (
	"os"

	"github.com/coreos/go-systemd/dbus"
	"github.com/influxdata/telegraf"
	// "github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// Systemd struct
type Systemd struct{}

// Description returns the plugin Description
func (s *Systemd) Description() string {
	return "Gathers state of systemd units"
}

var sampleConfig = `
## NOTE: this plugin has no options
`

var unitStates = []string{"active", "activating", "deactivating", "inactive", "failed"}

// SampleConfig returns the plugin SampleConfig
func (s *Systemd) SampleConfig() string {
	return sampleConfig
}

// Gather gets all metric fields and tags and returns any errors it encounters
func (s *Systemd) Gather(acc telegraf.Accumulator) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	conn, err := dbus.NewSystemdConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnits()
	if err != nil {
		return err
	}

	for _, unit := range units {
		for _, state := range unitStates {
			isActive := 0
			if state == unit.ActiveState {
				isActive = 1
			}
			fields := map[string]interface{}{
				"active": isActive,
			}
			tags := map[string]string{"host": hostname, "name": unit.Name, "state": state}

			acc.AddFields("systemd_units", fields, tags)
		}
	}
	return nil
}

func init() {
	inputs.Add("systemd_units", func() telegraf.Input {
		return &Systemd{}
	})
}
