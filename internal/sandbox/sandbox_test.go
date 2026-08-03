package sandbox

import (
	"context"
	"testing"
	"time"

	"github.com/sd3/linuxlab/internal/challenge"
)

// resetDockerProbe clears the in-process DockerAvailable cache so tests do
// not observe (or leave behind) stale probe state.
func resetDockerProbe() {
	dockerProbeMu.Lock()
	defer dockerProbeMu.Unlock()
	dockerProbeResult = false
	dockerProbeTime = time.Time{}
}

func setDockerProbe(result bool, at time.Time) {
	dockerProbeMu.Lock()
	defer dockerProbeMu.Unlock()
	dockerProbeResult = result
	dockerProbeTime = at
}

func TestDockerSandbox_ImplementsSandbox(t *testing.T) {
	var _ Sandbox = (*DockerSandbox)(nil) // compile-time check
}

func TestDockerAvailable(t *testing.T) {
	resetDockerProbe()
	defer resetDockerProbe()
	// Just verify it returns a bool without panicking
	result := DockerAvailable()
	t.Logf("DockerAvailable() = %v", result)
}

func TestDockerAvailable_UsesCachedResult(t *testing.T) {
	defer resetDockerProbe()

	// Whatever the real Docker state is, one of the two seeded values
	// contradicts it — DockerAvailable returning the seeded value in both
	// cases proves the cache is honored instead of re-probing.
	for _, cached := range []bool{true, false} {
		setDockerProbe(cached, time.Now())
		if got := DockerAvailable(); got != cached {
			t.Errorf("DockerAvailable() = %v, want cached value %v", got, cached)
		}
	}
}

func TestDockerAvailable_RefreshesExpiredCache(t *testing.T) {
	defer resetDockerProbe()

	setDockerProbe(true, time.Now().Add(-time.Hour))
	_ = DockerAvailable()

	dockerProbeMu.Lock()
	refreshed := time.Since(dockerProbeTime) < time.Minute
	dockerProbeMu.Unlock()
	if !refreshed {
		t.Error("expired cache entry was not re-probed")
	}
}

func TestNewSandbox_ContainersCategoryUsesLocal(t *testing.T) {
	// containers challenges without a compose file need the host docker
	// CLI, so they must not run inside a DockerSandbox container —
	// regardless of whether Docker is available.
	ch := &challenge.Challenge{
		ID:       "docker-ps",
		Category: "containers",
	}
	sb, err := NewSandbox(context.Background(), ch)
	if err != nil {
		t.Fatalf("NewSandbox: %v", err)
	}
	defer sb.Destroy(context.Background())

	if _, ok := sb.(*LocalSandbox); !ok {
		t.Errorf("containers challenge should get *LocalSandbox, got %T", sb)
	}
}

func TestNewSandbox_LocalFallback(t *testing.T) {
	resetDockerProbe()
	defer resetDockerProbe()
	if DockerAvailable() {
		t.Skip("Docker is available; cannot test local fallback")
	}
	ch := &challenge.Challenge{
		ID:       "test-challenge",
		Category: "linux-basics",
	}
	sb, err := NewSandbox(context.Background(), ch)
	if err != nil {
		t.Fatalf("NewSandbox: %v", err)
	}
	defer sb.Destroy(context.Background())

	if _, ok := sb.(*LocalSandbox); !ok {
		t.Errorf("expected *LocalSandbox, got %T", sb)
	}
}
