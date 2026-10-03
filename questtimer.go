// Package questtimer 是任务限时：暂停顺延、到期判定与超时结算幂等。
package questtimer

// Timer 是任务计时台账。
type Timer struct {
	deadline map[string]int64
	pausedAt map[string]int64
	settled  map[string]bool
	memo     map[string]bool
	scanned  int
}

// New 建台账。
func New() *Timer {
	return &Timer{deadline: map[string]int64{}, pausedAt: map[string]int64{},
		settled: map[string]bool{}, memo: map[string]bool{}}
}

// Start 开始一个限时任务。
func (t *Timer) Start(quest string, durationMs, nowMs int64) {
	t.deadline[quest] = nowMs + durationMs
	delete(t.pausedAt, quest)
	delete(t.settled, quest)
}

// Pause 暂停计时，暂停期间不计时。
func (t *Timer) Pause(quest string, nowMs int64) bool {
	if _, ok := t.deadline[quest]; !ok {
		return false
	}
	if _, ok := t.pausedAt[quest]; ok {
		return false
	}
	t.pausedAt[quest] = nowMs
	return true
}

// Resume 恢复计时，把暂停时长补到截止时刻上。
func (t *Timer) Resume(quest string, nowMs int64) bool {
	pausedAt, ok := t.pausedAt[quest]
	if !ok {
		return false
	}
	t.deadline[quest] += nowMs - pausedAt
	delete(t.pausedAt, quest)
	return true
}

// Expired 判断任务是否超时：当前时刻不小于截止时刻且没有暂停。
func (t *Timer) Expired(quest string, nowMs int64) bool {
	t.scanned++
	deadline, ok := t.deadline[quest]
	if !ok {
		return false
	}
	if _, paused := t.pausedAt[quest]; paused {
		return false
	}
	return nowMs >= deadline
}

// Settle 结算一次超时，返回是否成功；同 reqID 只结算一次。
func (t *Timer) Settle(quest string, nowMs int64, reqID string) bool {
	if t.memo[reqID] || t.settled[quest] {
		return false
	}
	if !t.Expired(quest, nowMs) {
		return false
	}
	t.memo[reqID] = true
	t.settled[quest] = true
	return true
}

// Paused 判断任务是否处于暂停中。
func (t *Timer) Paused(quest string) bool {
	_, ok := t.pausedAt[quest]
	return ok
}

// Scanned 是累计扫描量（规模观测）。
func (t *Timer) Scanned() int { return t.scanned }
