package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type ALUInstruction struct {
	op         string
	left       string // left hand of the operation, always a register
	right      string // right hand side of the operation if it's a register (blank for values)
	rightValue int64
}

func numberToArray(num int64) []int {
	result := make([]int, 14)
	for i := 0; i < 14; i++ {
		result[len(result)-i-1] = int(num%9 + 1)
		num /= 9
	}
	return result
}

func arrayToNumber(a []int) int64 {
	result := int64(0)
	for i := 0; i < len(a); i++ {
		result = (result * 9) + int64(a[i]-1)
	}
	return result
}

type RegisterMap map[string]int64

// Run the program, return values of registers
func RunProgram(program []ALUInstruction, input []int, registers RegisterMap) (RegisterMap, bool) {
	rval := func(instr ALUInstruction) int64 {
		if instr.right != "" {
			return registers[instr.right]
		} else {
			return instr.rightValue
		}
	}

	for _, instr := range program {
		switch instr.op {
		case "inp":
			value := int64(input[0])
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
				instr.rightValue = int64(r)
			}
		}
		program = append(program, instr)
	}
	return program
}

func registersToString(r map[string]int64) string {
	return fmt.Sprintf("x: %3d y: %3d w: %3d z: %10d", r["x"], r["y"], r["w"], r["z"])
}

func TestProgram(program []ALUInstruction) {
	for j := 1; j <= 9; j++ {
		for i := 0; i < 1000; i++ {
			registers := make(RegisterMap)
			registers["z"] = int64(i)
			input := []int{j}

			if result, ok := RunProgram(program, input, registers); !ok {
				fmt.Println("ERROR with input", input, registers)
			} else {
				z := result["z"]
				if z == 0 {
					fmt.Printf("Z in: %3d Input: %v Z out: %v\n", i, input, z)
				}
			}
		}
	}

	r := make(RegisterMap)
	r["z"] = 1123343435
	fmt.Println(RunProgram(program, []int{9}, r))
}

type SolverResult struct {
	zValue int
	wValue int
}

func SolveProgram(program []ALUInstruction, wantZValue int64) []SolverResult {
	result := []SolverResult{}
	for w := 1; w <= 9; w++ {
		for i := 0; i < 1000; i++ {
			registers := make(RegisterMap)
			registers["z"] = int64(i)
			input := []int{w}

			if result, ok := RunProgram(program, input, registers); !ok {
				fmt.Println("ERROR with input", input, registers)
			} else {
				z := result["z"]
				if z == wantZValue {
					fmt.Printf("Z in: %3d Input: %v Z out: %v\n", i, input, z)
					result = append(result, SolverResult{z, w})
				}
			}
		}
	}
	return result
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

	for i := range programs {
		pos := len(programs) - i - 1
		res := SolveProgram(programs[pos], 0)
		fmt.Println(res)
		break
	}
}
