package uniter_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/juju-core/state"
	"launchpad.net/juju-core/testing"
	"launchpad.net/juju-core/worker/uniter"
	stdtesting "testing"
	"time"
)

func Test(t *stdtesting.T) { testing.ZkTestPackage(t) }

type HookQueueSuite struct{}

var _ = Suite(&HookQueueSuite{})

type msi map[string]int

type hookQueueTest struct {
	initial *uniter.RelationState
	steps   []checker
}

func fullTest(steps ...checker) hookQueueTest {
	return hookQueueTest{&uniter.RelationState{"", 21345, nil, ""}, steps}
}

func reconcileTest(members msi, joined string, steps ...checker) hookQueueTest {
	return hookQueueTest{&uniter.RelationState{"", 21345, members, joined}, steps}
}

var hookQueueTests = []hookQueueTest{
	fullTest(
		// Empty initial change causes no hooks.
		send{nil, nil},
	), fullTest(
		// Joined and changed are both run when unit is first detected.
		send{msi{"u/0": 0}, nil},
		expect{"joined", "u/0", 0, msi{"u/0": 0}},
		expect{"changed", "u/0", 0, msi{"u/0": 0}},
	), fullTest(
		// Automatic changed is run with latest settings.
		send{msi{"u/0": 0}, nil},
		expect{"joined", "u/0", 0, msi{"u/0": 0}},
		send{msi{"u/0": 7}, nil},
		expect{"changed", "u/0", 7, msi{"u/0": 7}},
	), fullTest(
		// Joined is also run with latest settings.
		send{msi{"u/0": 0}, nil},
		send{msi{"u/0": 7}, nil},
		expect{"joined", "u/0", 7, msi{"u/0": 7}},
		expect{"changed", "u/0", 7, msi{"u/0": 7}},
	), fullTest(
		// Nothing happens if a unit departs before its joined is run.
		send{msi{"u/0": 0}, nil},
		send{msi{"u/0": 7}, nil},
		send{nil, []string{"u/0"}},
	), fullTest(
		// A changed is run after a joined, even if a departed is known.
		send{msi{"u/0": 0}, nil},
		expect{"joined", "u/0", 0, msi{"u/0": 0}},
		send{nil, []string{"u/0"}},
		expect{"changed", "u/0", 0, msi{"u/0": 0}},
		expect{"departed", "u/0", 0, msi{}},
	), fullTest(
		// A departed replaces a changed.
		send{msi{"u/0": 0}, nil},
		advance{2},
		send{msi{"u/0": 7}, nil},
		send{nil, []string{"u/0"}},
		expect{"departed", "u/0", 7, msi{}},
	), fullTest(
		// Changed events are ignored if the version has not changed.
		send{msi{"u/0": 0}, nil},
		advance{2},
		send{msi{"u/0": 0}, nil},
	), fullTest(
		// Multiple changed events are compacted into one.
		send{msi{"u/0": 0}, nil},
		advance{2},
		send{msi{"u/0": 3}, nil},
		send{msi{"u/0": 7}, nil},
		send{msi{"u/0": 79}, nil},
		expect{"changed", "u/0", 79, msi{"u/0": 79}},
	), fullTest(
		// Multiple changed events are elided.
		send{msi{"u/0": 0}, nil},
		advance{2},
		send{msi{"u/0": 3}, nil},
		send{msi{"u/0": 7}, nil},
		send{msi{"u/0": 79}, nil},
		expect{"changed", "u/0", 79, msi{"u/0": 79}},
	), fullTest(
		// Latest hooks are run in the original unit order.
		send{msi{"u/0": 0, "u/1": 1}, nil},
		advance{4},
		send{msi{"u/0": 3}, nil},
		send{msi{"u/1": 7}, nil},
		send{nil, []string{"u/0"}},
		expect{"departed", "u/0", 3, msi{"u/1": 7}},
		expect{"changed", "u/1", 7, msi{"u/1": 7}},
	), fullTest(
		// Test everything we can think of at the same time.
		send{msi{"u/0": 0, "u/1": 0, "u/2": 0, "u/3": 0, "u/4": 0}, nil},
		advance{6},
		// u/0, u/1, u/2 are now up to date; u/3, u/4 are untouched.
		send{msi{"u/0": 1}, nil},
		send{msi{"u/1": 1, "u/2": 1, "u/3": 1, "u/5": 0}, []string{"u/0", "u/4"}},
		send{msi{"u/3": 2}, nil},
		// - Finish off the rest of the initial state, ignoring u/4, but using
		// the latest known settings.
		expect{"joined", "u/3", 2, msi{"u/0": 1, "u/1": 1, "u/2": 1, "u/3": 2}},
		expect{"changed", "u/3", 2, msi{"u/0": 1, "u/1": 1, "u/2": 1, "u/3": 2}},
		// - u/0 was queued for change by the first RUC, but this change is
		// no longer relevant; it's departed in the second RUC, so we run
		// that hook instead.
		expect{"departed", "u/0", 1, msi{"u/1": 1, "u/2": 1, "u/3": 2}},
		// - Handle the remaining changes in the second RUC, still ignoring u/4.
		// We do run new changed hooks for u/1 and u/2, because the latest settings
		// are newer than those used in their original changed events.
		expect{"changed", "u/1", 1, msi{"u/1": 1, "u/2": 1, "u/3": 2}},
		expect{"changed", "u/2", 1, msi{"u/1": 1, "u/2": 1, "u/3": 2}},
		expect{"joined", "u/5", 0, msi{"u/1": 1, "u/2": 1, "u/3": 2, "u/5": 0}},
		expect{"changed", "u/5", 0, msi{"u/1": 1, "u/2": 1, "u/3": 2, "u/5": 0}},
		// - Ignore the third RUC, because the original joined/changed on u/3
		// was executed after we got the latest settings version.
	), reconcileTest(
		// Check that matching settings versions cause no changes.
		msi{"u/0": 0}, "",
		send{msi{"u/0": 0}, nil},
	), reconcileTest(
		// Check that new settings versions cause appropriate changes.
		msi{"u/0": 0}, "",
		send{msi{"u/0": 1}, nil},
		expect{"changed", "u/0", 1, msi{"u/0": 1}},
	), reconcileTest(
		// Check that a just-joined unit gets its changed hook run first.
		msi{"u/0": 0}, "u/0",
		send{msi{"u/0": 0}, nil},
		expect{"changed", "u/0", 0, msi{"u/0": 0}},
	), reconcileTest(
		// Check that missing units are queued for depart as early as possible.
		msi{"u/0": 0}, "",
		send{msi{"u/1": 0}, nil},
		expect{"departed", "u/0", 0, msi{}},
		expect{"joined", "u/1", 0, msi{"u/1": 0}},
		expect{"changed", "u/1", 0, msi{"u/1": 0}},
	), reconcileTest(
		// Double-check that a just-joined unit gets its changed hook run first,
		// even when it's due to depart.
		msi{"u/0": 0}, "u/0",
		send{nil, nil},
		expect{"changed", "u/0", 0, msi{"u/0": -1}},
		expect{"departed", "u/0", 0, msi{}},
	), reconcileTest(
		// Check that missing units don't slip in front of required changed hooks.
		msi{"u/0": 0}, "u/0",
		send{msi{"u/1": 0}, nil},
		expect{"changed", "u/0", 0, msi{"u/0": -1}},
		expect{"departed", "u/0", 0, msi{}},
		expect{"joined", "u/1", 0, msi{"u/1": 0}},
		expect{"changed", "u/1", 0, msi{"u/1": 0}},
	),
}

