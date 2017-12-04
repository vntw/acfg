package ac

import (
	"testing"

	"github.com/venyii/acsrvmanager/server/ac/server"
)

func TestNewMemoryInstanceManager(t *testing.T) {
	im := NewMemoryInstanceManager()

	if im.instances == nil {
		t.Error("expected empty instances, got nil")
	}
}

func TestGetInstancesWithNoData(t *testing.T) {
	im := NewMemoryInstanceManager()

	if len(im.GetInstances()) != 0 {
		t.Fatalf("got wrong server log count: got %v want 0", len(im.GetInstances()))
	}
}

func TestGetInstances(t *testing.T) {
	im := NewMemoryInstanceManager()
	instance := server.Instance{}
	im.addInstance(instance)

	if len(im.GetInstances()) != 1 {
		t.Fatalf("got wrong server log count: got %v want 1", len(im.GetInstances()))
	}
}

func TestStopInstanceNotFound(t *testing.T) {
	im := NewMemoryInstanceManager()
	err := im.StopInstance("unknown")

	if err == nil {
		t.Fatal("expected error to be returned, got nil")
	}

	if err.Error() != "instance not found" {
		t.Fatalf("unexpected error msg, got %s", err.Error())
	}
}
