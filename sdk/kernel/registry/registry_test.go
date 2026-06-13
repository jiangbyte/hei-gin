package registry

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRouteRejectsDuplicate(t *testing.T) {
	ResetForTest()

	fn := func(r *gin.Engine) {}
	RegisterRoute(fn)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected duplicate route registration panic")
		}
	}()
	RegisterRoute(fn)
}

func TestRegisterMiddlewareRejectsAfterFreeze(t *testing.T) {
	ResetForTest()
	Freeze()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected frozen middleware registry panic")
		}
	}()
	RegisterMiddleware(func(r *gin.Engine) {})
}

func TestSnapshotStateShowsFrozenRegistry(t *testing.T) {
	ResetForTest()
	RegisterRoute(func(r *gin.Engine) {})
	Freeze()

	snapshot := SnapshotState()
	if !snapshot.Frozen || len(snapshot.Routes) != 1 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
}
