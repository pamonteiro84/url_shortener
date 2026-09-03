package shortcode

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://example.com", "EAaArVRs"},
		{"https://google.com", "BQRvJsg-"},
		{"", "47DEQpj8"},
		{"https://exemplo.pt/café", "6iS07GCC"},
		{"https://example.com/", "DxFdsGK3"},
	}

	for _, tt := range tests {
		got := Generate(tt.input)
		if got != tt.want {
			t.Errorf("Generate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGenerate_Deterministic(t *testing.T) {
	input := "https://example.com"

	first := Generate(input)
	second := Generate(input)

	if first != second {
		t.Errorf("Generate(%q) is not deterministic: got %q and %q", input, first, second)
	}
}

func TestGenerate_Length(t *testing.T) {
	got := Generate("https://example.com")

	if len(got) != 8 {
		t.Errorf("Generate() length = %d, want 8", len(got))
	}
}
