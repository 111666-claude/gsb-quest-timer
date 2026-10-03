# quest-timer

任务限时：暂停顺延、到期判定与超时结算幂等，只用 Go 标准库。

```
go run ./cmd/quest-timer --sample=pause
go run ./cmd/quest-timer --sample=edge
go run ./cmd/quest-timer --sample=settle
go run ./cmd/quest-timer --sample=idem
go run ./cmd/quest-timer --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **暂停**：暂停期间不计时；`Resume` 时把这段暂停时长补到截止时刻上。
- **到期**：当前时刻不小于截止时刻且没有暂停时算超时；暂停中永远不算超时。
- **结算**：只有已经超时的任务才能结算，否则返回 false。
- **幂等**：同一个 `reqID` 只结算一次，重复调用返回 false。
- **不变量**：已结算的任务不会因为继续计时而改变状态；同一串开始、暂停、恢复与
  结算序列重复执行结果相同。
- **代价**：`Expired` 不许遍历全部任务，`Scanned` 不随任务数线性增长。
- **规模**：10 万个任务、每秒 10 万次查询，单次摊还 O(1)。

## 输出契约（不改格式）

```
expired=..
settled=..
second=..
scanned=..
```

## 目录

```
questtimer.go          暂停、到期与结算
questtimer_test.go     单元测试
cmd/quest-timer/       命令行入口
```
