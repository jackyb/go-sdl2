package gfx

import (
	"testing"
)

func TestMin(t *testing.T) {
	if a := []int{4, 2, 3, 5, 1, 6, 5, 7}; min(a...) != 1 {
		t.Fail()
	}

	if a := []int{5, 6, -2, 1, 8, 1, 4, 3}; min(a...) != -2 {
		t.Fail()
	}
}

func TestFPSmanager(t *testing.T) {
	fps_manager := FPSmanager{}
	InitFramerate(&fps_manager)

	if !SetFramerate(&fps_manager, 60) {
		t.Errorf("gfx.SetFramerate failed")
	}

	if fps, ok := GetFramerate(&fps_manager); !ok || fps != 60 {
		t.Errorf("gfx.GetFramerate() = (%v, %v) - want (60, true)", fps, ok)
	}

	// frame count is initialized to 0
	if count, ok := GetFramecount(&fps_manager); !ok || count != 0 {
		t.Errorf("gfx.GetFramecount() = (%v, %v) - want (0, true)", count, ok)
	}

	// count one frame
	FramerateDelay(&fps_manager)

	// frame count should now be 1
	if count, ok := GetFramecount(&fps_manager); !ok || count != 1 {
		t.Errorf("gfx.Framecount() = (%v, %v) - want (1, true)", count, ok)
	}
}
