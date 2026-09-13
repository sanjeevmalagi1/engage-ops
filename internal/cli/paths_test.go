package cli

import "testing"

func TestResolvePaths(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantDir   string
		wantState string
	}{
		{
			name:      "no env leaves defaults untouched",
			args:      nil,
			wantDir:   ".",
			wantState: "engageops.tfstate.json",
		},
		{
			name:      "env alone derives both paths",
			args:      []string{"--env", "staging"},
			wantDir:   "environments/staging",
			wantState: "engageops.staging.tfstate.json",
		},
		{
			name:      "explicit dir wins over env",
			args:      []string{"--env", "staging", "--dir", "custom"},
			wantDir:   "custom",
			wantState: "engageops.staging.tfstate.json",
		},
		{
			name:      "explicit state wins over env",
			args:      []string{"--env", "staging", "--state", "custom.json"},
			wantDir:   "environments/staging",
			wantState: "custom.json",
		},
		{
			name:      "explicit dir and state both win over env",
			args:      []string{"--env", "staging", "--dir", "custom", "--state", "custom.json"},
			wantDir:   "custom",
			wantState: "custom.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := NewRootCommand()
			if err := root.ParseFlags(tc.args); err != nil {
				t.Fatalf("parse flags: %v", err)
			}

			gotDir, gotState := resolvePaths(root)
			if gotDir != tc.wantDir {
				t.Errorf("dir = %q, want %q", gotDir, tc.wantDir)
			}
			if gotState != tc.wantState {
				t.Errorf("state = %q, want %q", gotState, tc.wantState)
			}
		})
	}
}
