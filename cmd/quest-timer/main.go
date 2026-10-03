// Command quest-timer 跑任务限时样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/questtimer"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("quest-timer", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "pause", "pause / edge / settle / idem / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "pause":
		timer := questtimer.New()
		timer.Start("q1", 1000, 0)
		timer.Pause("q1", 200)
		timer.Resume("q1", 700)
		fmt.Fprintf(stdout, "expired=%v\n", timer.Expired("q1", 1400))
	case "edge":
		timer := questtimer.New()
		timer.Start("q1", 1000, 0)
		fmt.Fprintf(stdout, "expired=%v\n", timer.Expired("q1", 1000))
	case "settle":
		timer := questtimer.New()
		timer.Start("q1", 1000, 0)
		fmt.Fprintf(stdout, "settled=%v\n", timer.Settle("q1", 500, "r1"))
	case "idem":
		timer := questtimer.New()
		timer.Start("q1", 1000, 0)
		timer.Settle("q1", 1200, "same")
		fmt.Fprintf(stdout, "second=%v\n", timer.Settle("q1", 1300, "same"))
	case "scan":
		timer := questtimer.New()
		for index := 0; index < 3000; index++ {
			timer.Start(fmt.Sprintf("q-%d", index), 1000, 0)
		}
		for index := 0; index < 3000; index++ {
			timer.Expired("q-2999", 500)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", timer.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample pause|edge|settle|idem|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
