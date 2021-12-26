package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type ALUInstruction struct {
	op         string
	left       string // left hand of the operation, always a register
	right      string // right hand side of the operation if it's a register (blank for values)
	rightValue int
}

func numberToArray(num int) []int {
	result := make([]int, 14)
	for i := 0; i < 14; i++ {
		result[len(result)-i-1] = int(num%9 + 1)
		num /= 9
	}
	return result
}

func arrayToNumber(a []int) int {
	result := 0
	for i := 0; i < len(a); i++ {
		result = (result * 9) + (a[i] - 1)
	}
	return result
}

type RegisterMap map[string]int

// Run the program, return values of registers
func RunProgram(program []ALUInstruction, input []int, registers RegisterMap) (RegisterMap, bool) {
	rval := func(instr ALUInstruction) int {
		if instr.right != "" {
			return registers[instr.right]
		} else {
			return instr.rightValue
		}
	}

	for _, instr := range program {
		switch instr.op {
		case "inp":
			value := input[0]
			input = input[1:]
			registers[instr.left] = value
		case "add":
			registers[instr.left] += rval(instr)
		case "mul":
			registers[instr.left] *= rval(instr)
		case "div":
			rv := rval(instr)
			if rv == 0 {
				fmt.Println("MONAD error: division by zero", instr)
				return registers, false
			}
			registers[instr.left] /= rv
		case "mod":
			rv := rval(instr)
			if rv <= 0 || registers[instr.left] < 0 {
				fmt.Println("MONAD error: modulo args", instr)
				return registers, false
			}
			registers[instr.left] %= rval(instr)
		case "eql":
			if registers[instr.left] == rval(instr) {
				registers[instr.left] = 1
			} else {
				registers[instr.left] = 0
			}
		default:
			panic(fmt.Sprintf("invalid operation: %+v", instr))
		}
	}
	return registers, true
}

func ReadProgram(input []string) []ALUInstruction {
	program := []ALUInstruction{}

	for _, line := range input {
		if line == "" {
			continue
		}
		s := strings.Split(line, " ")
		instr := ALUInstruction{}
		instr.op = s[0]
		if s[0] == "inp" {
			instr.left = s[1]
		} else {
			instr.left = s[1]
			if unicode.IsLetter(rune(s[2][0])) {
				instr.right = s[2]
			} else {
				r, _ := strconv.Atoi(s[2])
				instr.rightValue = r
			}
		}
		program = append(program, instr)
	}
	return program
}

type SolverResult struct {
	z int
	w int
}

func SolveProgram(program []ALUInstruction, wantZValue int) []SolverResult {
	solutions := []SolverResult{}
	maxZValue := wantZValue*26 + 100
	minZValue := wantZValue * 26
	if minZValue > 100 {
		minZValue -= 100
	}
	for w := 1; w <= 9; w++ {
		for z := minZValue; z <= maxZValue; z++ {
			registers := make(RegisterMap)
			registers["z"] = z
			input := []int{w}

			if result, ok := RunProgram(program, input, registers); !ok {
				fmt.Println("ERROR with input", input, registers)
			} else {
				zResult := result["z"]
				if zResult == wantZValue {
					//fmt.Printf("Z in: %3d Input: %v Z out: %v\n", z, input, zResult)
					solutions = append(solutions, SolverResult{z, w})
					break
				}
			}
		}
	}

	sort.Slice(solutions, func(i, j int) bool { return solutions[i].w > solutions[j].w })
	return solutions
}

func Solve(programs [][]ALUInstruction, forValue int, digits []int) bool {
	if len(programs) == 0 {
		fmt.Println("Solution:", digits)
		return true
	}

	pos := len(programs) - 1
	res := SolveProgram(programs[pos], forValue)

	if len(res) > 0 {
		fmt.Printf("Program %d solved for %d: %+v\n", pos, forValue, res)

		for _, r := range res {
			digits = append(digits, r.w)
			if Solve(programs[:len(programs)-1], r.z, digits) {
				return true
			}
			digits = digits[:len(digits)-1]
		}
	}

	return false
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	bigProgram := ReadProgram(strings.Split(string(buf), "\n"))

	programs := make([][]ALUInstruction, 0)
	var currentProgram []ALUInstruction

	for _, instr := range bigProgram {
		if instr.op == "inp" {
			if currentProgram != nil {
				programs = append(programs, currentProgram)
			}
			currentProgram = make([]ALUInstruction, 0)
		}
		currentProgram = append(currentProgram, instr)
	}
	programs = append(programs, currentProgram)

	//TestProgram(programs)
	//return

	r := []int{}
	if Solve(programs, 0, r) {
		fmt.Println("Answer =", r)
	} else {
		fmt.Println("No answer.")
	}
}
