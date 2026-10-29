package services

import (
	"strings"
	"testing"

	"minitask/models"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"válido", "Hola", "Hola", false},
		{"recorta espacios", "  Hola  ", "Hola", false},
		{"muy corto", "ab", "", true},
		{"vacío", "", "", true},
		{"solo espacios", "   ", "", true},
		{"mínimo permitido", "abc", "abc", false},
		{"máximo permitido", strings.Repeat("a", 100), strings.Repeat("a", 100), false},
		{"supera máximo", strings.Repeat("a", 101), "", true},
		{"Unicode", "café", "café", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTitle(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"vacía", "", "", false},
		{"conserva espacios", "  descripción  ", "  descripción  ", false},
		{"exactamente 1000 caracteres", strings.Repeat("a", 1000), strings.Repeat("a", 1000), false},
		{"supera 1000 caracteres", strings.Repeat("a", 1001), "", true},
		{"1000 caracteres ñ", strings.Repeat("ñ", 1000), strings.Repeat("ñ", 1000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateDescription(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateEnum(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		allowed []string
		field   string
		wantErr bool
	}{
		{
			name:    "estado permitido",
			value:   models.StatusPending,
			allowed: []string{models.StatusPending, models.StatusInProgress, models.StatusCompleted},
			field:   "estado",
			wantErr: false,
		},
		{
			name:    "estado desconocido",
			value:   "archived",
			allowed: []string{models.StatusPending, models.StatusInProgress, models.StatusCompleted},
			field:   "estado",
			wantErr: true,
		},
		{
			name:    "estado vacío",
			value:   "",
			allowed: []string{models.StatusPending, models.StatusInProgress, models.StatusCompleted},
			field:   "estado",
			wantErr: true,
		},
		{
			name:    "mayúsculas no permitidas",
			value:   "PENDING",
			allowed: []string{models.StatusPending, models.StatusInProgress, models.StatusCompleted},
			field:   "estado",
			wantErr: true,
		},
		{
			name:    "prioridad permitida",
			value:   models.PriorityHigh,
			allowed: []string{models.PriorityLow, models.PriorityMedium, models.PriorityHigh},
			field:   "prioridad",
			wantErr: false,
		},
		{
			name:    "prioridad desconocida",
			value:   "urgent",
			allowed: []string{models.PriorityLow, models.PriorityMedium, models.PriorityHigh},
			field:   "prioridad",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnum(tt.value, tt.allowed, tt.field)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
