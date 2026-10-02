// Command anim-blend 跑动画混合样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/anim"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("anim-blend", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "override", "override / normalize / missing / transition / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "override":
		mixer := anim.New("idle")
		mixer.AddClip("run")
		mixer.AddClip("aim")
		base := mixer.AddLayer()
		mixer.Play(base, "run", 1)
		top := mixer.AddLayer()
		mixer.Play(top, "aim", 1)
		out := mixer.Sample()
		fmt.Fprintf(stdout, "run=%.2f aim=%.2f\n", out["run"], out["aim"])
	case "normalize":
		mixer := anim.New("idle")
		mixer.AddClip("a")
		mixer.AddClip("b")
		layer := mixer.AddLayer()
		mixer.Play(layer, "a", 3)
		mixer.Play(layer, "b", 1)
		out := mixer.Sample()
		fmt.Fprintf(stdout, "a=%.2f b=%.2f sum=%.2f\n", out["a"], out["b"], out["a"]+out["b"])
	case "missing":
		mixer := anim.New("idle")
		mixer.AddClip("idle")
		base := mixer.AddLayer()
		mixer.Play(base, "idle", 1)
		top := mixer.AddLayer()
		mixer.Play(top, "jump", 1)
		out := mixer.Sample()
		fmt.Fprintf(stdout, "idle=%.2f missing=%d\n", out["idle"], mixer.Missing())
	case "transition":
		fast := anim.New("idle")
		fast.AddClip("run")
		fast.AddClip("idle")
		layerA := fast.AddLayer()
		fast.BeginTransition(layerA, "run", "idle", 300)
		for index := 0; index < 10; index++ {
			fast.Tick(15)
		}
		slow := anim.New("idle")
		slow.AddClip("run")
		slow.AddClip("idle")
		layerB := slow.AddLayer()
		slow.BeginTransition(layerB, "run", "idle", 300)
		for index := 0; index < 3; index++ {
			slow.Tick(50)
		}
		fmt.Fprintf(stdout, "a=%.3f b=%.3f\n", fast.Sample()["idle"], slow.Sample()["idle"])
	case "scan":
		mixer := anim.New("idle")
		for index := 0; index < 3000; index++ {
			mixer.AddClip(fmt.Sprintf("c-%d", index))
		}
		layer := mixer.AddLayer()
		mixer.Play(layer, "c-0", 1)
		for index := 0; index < 3000; index++ {
			mixer.Sample()
		}
		fmt.Fprintf(stdout, "scanned=%d\n", mixer.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample override|normalize|missing|transition|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
