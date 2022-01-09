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

func numberToArray(num int64, n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[len(result)-i-1] = int(num%9 + 1)
		num /= 9
	}
	return result
}

func formatNumber(n int64) string {
	a := numberToArray(n, 14)
	res := ""
	for _, v := range a {
		res += fmt.Sprintf("%d", v)
	}
	return res
}

func arrayToNumber(a []int) int64 {
	result := int64(0)
	for i := 0; i < len(a); i++ {
		result = (result * 9) + int64(a[i]-1)
	}
	return result
}

type RegisterMap map[string]int64

func RunWithZW(program []ALUInstruction, z int64, w int) int64 {
	res, ok := RunProgram(program, []int{w}, map[string]int64{"z": z})
	if !ok {
		panic("program failed")
	}
	return res["z"]
}

// Run the program, return values of registers
func RunProgram(program []ALUInstruction, input []int, registers RegisterMap) (RegisterMap, bool) {
	rval := func(instr ALUInstruction) int64 {
		if instr.right != "" {
			return registers[instr.right]
		} else {
			return int64(instr.rightValue)
		}
	}

	for _, instr := range program {
		switch instr.op {
		case "inp":
			value := input[0]
			input = input[1:]
			registers[instr.left] = int64(value)
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

func Pow(a int64, n int) int64 {
	res := int64(1)
	for ; n > 0; n-- {
		res *= a
	}
	return res
}

type ZKey struct {
	digit int
	z     int64
}

type ZVal struct {
	min, max int64
}

type ZMap map[ZKey]ZVal

func Solve(programs [][]ALUInstruction, findSmallestSerial bool) (int64, int64) {
	zValues := make(ZMap)
	zValues[ZKey{}] = ZVal{}

	result := []ZVal{}

	for p := range programs {
		maxZ := Pow(26, len(programs)-p) // Overflow on 1st run, but that's OK
		newZValues := make(ZMap)
		filtered := 0
		maxValue := int64(0)

		for digit := 1; digit <= 9; digit++ {
			for k, v := range zValues {
				rz := RunWithZW(programs[p], k.z, digit)
				if rz < maxZ {
					if rz > maxValue {
						maxValue = rz
					}
					minVal := v.min*9 + int64(digit-1)
					maxVal := v.max*9 + int64(digit-1)
					key := ZKey{digit, rz}
					if cachedVal, ok := newZValues[key]; ok {
						if cachedVal.min < minVal {
							minVal = cachedVal.min
						}
						if cachedVal.max > maxVal {
							maxVal = cachedVal.max
						}
					}
					val := ZVal{minVal, maxVal}
					newZValues[key] = val

					if p == 13 && rz == 0 {
						result = append(result, val)
					}
				} else {
					filtered++
				}
			}
		}

		fmt.Println("Program", p, "zvalues=", len(newZValues), "maxZ", maxZ, "filtered", filtered, "maxVal", maxValue)
		zValues = newZValues
	}

	sort.Slice(result, func(i, j int) bool { return result[i].min < result[j].min })
	minVal := result[0].min
	sort.Slice(result, func(i, j int) bool { return result[i].max > result[j].max })
	maxVal := result[0].max

	return minVal, maxVal
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

	smallest, largest := Solve(programs, false)
	fmt.Println("Smallest =", formatNumber(smallest), "Largest =", formatNumber(largest))
}
