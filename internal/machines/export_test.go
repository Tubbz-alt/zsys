package machines

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/ubuntu/zsys/internal/config"
	"github.com/ubuntu/zsys/internal/testutils"
)

const (
	RevertUserDataTag       = zfsRevertUserDataTag
	AutomatedSnapshotPrefix = automatedSnapshotPrefix
)

// WithTime allows overriding default time implementations with a mock
func WithTime(time Nower) func(o *options) error {
	return func(o *options) error {
		o.time = time
		return nil
	}
}

// Import from json to export the private fields
func (ms *Machines) UnmarshalJSON(b []byte) error {
	mt := Machinesdump{}

	if err := json.Unmarshal(b, &mt); err != nil {
		return err
	}

	ms.all = mt.All
	ms.cmdline = mt.Cmdline
	ms.current = mt.Current
	ms.nextState = mt.NextState
	ms.allSystemDatasets = mt.AllSystemDatasets
	ms.allUsersDatasets = mt.AllUsersDatasets
	ms.allPersistentDatasets = mt.AllPersistentDatasets
	ms.unmanagedDatasets = mt.UnmanagedDatasets

	if ms.current != nil {
		for k, machine := range mt.All {
			if machine.ID != ms.current.ID {
				continue
			}
			// restore current machine pointing to the same element than the hashmap
			ms.current = mt.All[k]
		}
	}

	return nil
}

// MakeComparable prepares Machines by resetting private fields that change at each invocation
func (ms *Machines) MakeComparable() {
	ds := sortedDatasets(ms.allSystemDatasets)
	sort.Sort(ds)
	ms.allSystemDatasets = ds

	ds = sortedDatasets(ms.allUsersDatasets)
	sort.Sort(ds)
	ms.allUsersDatasets = ds

	ds = sortedDatasets(ms.allPersistentDatasets)
	sort.Sort(ds)
	ms.allPersistentDatasets = ds

	ds = sortedDatasets(ms.unmanagedDatasets)
	sort.Sort(ds)
	ms.unmanagedDatasets = ds

	ms.z = nil
	ms.time = nil
	ms.conf = config.ZConfig{}
}

// SplitSnapshotName calls internal splitSnapshotName to split a snapshot name in base and id of a snapshot
func SplitSnapshotName(s string) (string, string) { return splitSnapshotName(s) }

// AllMachines exports machines lists for tests
func (ms *Machines) AllMachines() map[string]*Machine { return ms.all }

func (ms Machines) CopyForTests(t *testing.T) (copy Machines) {
	t.Helper()

	testutils.Deepcopy(t, &copy, ms)
	return copy
}
