package main

import (
	"flag"
	"fmt"
	"os"
)

func getenv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {
	apiAddr := flag.String("api", getenv("FYST_DB_ADDR", ":5900"), "DB API address")
	fmt.Println(*apiAddr)

	filescmd := flag.NewFlagSet("files", flag.ExitOnError)
	id := filescmd.String("id", "", "Files to query")

	mulcmd := flag.NewFlagSet("mul", flag.ExitOnError)
	aMul := mulcmd.Int("a", 0, "The value of a")
	bMul := mulcmd.Int("b", 0, "The value of b")

	switch os.Args[1] {
	case "add":
		addcmd.Parse(os.Args[2:])
		fmt.Println(*aAdd + *bAdd)
	case "mul":
		mulcmd.Parse(os.Args[2:])
		fmt.Println(*(aMul) * (*bMul))
	default:
		fmt.Println("expected add or mul command")
		os.Exit(1)
	}

}
