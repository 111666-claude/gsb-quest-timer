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
		t.Fatal("暂停 500ms 应该顺延截止时刻，1400 不该超时")
	}
	if !timer.Expired("q1", 1500) {
		t.Fatal("顺延后的截止时刻 1500 应该超时")
	}
}

func TestExpiredAtExactDeadline(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Expired("q1", 1000) {
		t.Fatal("正好到点应该算超时")
	}
	if timer.Expired("q1", 999) {
		t.Fatal("到点前一刻不该算超时")
	}
}

func TestPausedNeverExpires(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	timer.Pause("q1", 200)
	if timer.Expired("q1", 100000) {
		t.Fatal("暂停中永远不算超时")
	}
}

func TestSettleRequiresExpired(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if timer.Settle("q1", 500, "r1") {
		t.Fatal("没超时不该结算成功")
	}
	if !timer.Settle("q1", 1000, "r1") {
		t.Fatal("已超时应该结算成功")
	}
}

func TestSettleIdempotentByReqID(t *testing.T) {
	timer := New()
	timer.Start("q1", 1000, 0)
	if !timer.Settle("q1", 1200, "same") {
		t.Fatal("首次结算应该成功")
	}
	if timer.Settle("q1", 1300, "same") {
		t.Fatal("同一 reqID 重复结算应该返回 false")
	}
	if timer.Settle("q1", 1400, "other") {
		t.Fatal("已结算的任务不该再次结算")
	}
}

func TestScannedDoesNotScaleWithQuests(t *testing.T) {
	timer := New()
	for index := 0; index < 100000; index++ {
		timer.Start(string(rune(index)), 1000, 0)
	}
	before := timer.Scanned()
	for index := 0; index < 3000; index++ {
		timer.Expired("q", 500)
	}
	if got := timer.Scanned() - before; got > 6000 {
		t.Fatalf("3000 次查询扫描量 %d 超过 6000", got)
	}
}