func (s *HookQueueSuite) TestHookQueue(c *C) {
	for i, t := range hookQueueTests {
		c.Logf("test %d", i)
		out := make(chan uniter.HookInfo)
		in := make(chan state.RelationUnitsChange)
		ruw := &RUW{in, false}
		q := uniter.NewHookQueue(t.initial, out, ruw)
		for i, step := range t.steps {
			c.Logf("  step %d", i)
			step.check(c, in, out)
		}
		expect{}.check(c, in, out)
		q.Stop()
		c.Assert(ruw.stopped, Equals, true)
	}
}

var brokenHookQueueTests = []hookQueueTest{
	fullTest(
		// Empty state just gets a broken hook.
		expect{hook: "broken"},
	), reconcileTest(
		// Each current member is departed before broken is sent.
		msi{"u/1": 7, "u/4": 33}, "",
		expect{"departed", "u/1", 7, msi{"u/4": -1}},
		expect{"departed", "u/4", 33, msi{}},
		expect{hook: "broken"},
	), reconcileTest(
		// If there's a pending changed, that must still be respected.
		msi{"u/0": 3}, "u/0",
		expect{"changed", "u/0", 3, msi{"u/0": -1}},
		expect{"departed", "u/0", 3, msi{}},
		expect{hook: "broken"},
	),
}

