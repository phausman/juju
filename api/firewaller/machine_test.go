// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package firewaller_test

import (
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/firewaller"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/watcher/watchertest"
	"github.com/juju/juju/state"
)

const allEndpoints = ""

type machineSuite struct {
	firewallerSuite

	apiMachine *firewaller.Machine
}

var _ = gc.Suite(&machineSuite{})

func (s *machineSuite) SetUpTest(c *gc.C) {
	s.firewallerSuite.SetUpTest(c)

	var err error
	s.apiMachine, err = s.firewaller.Machine(s.machines[0].Tag().(names.MachineTag))
	c.Assert(err, jc.ErrorIsNil)
}

func (s *machineSuite) TearDownTest(c *gc.C) {
	s.firewallerSuite.TearDownTest(c)
}

func (s *machineSuite) assertRemoveMachinePortsDoc(c *gc.C, p state.MachinePortRanges) {
	type remover interface {
		Remove() error
	}

	portRemover, supports := p.(remover)
	c.Assert(supports, jc.IsTrue, gc.Commentf("machine ports interface does not implement Remove()"))
	c.Assert(portRemover.Remove(), jc.ErrorIsNil)
}

func (s *machineSuite) TestMachine(c *gc.C) {
	apiMachine42, err := s.firewaller.Machine(names.NewMachineTag("42"))
	c.Assert(err, gc.ErrorMatches, "machine 42 not found")
	c.Assert(err, jc.Satisfies, params.IsCodeNotFound)
	c.Assert(apiMachine42, gc.IsNil)

	apiMachine0, err := s.firewaller.Machine(s.machines[0].Tag().(names.MachineTag))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(apiMachine0, gc.NotNil)
}

func (s *machineSuite) TestTag(c *gc.C) {
	c.Assert(s.apiMachine.Tag(), gc.Equals, names.NewMachineTag(s.machines[0].Id()))
}

func (s *machineSuite) TestInstanceId(c *gc.C) {
	// Add another, not provisioned machine to test
	// CodeNotProvisioned.
	newMachine, err := s.State.AddMachine("quantal", state.JobHostUnits)
	c.Assert(err, jc.ErrorIsNil)
	apiNewMachine, err := s.firewaller.Machine(newMachine.Tag().(names.MachineTag))
	c.Assert(err, jc.ErrorIsNil)
	_, err = apiNewMachine.InstanceId()
	c.Assert(err, gc.ErrorMatches, "machine 3 not provisioned")
	c.Assert(err, jc.Satisfies, params.IsCodeNotProvisioned)

	instanceId, err := s.apiMachine.InstanceId()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(instanceId, gc.Equals, instance.Id("i-manager"))
}

func (s *machineSuite) TestWatchUnits(c *gc.C) {
	w, err := s.apiMachine.WatchUnits()
	c.Assert(err, jc.ErrorIsNil)
	wc := watchertest.NewStringsWatcherC(c, w, s.BackingState.StartSync)
	defer wc.AssertStops()

	// Initial event.
	wc.AssertChange("wordpress/0")
	wc.AssertNoChange()

	// Change something other than the life cycle and make sure it's
	// not detected.
	err = s.machines[0].SetPassword("foo")
	c.Assert(err, gc.ErrorMatches, "password is only 3 bytes long, and is not a valid Agent password")
	wc.AssertNoChange()

	err = s.machines[0].SetPassword("foo-12345678901234567890")
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertNoChange()

	// Unassign unit 0 from the machine and check it's detected.
	err = s.units[0].UnassignFromMachine()
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertChange("wordpress/0")
	wc.AssertNoChange()
}

func (s *machineSuite) TestActiveSubnets(c *gc.C) {
	// No ports opened at first, no active subnets.
	subnets, err := s.apiMachine.ActiveSubnets()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(subnets, gc.HasLen, 0)

	// Open a port and check again.
	mustOpenPortRanges(c, s.State, s.units[0], allEndpoints, []network.PortRange{
		network.MustParsePortRange("1234/tcp"),
	})
	subnets, err = s.apiMachine.ActiveSubnets()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(subnets, jc.DeepEquals, []names.SubnetTag{{}})

	// Remove all ports, no more active subnets.
	machPortRanges, err := s.machines[0].OpenedPortRanges()
	c.Assert(err, jc.ErrorIsNil)
	s.assertRemoveMachinePortsDoc(c, machPortRanges)
	subnets, err = s.apiMachine.ActiveSubnets()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(subnets, gc.HasLen, 0)
}

func (s *machineSuite) TestOpenedPorts(c *gc.C) {
	unitTag := s.units[0].Tag().(names.UnitTag)

	// No ports opened at first.
	machPortRanges, err := s.machines[0].OpenedPortRanges()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machPortRanges.UniquePortRanges(), gc.HasLen, 0)

	// Open a port and check again.
	mustOpenPortRanges(c, s.State, s.units[0], allEndpoints, []network.PortRange{
		network.MustParsePortRange("1234/tcp"),
	})
	ports, err := s.apiMachine.OpenedPorts(names.SubnetTag{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(ports, jc.DeepEquals, map[network.PortRange]names.UnitTag{
		{FromPort: 1234, ToPort: 1234, Protocol: "tcp"}: unitTag,
	})
}

func (s *machineSuite) TestIsManual(c *gc.C) {
	answer, err := s.machines[0].IsManual()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(answer, jc.IsFalse)

	m, err := s.State.AddOneMachine(state.MachineTemplate{
		Series:     "quantal",
		Jobs:       []state.MachineJob{state.JobHostUnits},
		InstanceId: "2",
		Nonce:      "manual:",
	})
	c.Assert(err, jc.ErrorIsNil)
	answer, err = m.IsManual()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(answer, jc.IsTrue)

}

func mustOpenPortRanges(c *gc.C, st *state.State, u *state.Unit, endpointName string, portRanges []network.PortRange) {
	unitPortRanges, err := u.OpenedPortRanges()
	c.Assert(err, jc.ErrorIsNil)

	for _, pr := range portRanges {
		unitPortRanges.Open(endpointName, pr)
	}

	c.Assert(st.ApplyOperation(unitPortRanges.Changes()), jc.ErrorIsNil)
}

func mustClosePortRanges(c *gc.C, st *state.State, u *state.Unit, endpointName string, portRanges []network.PortRange) {
	unitPortRanges, err := u.OpenedPortRanges()
	c.Assert(err, jc.ErrorIsNil)

	for _, pr := range portRanges {
		unitPortRanges.Close(endpointName, pr)
	}

	c.Assert(st.ApplyOperation(unitPortRanges.Changes()), jc.ErrorIsNil)
}
