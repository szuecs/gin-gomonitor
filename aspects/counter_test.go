package ginmon

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestMode string = "test"

const checkMark = "\u2713"
const ballotX = "\u2717"

func Test_Inc(t *testing.T) {
	c := CounterAspect{1}
	expect := 2
	c.Inc()
	if assert.Equal(t, c.Count, expect, "Incrementation of counter does not work, expect %d but got %d %s",
		expect, c.Count, ballotX) {
		t.Logf("Incrementation of counter works, expect %d and git %d %s",
			expect, c.Count, checkMark)
	}
}

func Test_GetStats(t *testing.T) {
	c := CounterAspect{1}
	if assert.NotNil(t, c.GetStats(), "Return of counter getstats should not be nil") {
		t.Logf("Should be an interface %s", checkMark)
	}
}

func Test_Name(t *testing.T) {
	c := CounterAspect{1}
	expect := "Counter"
	if assert.Equal(t, c.Name(), expect, "Return of counter name does not work, expect %s but got %s %s",
		expect, c.Name(), ballotX) {
		t.Logf("Return of counter name works, expect %s and got %s %s",
			expect, c.Name(), checkMark)
	}
}

func Test_InRoot(t *testing.T) {
	c := CounterAspect{1}
	expect := false
	if assert.Equal(t, c.InRoot(), expect, "Expect %v but got %v %s",
		expect, c.InRoot(), ballotX) {
		t.Logf("Expect %v and got %v %s",
			expect, c.InRoot(), checkMark)
	}
}

func Test_CounterHandler(t *testing.T) {
	gin.SetMode(TestMode)
	router := gin.New()
	cnt := &CounterAspect{1}
	router.Use(CounterHandler(cnt))
	tryRequest(router, "GET", "/")
	if assert.Equal(t, cnt.Count, 2, "Incrementation of counter does not work, expect %d but got %d %s", 2, cnt.Count, ballotX) {
		t.Logf("CounterHandler works, expect %d and got %d %s", 2, cnt.Count, checkMark)
	}
}

func tryRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
