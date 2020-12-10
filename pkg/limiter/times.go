package limiter

import (
	"github.com/morgine/moon/pkg/x_time"
	"sync"
	"time"
)

// 根据限制次数获取限制时间及限制解除时间
type LimitTimeProvider func(times int) (limitIn, clearIn time.Duration)

// 次数限制器，可用于 IP 封禁或用户账户登陆封禁
type TimesLimiter struct {
	limits   map[string]*Limit
	provider LimitTimeProvider
	mu       sync.Mutex
}

func NewTimesLimiter(provider LimitTimeProvider) *TimesLimiter {
	return &TimesLimiter{
		limits:   map[string]*Limit{},
		provider: provider,
	}
}

type Limit struct {
	Times   int        // 被限制次数
	LimitAt *time.Time // 限制解除时间
	ClearAt *time.Time // 限制清除时间
}

// 解除限制
func (ts *TimesLimiter) RemoveLimit(key string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.limits, key)
}

// 获得限制剩余时间
func (ts *TimesLimiter) CheckLimit(key string) (limitIn, clearIn time.Duration) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	l := ts.limits[key]
	if l != nil {
		now := x_time.Now()
		if l.LimitAt != nil {
			limitIn = l.LimitAt.Sub(now)
		}
		if l.ClearAt != nil {
			clearIn = l.ClearAt.Sub(now)
		}
	}
	return
}

// 增加一次次数，并返回限制剩余时间
func (ts *TimesLimiter) AddOneTimes(key string) (limitIn, clearIn time.Duration) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	l := ts.limits[key]
	now := x_time.Now()
	if l == nil || l.ClearAt == nil || l.ClearAt.Before(now) {
		if l == nil {
			l = &Limit{}
			ts.limits[key] = l
		}
		// 随机查询 3 个限制是否过期，相当于垃圾清理器，不同于全局垃圾清理，这种策略只能清理 2/3 的垃圾，
		// 但不会因为大量清理垃圾而导致程序卡顿
		ts.removeExpired(3)
	}
	l.Times++
	limitIn, clearIn = ts.provider(l.Times)
	if limitIn > 0 {
		limitAt := now.Add(limitIn)
		l.LimitAt = &limitAt
	} else {
		l.LimitAt = nil
	}
	if clearIn > 0 {
		clearAt := now.Add(clearIn)
		l.ClearAt = &clearAt
	} else {
		l.LimitAt = nil
	}
	return
}

// 随机查询 n 个限制是否达到清除时间, 达到清除时间则清除该限制，需要上锁
func (ts *TimesLimiter) removeExpired(random int) {
	now := x_time.Now()
	for key, limit := range ts.limits {
		if limit.ClearAt != nil && limit.ClearAt.Before(now) {
			delete(ts.limits, key)
		}
		random--
		if random <= 0 {
			return
		}
	}
}
