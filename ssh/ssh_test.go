package ssh

import (
    "testing"
)

func TestConekta(t *testing.T)  {
    t.Parallel()
    commands := []string{
		"echo test", `for i in $(ls); do echo "$i"; done`, "ls",
	}
    for _, cmd := range commands {
		out, err := Conekta("gophers", "gophers", cmd)
		if err != nil {
			t.Errorf("Run failed: %s", err)
		}
		if out == "" {
			t.Errorf("Output was empty for command: %s", cmd)
		}
	}
}