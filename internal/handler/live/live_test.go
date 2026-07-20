package live

import "testing"

func TestApplyLeaderSkillBonus(t *testing.T) {
	tests := []struct {
		name                          string
		effectType                    int
		effectValue                   float64
		smile, pure, cool             float64
		wantSmile, wantPure, wantCool float64
	}{
		{"smile", 1, 12, 10000, 20000, 30000, 1200, 0, 0},
		{"pure", 2, 12, 10000, 20000, 30000, 0, 2400, 0},
		{"cool", 3, 12, 10000, 20000, 30000, 0, 0, 3600},
		{"smile-to-pure", 112, 12.5, 1001, 20000, 30000, 0, 126, 0},
		{"smile-to-cool", 113, 12, 10000, 20000, 30000, 0, 0, 1200},
		{"pure-to-smile", 121, 12, 10000, 20000, 30000, 2400, 0, 0},
		{"pure-to-cool", 123, 12, 10000, 20000, 30000, 0, 0, 2400},
		{"cool-to-smile", 131, 12, 10000, 20000, 30000, 3600, 0, 0},
		{"cool-to-pure", 132, 12, 10000, 20000, 30000, 0, 3600, 0},
		{"unsupported", 999, 12, 10000, 20000, 30000, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smile, pure, cool := applyLeaderSkillBonus(tt.effectType, tt.effectValue, tt.smile, tt.pure, tt.cool)
			if smile != tt.wantSmile || pure != tt.wantPure || cool != tt.wantCool {
				t.Errorf("applyLeaderSkillBonus(%d) = (%v, %v, %v), want (%v, %v, %v)", tt.effectType, smile, pure, cool, tt.wantSmile, tt.wantPure, tt.wantCool)
			}
		})
	}
}
