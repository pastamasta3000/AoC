//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func slicesMap[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

// homegrown impl of maps.Clear until it's standardised
func mapsClear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// pattern matching NFA that runs in O(n)
func countPossible(s []byte, c []int) int {
	pos := 0
	// state is a tuple of 4 values
	cstates := map[[4]int]int{{0, 0, 0, 0}: 1}
	nstates := map[[4]int]int{}
	for len(cstates) > 0 {
		for state, num := range cstates {
			si, ci, cc, expdot := state[0], state[1], state[2], state[3]
			if si == len(s) {
				if ci == len(c) {
					pos += num
				}
				continue
			}
			switch {
			case (s[si] == '#' || s[si] == '?') && ci < len(c) && expdot == 0:
				// we are still looking for broken springs
				if s[si] == '?' && cc == 0 {
					// we are not in a run of broken springs, so ? can be working
					nstates[[4]int{si + 1, ci, cc, expdot}] += num
				}
				cc++
				if cc == c[ci] {
					// we've found the full next contiguous section of broken springs
					ci++
					cc = 0
					expdot = 1 // we only want a working spring next
				}
				nstates[[4]int{si + 1, ci, cc, expdot}] += num
			case (s[si] == '.' || s[si] == '?') && cc == 0:
				// we are not in a contiguous run of broken springs
				expdot = 0
				nstates[[4]int{si + 1, ci, cc, expdot}] += num
			}
		}
		cstates, nstates = nstates, cstates
		mapsClear(nstates)
	}
	return pos
}

func main() {
	paths := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := scanner.Text()
		b, a, _ := strings.Cut(str, " ")
		bn, an := "", ""
		for i := 0; i < 5; i++ {
			bn, an = bn+b+"?", an+a+","
		}
		b, a = strings.TrimSuffix(bn, "?"), strings.TrimSuffix(an, ",")
		s := []byte(b)
		c := slicesMap(strings.Split(a, ","), atoi)
		p := countPossible(s, c)
		paths += p
	}
	fmt.Println(paths)
}
