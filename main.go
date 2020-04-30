package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tjarratt/babble"
)

// タイムリミットのデフォルト 10 秒
const timeLimit = 10

// スコア初期値
var score int = 0

var (
	babbler  babble.Babbler
	question string
	d        time.Duration
)

func init() {
	flag.DurationVar(&d, "d", timeLimit*time.Second, "duration flag")
	flag.Parse()

	babbler = babble.NewBabbler()
	babbler.Count = 1
}

func main() {
	_main()
}

func _main() {
	// ゲーム開始前の 3,2,1 Go 表示
	countdown()

	// タイムアウト処理付き context
	bc := context.Background()
	ctx, cancel := context.WithTimeout(bc, d)
	defer cancel()

	start := time.Now()

	// 第一問
	q()

	ch := input(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n\ntime up !")
			fmt.Println("score:", score)
			return
		case v := <-ch:
			if v == question {
				score++
				fmt.Println("(^-^) good !")
			} else {
				fmt.Println("(>_<) oops...")
			}
			end := time.Now()
			fmt.Printf("%d秒経過\n", int((end.Sub(start)).Seconds()))
			q()
		}
	}
}

func countdown() {
	for i := 3; i > 0; i-- {
		fmt.Print(i)
		fmt.Print(" ")
		time.Sleep(time.Second)
	}
	fmt.Println("Go !")
}

func q() {
	question = babbler.Babble()
	fmt.Println("\ntype this: ", question)
	fmt.Print("> ")
}

func input(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		// 標準入力から一行ずつ文字を読み込む
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- s.Text()
		}
		close(ch)
	}()
	return ch
}
