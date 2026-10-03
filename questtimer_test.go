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
