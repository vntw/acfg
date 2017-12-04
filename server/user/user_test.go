package user

import (
	"testing"
)

func resetUsers() {
	users = []*User{}
}

func TestAddConfigUsers(t *testing.T) {
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %#v", len(users))
	}

	AddConfigUsers(map[string]string{
		"u1": "p1",
		"u2": "p2",
	})

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %#v", len(users))
	}

	resetUsers()
}

func TestMatchUser(t *testing.T) {
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %#v", len(users))
	}

	addUser("u1", "p1")

	if _, err := MatchUser("u1", "p1"); err != nil {
		t.Errorf("expected to match user, got %v", err)
	}
	if _, err := MatchUser("u1", "p2"); err == nil {
		t.Error("expected not to match user with same username and different password")
	}
	if _, err := MatchUser("u2", "p1"); err == nil {
		t.Error("expected not to match user with different username and same password")
	}
	if _, err := MatchUser("u2", "p2"); err == nil {
		t.Error("expected not to match user with different data")
	}

	resetUsers()
}
