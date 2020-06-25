package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
)

type polar struct {
	radius float64
	θ      float64
}

type cartesian struct {
	x float64
	y float64
}

var prompt = "Enter a radius and angle(in degree), eg. 5 30, " + "or %s to quit."

func init() {
	if runtime.GOOS == "windows" {
		prompt = fmt.Sprintf(prompt, "CRTL + Z, Enter")
	} else {
		prompt = fmt.Sprintf(prompt, "CRTL +D")
	}
}

func main() {
	question := make(chan polar)
	defer close(question)

	answer := createSolver(question)
	defer close(answer)

	interact(question, answer)
}

func createSolver(question chan polar) chan cartesian {
	answer := make(chan cartesian)
	go func() {
		for {
			polarcord := <-question
			θ := polarcord.θ * math.Pi
			x := polarcord.radius * math.Cos(θ)
			y := polarcord.radius * math.Sin(θ)

			answer <- cartesian{x, y}
		}
	}()
	return answer
}

const result = "Polar radius = %.02f θ = %.02f° → Cartesian x = %.02f y = %.02f\n"

func interact(question chan polar, answer chan cartesian) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		fmt.Print("Radius and angle: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			break
		}
		var radius, θ float64

		if _, err = fmt.Sscanf(line, "%f %f", &radius, &θ); err != nil {
			fmt.Fprintln(os.Stderr, "invalid input")
			continue
		}
		question <- polar{radius, θ}
		cord := <-answer
		fmt.Printf(result, radius, θ, cord.x, cord.y)
	}
	fmt.Println()
}
