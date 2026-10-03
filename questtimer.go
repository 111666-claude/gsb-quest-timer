// Package questtimer 是任务限时：暂停顺延、到期判定与超时结算幂等。
// 缺陷：暂停不影响计时、恢复不补偿暂停时长、到期判定漏掉等号、暂停中也算过期、
// 结算不看是否过期也不幂等、每次查询遍历全部任务。
package questtimer

// Timer 是任务计时台账。
type Timer struct {
	deadline map[string]int64
	pausedAt map[string]int64
	pausedMs map[string]int64
	memo     map[string]bool
	scanned  int
}

// New 建台账。
func New() *Timer {
	return &Timer{deadline: map[string]int64{}, pausedAt: map[string]int64{},
		pausedMs: map[string]int64{}, memo: map[string]bool{}}
}

// Start 开始一个限时任务。
func (t *Timer) Start(quest string, durationMs, nowMs int64) {
	t.deadline[quest] = nowMs + durationMs
}

// Pause 暂停计时。
func (t *Timer) Pause(quest string, nowMs int64) bool {
	t.scanned += len(t.deadline)
	t.pausedAt[quest] = 0
	return true
}

// Resume 恢复计时。
func (t *Timer) Resume(quest string, nowMs int64) bool {
	t.scanned += len(t.deadline)
	delete(t.pausedAt, quest)
	return true
}

// Expired 判断任务是否超时。
func (t *Timer) Expired(quest string, nowMs int64) bool {
	t.scanned += len(t.deadline)
	deadline, ok := t.deadline[quest]
	if !ok {
		return false
	}
	return nowMs > deadline
}

// Settle 结算一次超时，返回是否成功。
func (t *Timer) Settle(quest string, nowMs int64, reqID string) bool {
	t.scanned += len(t.deadline)
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