func (s *HookQueueSuite) TestBrokenHookQueue(c *C) {
	for i, t := range brokenHookQueueTests {
		c.Logf("test %d", i)
		out := make(chan uniter.HookInfo)
		q := uniter.NewBrokenHookQueue(t.initial, out)
		for i, step := range t.steps {
			c.Logf("  step %d", i)
			step.check(c, nil, out)
		}
		expect{}.check(c, nil, out)
		q.Stop()
	}
}

// RUW exists entirely to send RelationUnitsChanged events to a tested
// HookQueue in a synchronous and predictable fashion.
type RUW struct {
	in      chan state.RelationUnitsChange
	stopped bool
}

func (w *RUW) Changes() <-chan state.RelationUnitsChange {
	return w.in
}

func (w *RUW) Stop() error {
	close(w.in)
	w.stopped = true
	return nil
}

func (w *RUW) Err() error {
	return nil
}

type checker interface {
	check(c *C, in chan state.RelationUnitsChange, out chan uniter.HookInfo)
}

type send struct {
	changed  msi
	departed []string
}

func (d send) check(c *C, in chan state.RelationUnitsChange, out chan uniter.HookInfo) {
	ruc := state.RelationUnitsChange{Changed: map[string]state.UnitSettings{}}
	for name, version := range d.changed {
		ruc.Changed[name] = state.UnitSettings{
			Version:  version,
			Settings: settings(name, version),
		}
	}
	for _, name := range d.departed {
		ruc.Departed = append(ruc.Departed, name)
	}
	in <- ruc
}

type advance struct {
	count int
}

func (d advance) check(c *C, in chan state.RelationUnitsChange, out chan uniter.HookInfo) {
	for i := 0; i < d.count; i++ {
		select {
		case <-out:
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("timed out waiting for event %d", i)
		}
	}
}

type expect struct {
	hook, unit string
	version    int
	members    msi
}

func (d expect) check(c *C, in chan state.RelationUnitsChange, out chan uniter.HookInfo) {
	if d.hook == "" {
		select {
		case unexpected := <-out:
			c.Fatalf("got %#v", unexpected)
		case <-time.After(200 * time.Millisecond):
		}
		return
	}
	expect := uniter.HookInfo{
		RelationId:    21345,
		HookKind:      d.hook,
		RemoteUnit:    d.unit,
		ChangeVersion: d.version,
	}
	if d.members != nil {
		expect.Members = map[string]map[string]interface{}{}
		for name, version := range d.members {
			expect.Members[name] = settings(name, version)
		}
	}
	select {
	case actual := <-out:
		c.Assert(actual, DeepEquals, expect)
	case <-time.After(200 * time.Millisecond):
		c.Fatalf("timed out waiting for %#v", expect)
	}
}

func settings(name string, version int) map[string]interface{} {
	if version == -1 {
		// Accommodate required events for units no longer present in the
		// relation, whose settings will not be available through the stream
		// of RelationUnitsChanged events.
		return nil
	}
	return map[string]interface{}{
		"unit-name":        name,
		"settings-version": version,
	}
}
