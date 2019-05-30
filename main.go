package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

const Sum = 0
const Sub = 1
const Mult = 2
const rowsToHell = 25
const gameTime = 5 // will be diveded by 1 second

type Expression struct {
	Column int
	Row    float32

	String string
	Result int
}

func resolve(a, b, opr int) int {
	switch opr {
	case Sum:
		return a + b
	case Sub:
		return a - b
	case Mult:
		return a * b
	}
	return a + b
}

func operatorToString(opr int) string {
	switch opr {
	case Sum:
		return "+"
	case Sub:
		return "-"
	case Mult:
		return "*"
	}
	return "+"
}

func NewExpression() *Expression {
	opr := rand.Intn(3)
	a := 0
	b := 0
	if opr == Mult {
		a = rand.Intn(9) + 1
		b = rand.Intn(9) + 1
	} else {
		a = rand.Intn(98) + 1
		b = rand.Intn(98) + 1
	}

	e := Expression{
		Column: rand.Intn(6) + 1,
		Row:    2,
		String: strconv.Itoa(a) + operatorToString(opr) + strconv.Itoa(b),
		Result: resolve(a, b, opr)}

	return &e
}

func configSignals(score *int, seed int64) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGQUIT, os.Interrupt)
	go func() {
		<-sigc
		fmt.Println("")
		fmt.Println("score:", *score, "- seed", seed)
		os.Exit(0)
	}()
}

func drawHell(fire string) {
	tm.MoveCursor(0, rowsToHell)
	tm.Println(fire)
}

func drawInput(g string) {
	tm.Print("answer:")

	var gsize int
	if g == "" {
		gsize = 0
	} else {
		gsize = len(g)
	}
	if gsize > 0 {
		tm.Print(g)
	}

	for i := 0; i < 5-gsize; i++ {
		tm.Print("_")
	}

	tm.MoveCursorBackward(5 - gsize)

}

func hell() string {
	firec := "🔥"
	fire := "🔥"
	for i := 0; i < tm.Width()-1; i++ {
		fire += firec
	}
	return fire
}

func getDif() int {
	return 5
}

func generateExpression(dif int, ee []Expression) []Expression {
	r := rand.Intn(500)
	if dif/gameTime >= r {
		ne := NewExpression()
		return append(ee, *ne)
	}

	return ee
}
func hasExpressionsInitialRow(ee []Expression) bool {
	for i := 0; i < len(ee); i++ {
		if int(ee[i].Row) == 2 || int(ee[i].Row) == 3 {
			return true
		}
	}
	return false
}

func getSeed() int64 {
	if len(os.Args) == 2 {
		v, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic("Strange seed")
		}

		return int64(v)
	}

	return time.Now().Unix()
}

func main() {
	score := 0
	seed := getSeed()
	rand.Seed(seed)
	guess := ""
	h := hell()
	dif := getDif()

	configSignals(&score, seed)
	go listenKey(&guess)
	exp := NewExpression()
	ee := make([]Expression, 0)
	ee = append(ee, *exp)

	for {
		if !hasExpressionsInitialRow(ee) {
			ee = generateExpression(dif, ee)
			dif++
		}

		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Println("QuickMath 1.0 - Seed:", seed)
		drawExpressions(ee)
		drawHell(h)
		drawInput(guess)
		tm.Flush()
		time.Sleep(time.Second / 5)
		ee = verifyAnwser(&guess, &score, &ee)
		gameover(ee, score, seed)
		pushExpressions(ee)
	}
}

func gameover(ee []Expression, score int, seed int64) {
	for i := 0; i < len(ee); i++ {
		if int(ee[i].Row) == rowsToHell-1 {
			fmt.Println("")
			fmt.Println("Game Over")
			fmt.Println("score:", score, "- seed", seed)
			os.Exit(0)
		}
	}
}

func pushExpressions(ee []Expression) {
	for i := 0; i < len(ee); i++ {
		ee[i].Row += 0.2
	}
}

func drawExpressions(ee []Expression) {
	for i := 0; i < len(ee); i++ {
		tm.MoveCursor(ee[i].Column*10, int(ee[i].Row))
		tm.Print(ee[i].String)
	}
}

func verifyAnwser(g *string, score *int, ee *[]Expression) []Expression {
	nee := make([]Expression, 0)
	oee := *ee
	correct := false
	for i := 0; i < len(oee); i++ {
		ig, _ := strconv.Atoi(*g)
		if oee[i].Result != ig {
			nee = append(nee, oee[i])
		} else {
			correct = true
			*score++
		}
	}

	if correct {
		*g = ""
	}

	return nee
}

func listenKey(guess *string) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	allowedKeys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-"}
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			*guess = ""
		} else if key == keyboard.KeyEnter {
			*guess = ""
		}

		if contains(allowedKeys, string(char)) {
			*guess += string(char)
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
