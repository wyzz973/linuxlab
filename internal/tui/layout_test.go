package tui

import "testing"

func TestComputeLayout_Breakpoints(t *testing.T) {
	cases := []struct {
		w, h int
		want LayoutMode
	}{
		{0, 0, ModeUnsupported},
		{59, 24, ModeUnsupported},
		{60, 17, ModeUnsupported},
		{60, 18, ModeCompact},
		{79, 30, ModeCompact},
		{80, 23, ModeCompact},
		{80, 24, ModeStandard},
		{100, 30, ModeStandard},
		{119, 40, ModeStandard},
		{120, 29, ModeStandard}, // wide width but short: stays standard
		{120, 30, ModeWide},
		{140, 40, ModeWide},
	}
	for _, c := range cases {
		got := computeLayout(c.w, c.h)
		if got.Mode != c.want {
			t.Errorf("computeLayout(%d, %d).Mode = %d, want %d", c.w, c.h, got.Mode, c.want)
		}
	}
}

func TestComputeLayout_Budgets(t *testing.T) {
	spec := computeLayout(100, 30)
	if spec.HeaderH != 1 || spec.FooterH != 1 {
		t.Fatalf("header/footer = %d/%d, want 1/1", spec.HeaderH, spec.FooterH)
	}
	if spec.MainW != 100 {
		t.Fatalf("MainW = %d, want 100", spec.MainW)
	}
	if spec.MainH != 30-spec.HeaderH-spec.FooterH {
		t.Fatalf("MainH = %d, want %d", spec.MainH, 30-spec.HeaderH-spec.FooterH)
	}
}

func TestComputeLayout_WideInspector(t *testing.T) {
	for _, size := range [][2]int{{120, 30}, {140, 40}, {200, 50}} {
		spec := computeLayout(size[0], size[1])
		if spec.Mode != ModeWide {
			t.Fatalf("computeLayout(%d, %d).Mode = %d, want ModeWide", size[0], size[1], spec.Mode)
		}
	}
}

func TestComputeLayout_UnsupportedSkipsBudgets(t *testing.T) {
	spec := computeLayout(40, 10)
	if spec.Mode != ModeUnsupported {
		t.Fatalf("Mode = %d, want ModeUnsupported", spec.Mode)
	}
	if spec.MainH != 0 || spec.HeaderH != 0 || spec.FooterH != 0 {
		t.Fatalf("unsupported layout should not allocate slots: %+v", spec)
	}
}
