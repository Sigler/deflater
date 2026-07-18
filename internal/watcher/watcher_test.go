package watcher

import (
	"reflect"
	"testing"
)

func set(names ...string) map[string]bool {
	m := map[string]bool{}
	for _, n := range names {
		m[n] = true
	}
	return m
}

func TestFirstRunNeverAlerts(t *testing.T) {
	got := NewArrivals(set("LG.Junk"), nil, nil)
	if got != nil {
		t.Fatalf("nil snapshot means no baseline; got %v", got)
	}
}

func TestNewArrivalIsFlagged(t *testing.T) {
	got := NewArrivals(set("Old.App", "LG.ThinQ"), []string{"Old.App"}, nil)
	if !reflect.DeepEqual(got, []string{"LG.ThinQ"}) {
		t.Fatalf("got %v", got)
	}
}

func TestRoutineSystemPackagesAreIgnored(t *testing.T) {
	current := set("Old.App", "Microsoft.NET.Native.Framework.2.2", "Microsoft.VCLibs.140.00", "Microsoft.Windows.Something")
	got := NewArrivals(current, []string{"Old.App"}, nil)
	if got != nil {
		t.Fatalf("servicing packages must not alert; got %v", got)
	}
}

func TestManagedPackagesAreIgnored(t *testing.T) {
	// A package Deflater itself manages coming back is maintenance's
	// job to remove, not the watcher's job to announce.
	got := NewArrivals(set("Old.App", "Microsoft.BingNews"), []string{"Old.App"}, set("Microsoft.BingNews"))
	if got != nil {
		t.Fatalf("managed packages must not alert; got %v", got)
	}
}

func TestArrivalsAreSorted(t *testing.T) {
	got := NewArrivals(set("Zebra.App", "Alpha.App"), []string{}, nil)
	if !reflect.DeepEqual(got, []string{"Alpha.App", "Zebra.App"}) {
		t.Fatalf("got %v", got)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	snap := SnapshotOf(set("B.App", "A.App"))
	if !reflect.DeepEqual(snap, []string{"A.App", "B.App"}) {
		t.Fatalf("snapshot should be sorted: %v", snap)
	}
	if got := NewArrivals(set("B.App", "A.App"), snap, nil); got != nil {
		t.Fatalf("no changes means no arrivals; got %v", got)
	}
}
