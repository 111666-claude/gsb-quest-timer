package questtimer

import "testing"

func TestUnknownQuestIsNotExpired(t *testing.T) {
	if New().Expired("ghost", 100000) {
		t.Fatal("没有开始的任务不该算超时")
	}
}

func TestStartThenNotExpired(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if timer.Expired("q1", 500) {
		t.Fatal("时限内不该算超时")
	}
}

func TestPauseReturnsTrue(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Pause("q1", 100) {
		t.Fatal("暂停应该成功")
	}
}

func TestInitiallyNotPaused(t *testing.T) {
	if New().Paused("q1") {
		t.Fatal("没有暂停的任务不该是暂停状态")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestPauseExtendsDeadline(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	timer.Pause("q1", 200)
	timer.Resume("q1", 700)
	if timer.Expired("q1", 1400) {
		t.Fatal("暂停 500ms 后截止时刻应顺延到 1500")
	}
	if !timer.Expired("q1", 1500) {
		t.Fatal("顺延后的截止时刻到点应算超时")
	}
}

func TestPausedNeverExpires(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	timer.Pause("q1", 200)
	if timer.Expired("q1", 1000000) {
		t.Fatal("暂停中永远不算超时")
	}
}

func TestExpiredAtExactDeadline(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Expired("q1", 1000) {
		t.Fatal("正好到点应算超时")
	}
	if timer.Expired("q1", 999) {
		t.Fatal("到点前一毫秒不该算超时")
	}
}

func TestSettleRequiresExpired(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if timer.Settle("q1", 500, "r1") {
		t.Fatal("未超时不该结算成功")
	}
	if !timer.Settle("q1", 1000, "r1") {
		t.Fatal("超时后应能结算成功")
	}
}

func TestSettleIdempotentByReqID(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Settle("q1", 1200, "same") {
		t.Fatal("首次结算应成功")
	}
	if timer.Settle("q1", 1300, "same") {
		t.Fatal("同一 reqID 重复结算应返回 false")
	}
}

func TestSettledQuestSettlesOnlyOnce(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Settle("q1", 1200, "r1") {
		t.Fatal("首次结算应成功")
	}
	if timer.Settle("q1", 1300, "r2") {
		t.Fatal("已结算任务不应再次结算")
	}
}

func TestSettledQuestStaysExpired(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	timer.Settle("q1", 1200, "r1")
	timer.Pause("q1", 1300)
	timer.Resume("q1", 2000)
	if !timer.Expired("q1", 1500) {
		t.Fatal("已结算任务不应因继续计时而改变状态")
	}
}

func TestRepeatedSequenceIsDeterministic(t *testing.T) {
	run := func() (bool, bool, bool) {
		timer := New()
		timer.Start("q1", 1000, 0)
		timer.Pause("q1", 200)
		timer.Resume("q1", 700)
		expired := timer.Expired("q1", 1400)
		settled := timer.Settle("q1", 1500, "r1")
		second := timer.Settle("q1", 1600, "r1")
		return expired, settled, second
	}
	firstExpired, firstSettled, firstSecond := run()
	for i := 0; i < 100; i++ {
		expired, settled, second := run()
		if expired != firstExpired || settled != firstSettled || second != firstSecond {
			t.Fatal("同一串操作序列重复执行结果应相同")
		}
	}
}

func TestScannedDoesNotScaleWithQuestCount(t *testing.T) {
	timer := New()
	for index := 0; index < 100000; index++ {
		timer.Start(string(rune(index)), 1000, 0)
	}
	before := timer.Scanned()
	for index := 0; index < 3000; index++ {
		timer.Expired("q", 500)
	}
	if got := timer.Scanned() - before; got > 6000 {
		t.Fatalf("3000 次查询扫描量应在 6000 以内，实际 %d", got)
	}
}
