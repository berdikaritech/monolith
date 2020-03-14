package micro

import (
	"reflect"
	"testing"
)

type testClient struct{}

func (c *testClient) Name() string {
	return "client"
}

func TestClientSet_NewClientSet(t *testing.T) {
	tests := []struct {
		name string
		args []Client
		want int
	}{
		{"With no clients", []Client{}, 0},
		{"With a client", []Client{&testClient{}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewClientSet(tt.args...)
			if got := cs.Length(); got != tt.want {
				t.Errorf("ClientSet.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientSet_Get(t *testing.T) {
	c := &testClient{}
	cs := NewClientSet(c)

	tests := []struct {
		name string
		c    *ClientSet
		args string
		want Client
	}{
		{"A valid key", cs, c.Name(), c},
		{"An invalid key", cs, "invalid", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Get(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientSet.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientSet_Keys(t *testing.T) {
	c := &testClient{}

	tests := []struct {
		name string
		c    *ClientSet
		want []string
	}{
		{"An empty ClientSet", NewClientSet(), []string{}},
		{"A non-empty ClientSet", NewClientSet(c), []string{c.Name()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientSet.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}
