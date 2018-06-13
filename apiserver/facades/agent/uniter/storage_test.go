// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter_test

import (
	"errors"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v2"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/facades/agent/uniter"
	"github.com/juju/juju/apiserver/params"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/state"
	statetesting "github.com/juju/juju/state/testing"
	"github.com/juju/juju/testing"
)

type storageSuite struct {
	testing.BaseSuite
	called []string
}

var _ = gc.Suite(&storageSuite{})

func (s *storageSuite) TestWatchUnitStorageAttachments(c *gc.C) {
	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(names.Tag) bool {
			return true
		}, nil
	}
	unitTag := names.NewUnitTag("mysql/0")
	watcher := &mockStringsWatcher{
		changes: make(chan []string, 1),
	}
	watcher.changes <- []string{"storage/0", "storage/1"}
	state := &mockStorageState{
		watchStorageAttachments: func(u names.UnitTag) state.StringsWatcher {
			c.Assert(u, gc.DeepEquals, unitTag)
			return watcher
		},
	}

	storage, err := uniter.NewStorageAPI(state, state, state, state, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	watches, err := storage.WatchUnitStorageAttachments(params.Entities{
		Entities: []params.Entity{{unitTag.String()}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(watches, gc.DeepEquals, params.StringsWatchResults{
		Results: []params.StringsWatchResult{{
			StringsWatcherId: "1",
			Changes:          []string{"storage/0", "storage/1"},
		}},
	})
	c.Assert(resources.Get("1"), gc.Equals, watcher)
}

func (s *storageSuite) TestWatchStorageAttachmentVolume(c *gc.C) {
	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(names.Tag) bool {
			return true
		}, nil
	}
	unitTag := names.NewUnitTag("mysql/0")
	storageTag := names.NewStorageTag("data/0")
	machineTag := names.NewMachineTag("66")
	volumeTag := names.NewVolumeTag("104")
	volume := &mockVolume{tag: volumeTag}
	storageInstance := &mockStorageInstance{kind: state.StorageKindBlock}
	storageWatcher := &mockNotifyWatcher{
		changes: make(chan struct{}, 1),
	}
	storageWatcher.changes <- struct{}{}
	volumeWatcher := &mockNotifyWatcher{
		changes: make(chan struct{}, 1),
	}
	volumeWatcher.changes <- struct{}{}
	blockDevicesWatcher := &mockNotifyWatcher{
		changes: make(chan struct{}, 1),
	}
	blockDevicesWatcher.changes <- struct{}{}
	var calls []string
	state := &mockStorageState{
		assignedMachine: "66",
		storageInstance: func(s names.StorageTag) (state.StorageInstance, error) {
			calls = append(calls, "StorageInstance")
			c.Assert(s, gc.DeepEquals, storageTag)
			return storageInstance, nil
		},
		storageInstanceVolume: func(s names.StorageTag) (state.Volume, error) {
			calls = append(calls, "StorageInstanceVolume")
			c.Assert(s, gc.DeepEquals, storageTag)
			return volume, nil
		},
		watchStorageAttachment: func(s names.StorageTag, u names.UnitTag) state.NotifyWatcher {
			calls = append(calls, "WatchStorageAttachment")
			c.Assert(s, gc.DeepEquals, storageTag)
			c.Assert(u, gc.DeepEquals, unitTag)
			return storageWatcher
		},
		watchVolumeAttachment: func(m names.MachineTag, v names.VolumeTag) state.NotifyWatcher {
			calls = append(calls, "WatchVolumeAttachment")
			c.Assert(m, gc.DeepEquals, machineTag)
			c.Assert(v, gc.DeepEquals, volumeTag)
			return volumeWatcher
		},
		watchBlockDevices: func(m names.MachineTag) state.NotifyWatcher {
			calls = append(calls, "WatchBlockDevices")
			c.Assert(m, gc.DeepEquals, machineTag)
			return blockDevicesWatcher
		},
	}

	storage, err := uniter.NewStorageAPI(state, state, state, state, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	watches, err := storage.WatchStorageAttachments(params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag.String(),
			UnitTag:    unitTag.String(),
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(watches, gc.DeepEquals, params.NotifyWatchResults{
		Results: []params.NotifyWatchResult{{
			NotifyWatcherId: "1",
		}},
	})
	c.Assert(calls, gc.DeepEquals, []string{
		"StorageInstance",
		"StorageInstanceVolume",
		"WatchVolumeAttachment",
		"WatchBlockDevices",
		"WatchStorageAttachment",
	})
}

func (s *storageSuite) TestWatchStorageAttachmentFilesystem(c *gc.C) {
	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(names.Tag) bool {
			return true
		}, nil
	}
	unitTag := names.NewUnitTag("mysql/0")
	storageTag := names.NewStorageTag("data/0")
	machineTag := names.NewMachineTag("66")
	filesystemTag := names.NewFilesystemTag("104")
	filesystem := &mockFilesystem{tag: filesystemTag}
	storageInstance := &mockStorageInstance{kind: state.StorageKindFilesystem}
	storageWatcher := &mockNotifyWatcher{
		changes: make(chan struct{}, 1),
	}
	storageWatcher.changes <- struct{}{}
	filesystemWatcher := &mockNotifyWatcher{
		changes: make(chan struct{}, 1),
	}
	filesystemWatcher.changes <- struct{}{}
	var calls []string
	state := &mockStorageState{
		assignedMachine: "66",
		storageInstance: func(s names.StorageTag) (state.StorageInstance, error) {
			calls = append(calls, "StorageInstance")
			c.Assert(s, gc.DeepEquals, storageTag)
			return storageInstance, nil
		},
		storageInstanceFilesystem: func(s names.StorageTag) (state.Filesystem, error) {
			calls = append(calls, "StorageInstanceFilesystem")
			c.Assert(s, gc.DeepEquals, storageTag)
			return filesystem, nil
		},
		watchStorageAttachment: func(s names.StorageTag, u names.UnitTag) state.NotifyWatcher {
			calls = append(calls, "WatchStorageAttachment")
			c.Assert(s, gc.DeepEquals, storageTag)
			c.Assert(u, gc.DeepEquals, unitTag)
			return storageWatcher
		},
		watchFilesystemAttachment: func(m names.MachineTag, f names.FilesystemTag) state.NotifyWatcher {
			calls = append(calls, "WatchFilesystemAttachment")
			c.Assert(m, gc.DeepEquals, machineTag)
			c.Assert(f, gc.DeepEquals, filesystemTag)
			return filesystemWatcher
		},
	}

	storage, err := uniter.NewStorageAPI(state, state, state, state, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	watches, err := storage.WatchStorageAttachments(params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag.String(),
			UnitTag:    unitTag.String(),
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(watches, gc.DeepEquals, params.NotifyWatchResults{
		Results: []params.NotifyWatchResult{{
			NotifyWatcherId: "1",
		}},
	})
	c.Assert(calls, gc.DeepEquals, []string{
		"StorageInstance",
		"StorageInstanceFilesystem",
		"WatchFilesystemAttachment",
		"WatchStorageAttachment",
	})
}

func (s *storageSuite) TestDestroyUnitStorageAttachments(c *gc.C) {
	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(names.Tag) bool {
			return true
		}, nil
	}
	unitTag := names.NewUnitTag("mysql/0")
	var calls []string
	state := &mockStorageState{
		destroyUnitStorageAttachments: func(u names.UnitTag) error {
			calls = append(calls, "DestroyUnitStorageAttachments")
			c.Assert(u, gc.DeepEquals, unitTag)
			return nil
		},
	}

	storage, err := uniter.NewStorageAPI(state, state, state, state, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	errors, err := storage.DestroyUnitStorageAttachments(params.Entities{
		Entities: []params.Entity{{
			Tag: unitTag.String(),
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(calls, jc.DeepEquals, []string{"DestroyUnitStorageAttachments"})
	c.Assert(errors, jc.DeepEquals, params.ErrorResults{
		[]params.ErrorResult{{}},
	})
}

func (s *storageSuite) TestRemoveStorageAttachments(c *gc.C) {
	setMock := func(st *mockStorageState, f func(s names.StorageTag, u names.UnitTag) error) {
		st.remove = f
	}

	unitTag0 := names.NewUnitTag("mysql/0")
	unitTag1 := names.NewUnitTag("mysql/1")
	storageTag0 := names.NewStorageTag("data/0")
	storageTag1 := names.NewStorageTag("data/1")

	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			return tag == unitTag0
		}, nil
	}

	state := &mockStorageState{}
	setMock(state, func(s names.StorageTag, u names.UnitTag) error {
		c.Assert(u, gc.DeepEquals, unitTag0)
		if s == storageTag1 {
			return errors.New("badness")
		}
		return nil
	})

	storage, err := uniter.NewStorageAPI(state, state, state, state, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	errors, err := storage.RemoveStorageAttachments(params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag0.String(),
			UnitTag:    unitTag0.String(),
		}, {
			StorageTag: storageTag1.String(),
			UnitTag:    unitTag0.String(),
		}, {
			StorageTag: storageTag0.String(),
			UnitTag:    unitTag1.String(),
		}, {
			StorageTag: unitTag0.String(), // oops
			UnitTag:    unitTag0.String(),
		}, {
			StorageTag: storageTag0.String(),
			UnitTag:    storageTag0.String(), // oops
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(errors, jc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{
			{nil},
			{&params.Error{Message: "badness"}},
			{&params.Error{Code: params.CodeUnauthorized, Message: "permission denied"}},
			{&params.Error{Message: `"unit-mysql-0" is not a valid storage tag`}},
			{&params.Error{Message: `"storage-data-0" is not a valid unit tag`}},
		},
	})
}

const (
	addStorageCall = "mockAdd"
)

func (s *storageSuite) TestAddUnitStorageConstraintsErrors(c *gc.C) {
	unitTag0 := names.NewUnitTag("mysql/0")
	storageName0 := "data"
	storageName1 := "store"

	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			return tag == unitTag0
		}, nil
	}

	s.called = []string{}
	mockState := &mockStorageState{
		unitStorageConstraints: map[string]state.StorageConstraints{
			storageName0: {},
		},
	}

	storage, err := uniter.NewStorageAPI(mockState, mockState, mockState, mockState, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	size := uint64(10)
	count := uint64(0)
	errors, err := storage.AddUnitStorage(params.StoragesAddParams{
		Storages: []params.StorageAddParams{
			{
				UnitTag:     unitTag0.String(),
				StorageName: storageName0,
				Constraints: params.StorageConstraints{Pool: "matter"},
			}, {
				UnitTag:     unitTag0.String(),
				StorageName: storageName0,
				Constraints: params.StorageConstraints{Size: &size},
			}, {
				UnitTag:     unitTag0.String(),
				StorageName: storageName0,
				Constraints: params.StorageConstraints{},
			}, {
				UnitTag:     unitTag0.String(),
				StorageName: storageName0,
				Constraints: params.StorageConstraints{Count: &count},
			}, {
				UnitTag:     unitTag0.String(),
				StorageName: storageName1,
				Constraints: params.StorageConstraints{},
			},
		}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.called, gc.HasLen, 0)
	c.Assert(errors, jc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{
			{&params.Error{Message: `adding storage data for unit-mysql-0: only count can be specified`}},
			{&params.Error{Message: `adding storage data for unit-mysql-0: only count can be specified`}},
			{&params.Error{Message: `adding storage data for unit-mysql-0: count must be specified`}},
			{&params.Error{Message: `adding storage data for unit-mysql-0: count must be specified`}},
			{&params.Error{
				Code:    "not found",
				Message: "adding storage store for unit-mysql-0: storage \"store\" not found"}},
		},
	})
}

func (s *storageSuite) TestAddUnitStorage(c *gc.C) {
	setMockAdd := func(st *mockStorageState, f func(tag names.UnitTag, name string, cons state.StorageConstraints) error) {
		st.addUnitStorage = f
	}

	unitTag0 := names.NewUnitTag("mysql/0")
	storageName0 := "data"
	storageName1 := "store"

	unitPool := "real"
	size := uint64(3)
	unitSize := size * 2
	unitCount := uint64(100)
	testCount := uint64(10)

	resources := common.NewResources()
	getCanAccess := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			return tag == unitTag0
		}, nil
	}

	s.called = []string{}
	mockState := &mockStorageState{
		unitStorageConstraints: map[string]state.StorageConstraints{
			storageName0: {
				Pool:  unitPool,
				Size:  unitSize,
				Count: unitCount,
			},
			storageName1: {},
		},
	}

	setMockAdd(mockState, func(u names.UnitTag, name string, cons state.StorageConstraints) error {
		s.called = append(s.called, addStorageCall)
		c.Assert(u, gc.DeepEquals, unitTag0)
		if name == storageName1 {
			return errors.New("badness")
		}
		c.Assert(cons.Count, gc.Not(gc.Equals), unitCount)
		c.Assert(cons.Count, jc.DeepEquals, testCount)
		c.Assert(cons.Pool, jc.DeepEquals, unitPool)
		c.Assert(cons.Size, jc.DeepEquals, unitSize)
		return nil
	})

	storage, err := uniter.NewStorageAPI(mockState, mockState, mockState, mockState, resources, getCanAccess)
	c.Assert(err, jc.ErrorIsNil)
	errors, err := storage.AddUnitStorage(params.StoragesAddParams{
		Storages: []params.StorageAddParams{
			{
				UnitTag:     unitTag0.String(),
				StorageName: storageName0,
				Constraints: params.StorageConstraints{Count: &testCount},
			}, {
				UnitTag:     unitTag0.String(),
				StorageName: storageName1,
				Constraints: params.StorageConstraints{Count: &testCount},
			},
		}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.called, jc.SameContents, []string{
		addStorageCall,
		addStorageCall,
	})
	c.Assert(errors, jc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{
			{nil},
			{&params.Error{Message: "adding storage store for unit-mysql-0: badness"}},
		},
	})
}

type mockUnit struct {
	assignedMachine    string
	storageConstraints map[string]state.StorageConstraints
}

func (u *mockUnit) AssignedMachineId() (string, error) {
	return u.assignedMachine, nil
}

func (u *mockUnit) StorageConstraints() (map[string]state.StorageConstraints, error) {
	return u.storageConstraints, nil
}

type mockStorageState struct {
	unitStorageConstraints map[string]state.StorageConstraints
	assignedMachine        string

	uniter.Backend
	uniter.StorageStateInterface
	uniter.StorageVolumeInterface
	uniter.StorageFilesystemInterface
	destroyUnitStorageAttachments func(names.UnitTag) error
	remove                        func(names.StorageTag, names.UnitTag) error
	storageInstance               func(names.StorageTag) (state.StorageInstance, error)
	storageInstanceFilesystem     func(names.StorageTag) (state.Filesystem, error)
	storageInstanceVolume         func(names.StorageTag) (state.Volume, error)
	watchStorageAttachments       func(names.UnitTag) state.StringsWatcher
	watchStorageAttachment        func(names.StorageTag, names.UnitTag) state.NotifyWatcher
	watchFilesystemAttachment     func(names.MachineTag, names.FilesystemTag) state.NotifyWatcher
	watchVolumeAttachment         func(names.MachineTag, names.VolumeTag) state.NotifyWatcher
	watchBlockDevices             func(names.MachineTag) state.NotifyWatcher
	addUnitStorage                func(u names.UnitTag, name string, cons state.StorageConstraints) error
}

func (m *mockStorageState) Unit(name string) (uniter.Unit, error) {
	return &mockUnit{
		assignedMachine:    m.assignedMachine,
		storageConstraints: m.unitStorageConstraints}, nil
}

func (m *mockStorageState) DestroyUnitStorageAttachments(u names.UnitTag) error {
	return m.destroyUnitStorageAttachments(u)
}

func (m *mockStorageState) RemoveStorageAttachment(s names.StorageTag, u names.UnitTag) error {
	return m.remove(s, u)
}

func (m *mockStorageState) StorageInstance(s names.StorageTag) (state.StorageInstance, error) {
	return m.storageInstance(s)
}

func (m *mockStorageState) StorageInstanceFilesystem(s names.StorageTag) (state.Filesystem, error) {
	return m.storageInstanceFilesystem(s)
}

func (m *mockStorageState) StorageInstanceVolume(s names.StorageTag) (state.Volume, error) {
	return m.storageInstanceVolume(s)
}

func (m *mockStorageState) WatchStorageAttachments(u names.UnitTag) state.StringsWatcher {
	return m.watchStorageAttachments(u)
}

func (m *mockStorageState) WatchStorageAttachment(s names.StorageTag, u names.UnitTag) state.NotifyWatcher {
	return m.watchStorageAttachment(s, u)
}

func (m *mockStorageState) WatchFilesystemAttachment(mtag names.MachineTag, f names.FilesystemTag) state.NotifyWatcher {
	return m.watchFilesystemAttachment(mtag, f)
}

func (m *mockStorageState) WatchVolumeAttachment(mtag names.MachineTag, v names.VolumeTag) state.NotifyWatcher {
	return m.watchVolumeAttachment(mtag, v)
}

func (m *mockStorageState) WatchBlockDevices(mtag names.MachineTag) state.NotifyWatcher {
	return m.watchBlockDevices(mtag)
}

func (m *mockStorageState) AddStorageForUnit(tag names.UnitTag, name string, cons state.StorageConstraints) ([]names.StorageTag, error) {
	return nil, m.addUnitStorage(tag, name, cons)
}

type mockStringsWatcher struct {
	state.StringsWatcher
	changes chan []string
}

func (m *mockStringsWatcher) Changes() <-chan []string {
	return m.changes
}

type mockNotifyWatcher struct {
	tomb    tomb.Tomb
	changes chan struct{}
}

func (m *mockNotifyWatcher) Stop() error {
	m.Kill()
	return m.Wait()
}

func (m *mockNotifyWatcher) Kill() {
	m.tomb.Kill(nil)
}

func (m *mockNotifyWatcher) Wait() error {
	return m.tomb.Wait()
}

func (m *mockNotifyWatcher) Err() error {
	return m.tomb.Err()
}

func (m *mockNotifyWatcher) Changes() <-chan struct{} {
	return m.changes
}

type mockVolume struct {
	state.Volume
	tag names.VolumeTag
}

func (m *mockVolume) VolumeTag() names.VolumeTag {
	return m.tag
}

type mockFilesystem struct {
	state.Filesystem
	tag names.FilesystemTag
}

func (m *mockFilesystem) FilesystemTag() names.FilesystemTag {
	return m.tag
}

type mockStorageInstance struct {
	state.StorageInstance
	kind state.StorageKind
}

func (m *mockStorageInstance) Kind() state.StorageKind {
	return m.kind
}

type watchStorageAttachmentSuite struct {
	storageTag               names.StorageTag
	machineTag               names.MachineTag
	unitTag                  names.UnitTag
	st                       *fakeStorage
	storageInstance          *fakeStorageInstance
	volume                   *fakeVolume
	volumeAttachmentWatcher  *apiservertesting.FakeNotifyWatcher
	blockDevicesWatcher      *apiservertesting.FakeNotifyWatcher
	storageAttachmentWatcher *apiservertesting.FakeNotifyWatcher
}

var _ = gc.Suite(&watchStorageAttachmentSuite{})

func (s *watchStorageAttachmentSuite) SetUpTest(c *gc.C) {
	s.storageTag = names.NewStorageTag("osd-devices/0")
	s.machineTag = names.NewMachineTag("0")
	s.unitTag = names.NewUnitTag("ceph/0")
	s.storageInstance = &fakeStorageInstance{
		tag:   s.storageTag,
		owner: s.machineTag,
		kind:  state.StorageKindBlock,
	}
	s.volume = &fakeVolume{tag: names.NewVolumeTag("0")}
	s.volumeAttachmentWatcher = apiservertesting.NewFakeNotifyWatcher()
	s.blockDevicesWatcher = apiservertesting.NewFakeNotifyWatcher()
	s.storageAttachmentWatcher = apiservertesting.NewFakeNotifyWatcher()
	s.st = &fakeStorage{
		storageInstance: func(tag names.StorageTag) (state.StorageInstance, error) {
			return s.storageInstance, nil
		},
		storageInstanceVolume: func(tag names.StorageTag) (state.Volume, error) {
			return s.volume, nil
		},
		watchVolumeAttachment: func(names.MachineTag, names.VolumeTag) state.NotifyWatcher {
			return s.volumeAttachmentWatcher
		},
		watchBlockDevices: func(names.MachineTag) state.NotifyWatcher {
			return s.blockDevicesWatcher
		},
		watchStorageAttachment: func(names.StorageTag, names.UnitTag) state.NotifyWatcher {
			return s.storageAttachmentWatcher
		},
	}
}

func (s *watchStorageAttachmentSuite) TestWatchStorageAttachmentVolumeAttachmentChanges(c *gc.C) {
	s.testWatchBlockStorageAttachment(c, func() {
		s.volumeAttachmentWatcher.C <- struct{}{}
	})
}

func (s *watchStorageAttachmentSuite) TestWatchStorageAttachmentStorageAttachmentChanges(c *gc.C) {
	s.testWatchBlockStorageAttachment(c, func() {
		s.storageAttachmentWatcher.C <- struct{}{}
	})
}

func (s *watchStorageAttachmentSuite) TestWatchStorageAttachmentBlockDevicesChange(c *gc.C) {
	s.testWatchBlockStorageAttachment(c, func() {
		s.blockDevicesWatcher.C <- struct{}{}
	})
}

func (s *watchStorageAttachmentSuite) testWatchBlockStorageAttachment(c *gc.C, change func()) {
	s.testWatchStorageAttachment(c, change)
	s.st.CheckCallNames(c,
		"StorageInstance",
		"StorageInstanceVolume",
		"WatchVolumeAttachment",
		"WatchBlockDevices",
		"WatchStorageAttachment",
	)
}

func (s *watchStorageAttachmentSuite) testWatchStorageAttachment(c *gc.C, change func()) {
	w, err := uniter.WatchStorageAttachment(
		s.st,
		s.st,
		s.st,
		s.storageTag,
		s.machineTag,
		s.unitTag,
	)
	c.Assert(err, jc.ErrorIsNil)
	wc := statetesting.NewNotifyWatcherC(c, nopSyncStarter{}, w)
	wc.AssertOneChange()
	change()
	wc.AssertOneChange()
}
