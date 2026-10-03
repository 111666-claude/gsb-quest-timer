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
	return &Timer{
		deadline: map[string]int64{},
		pausedAt: map[string]int64{},
		settled:  map[string]bool{},
		memo:     map[string]bool{},
	}
}

// Start 开始一个限时任务。
func (t *Timer) Start(quest string, durationMs, nowMs int64) {
	t.deadline[quest] = nowMs + durationMs
}

// Pause 暂停计时；重复暂停返回 false。
func (t *Timer) Pause(quest string, nowMs int64) bool {
	if _, ok := t.deadline[quest]; !ok {
		return false
	}
	if _, ok := t.pausedAt[quest]; ok {
		return false
	}
	t.scanned++
	t.pausedAt[quest] = nowMs
	return true
}

// Resume 恢复计时，并把暂停时长补到截止时刻上。
func (t *Timer) Resume(quest string, nowMs int64) bool {
	pausedSince, ok := t.pausedAt[quest]
	if !ok {
		return false
	}
	t.scanned++
	t.deadline[quest] += nowMs - pausedSince
	delete(t.pausedAt, quest)
	return true
}

// Expired 判断任务是否超时：当前时刻不小于截止时刻且未暂停时算超时。
func (t *Timer) Expired(quest string, nowMs int64) bool {
	t.scanned++
	if t.settled[quest] {
		return true
	}
	if _, paused := t.pausedAt[quest]; paused {
		return false
	}
	deadline, ok := t.deadline[quest]
	if !ok {
		return false
	}
	return nowMs >= deadline
}

// Settle 结算一次超时：只有已经超时、未结算且 reqID 未用过时才成功。
func (t *Timer) Settle(quest string, nowMs int64, reqID string) bool {
	t.scanned++
	if !t.Expired(quest, nowMs) {
		return false
	}
	if t.settled[quest] || t.memo[reqID] {
		return false
	}
	t.settled[quest] = true
	t.memo[reqID] = true
	return true
}

// Paused 判断任务是否处于暂停中。
func (t *Timer) Paused(quest string) bool {
	_, ok := t.pausedAt[quest]
	return ok
}

// Scanned 是累计扫描量（规模观测）。
func (t *Timer) Scanned() int { return t.scanned }
