package limiter_test

import (
	"github.com/morgine/moon/pkg/limiter"
	"github.com/morgine/moon/pkg/x_time"
	"testing"
	"time"
)

func TestNewTimesLimiter(t *testing.T) {
	type testcase struct {
		times   int
		limitIn time.Duration
		clearIn time.Duration
	}
	// n 次及 n 次以上限制时间逐步递增
	limitBase := time.Minute     // 限制时间基数
	clearBase := 2 * time.Minute // 清除时间基数
	n := 5                       // 限制最小基数

	var testcases = []testcase{
		{times: 1, limitIn: 0, clearIn: 0},
		{times: 2, limitIn: 0, clearIn: 0},
		{times: 3, limitIn: 0, clearIn: 0},
		{times: 4, limitIn: 0, clearIn: 0},
		{times: 5, limitIn: 0, clearIn: 0},
		{times: 6, limitIn: limitBase * 1, clearIn: clearBase * 1},
		{times: 7, limitIn: limitBase * 2, clearIn: clearBase * 2},
		{times: 8, limitIn: limitBase * 3, clearIn: clearBase * 3},
	}

	callAt(time.Now(), func() {
		timesLimiter := limiter.NewTimesLimiter(func(times int) (limitIn, clearIn time.Duration) {
			if 0 < times && times < n {
				// 1-n 次不设置时间
				return 0, 0
			} else {
				limitIn = time.Duration(times-n) * limitBase
				clearIn = time.Duration(times-n) * clearBase
				return
			}
		})
		key := "user_01"
		for _, tc := range testcases {
			timesLimiter.RemoveLimit(key)
			for i := 0; i < tc.times; i++ {
				timesLimiter.AddOneTimes(key)
			}
			limitIn, clearIn := timesLimiter.CheckLimit(key)
			if limitIn != tc.limitIn {
				t.Errorf("%d need: %v, got: %v\n", tc.times, tc.limitIn, limitIn)
			}
			if clearIn != tc.clearIn {
				t.Errorf("%d need: %v, got: %v\n", tc.times, tc.clearIn, clearIn)
			}
		}
	})
}

// 重置当前时间
func callAt(t time.Time, call func()) {
	x_time.Now = func() time.Time {
		return t
	}
	call()
	x_time.Now = time.Now
}
