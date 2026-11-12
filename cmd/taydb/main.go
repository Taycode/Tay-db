// Command taydb is an interactive shell over a Tay-db log file.
//
//	go run ./cmd/taydb data.log
//	> set user:1 Tay
//	> get user:1
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Taycode/Tay-db/store"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: taydb <log file>")
		os.Exit(2)
	}

	db, err := store.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Printf("loaded %d keys from %s\n", db.Len(), os.Args[1])

	lines := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for lines.Scan() {
		run(db, lines.Text())
		fmt.Print("> ")
	}
	if err := lines.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "stopped reading input:", err)
	}
}

func run(db *store.Store, line string) {
	command, args, _ := strings.Cut(strings.TrimSpace(line), " ")
	args = strings.TrimSpace(args)
	switch command {
	case "set":
		key, value, ok := strings.Cut(args, " ")
		if !ok {
			fmt.Println("usage: set <key> <value>")
			return
		}
		if err := db.Set(key, []byte(strings.TrimSpace(value))); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("ok")
	case "get":
		value, err := db.Get(args)
		if errors.Is(err, store.ErrNotFound) {
			fmt.Println("(not found)")
			return
		}
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(string(value))
	case "":
	default:
		fmt.Println("commands: set <key> <value> | get <key>")
	}
}
