package pkg

import "testing"

func TestLoadEnv(t *testing.T) {
	env, err := LoadEnv()
	if err != nil {
		t.Fatal(err)
	}
	if env == nil {
		t.Errorf("env=%v, want: other", env)
	}
	t.Logf("env=%v\n", env)
}
