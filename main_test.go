package main

import "testing"

func TestGetPort(t *testing.T) {
	tests := []struct {
		name string
		port string
		want string
	}{
		{name: "default", port: "", want: "8080"},
		{name: "from env", port: "3000", want: "3000"},
		{name: "trims spaces", port: " 9000 ", want: "9000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)

			got := getPort()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
