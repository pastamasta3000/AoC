//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func slicesFilter[T any](ts []T, f func(T) bool) []T {
	us := make([]T, 0, len(ts))
	for _, v := range ts {
		if f(v) {
			us = append(us, v)
		}
	}
	return us
}

// run f for all elements and return matches
func gridFind[T any](grd [][]T, f func(T) bool) [][2]int {
	var res [][2]int
	for i, r := range grd {
		for j, c := range r {
			if f(c) {
				res = append(res, [2]int{i, j})
			}
		}
	}
	return res
}

// find values f accepts in kernel around i, j
func gridFindInRadiusCoords[T any](grd [][]T, krnl [][2]int, i, j int, f func(T, int, int) bool) [][2]int {
	var res [][2]int
	r, c := len(grd), len(grd[0])
	for _, o := range krnl {
		io, jo := i+o[0], j+o[1]
		if io < 0 || io >= r || jo < 0 || jo >= c {
			continue
		}
		if f(grd[io][jo], io, jo) {
			res = append(res, [2]int{io, jo})
		}
	}
	return res
}

func gridCrawlFromEdgesRecr[T any](grd [][]T, krnl [][2]int, i, j int, f func([][]T, int, int) bool) {
	if f(grd, i, j) {
		return
	}
	r, c := len(grd), len(grd[0])
	for _, o := range krnl {
		io, jo := i+o[0], j+o[1]
		if io < 0 || io >= r || jo < 0 || jo >= c {
			continue
		}
		gridCrawlFromEdgesRecr(grd, krnl, io, jo, f)
	}
}

func gridCrawlFromEdges[T any](grd [][]T, krnl [][2]int, f func([][]T, int, int) bool) {
	r, c := len(grd), len(grd[0])
	for i := 0; i < r; i++ {
		gridCrawlFromEdgesRecr(grd, krnl, i, 0, f)
		gridCrawlFromEdgesRecr(grd, krnl, i, c-1, f)
	}
	for j := 0; j < c; j++ {
		gridCrawlFromEdgesRecr(grd, krnl, 0, j, f)
		gridCrawlFromEdgesRecr(grd, krnl, r-1, j, f)
	}
}

func gridGrow[T any](grd [][]T) [][]T {
	var ngrd [][]T
	or, oc := len(grd), len(grd[0])
	nr, nc := (or*2)+1, (oc*2)+1
	ngrd = make([][]T, nr)
	for i := range ngrd {
		ngrd[i] = make([]T, nc)
	}
	for i := 1; i < nr; i += 2 {
		for j := 1; j < nc; j += 2 {
			ngrd[i][j] = grd[i/2][j/2]
		}
	}
	return ngrd
}

func gridShrink[T any](grd [][]T) [][]T {
	var ngrd [][]T
	or, oc := len(grd), len(grd[0])
	nr, nc := or/2, oc/2
	ngrd = make([][]T, nr)
	for i := range ngrd {
		ngrd[i] = make([]T, nc)
	}
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			ngrd[i][j] = grd[(i*2)+1][(j*2)+1]
		}
	}
	return ngrd
}

func gridPrint(grd [][]int) {
	for _, r := range grd {
		for _, c := range r {
			fmt.Printf("%d", c)
		}
		fmt.Println()
	}
}

func gridPrintBytes(grd [][]byte) {
	for _, r := range grd {
		for _, c := range r {
			if c == 0 {
				c = ' '
			}
			fmt.Printf("%s", string(c))
		}
		fmt.Println()
	}
}

var (
	dirn = [2]int{-1, 0}
	dirs = [2]int{1, 0}
	dire = [2]int{0, 1}
	dirw = [2]int{0, -1}
)

var dirsList = [4][2]int{dirn, dirs, dire, dirw}

var pipeEndsStr = map[byte]([2][2]int){
	'|': {dirn, dirs},
	'-': {dire, dirw},
}
var pipeEnds = map[byte]([2][2]int){
	'|': {dirn, dirs},
	'-': {dire, dirw},
	'L': {dirn, dire},
	'J': {dirn, dirw},
	'7': {dirw, dirs},
	'F': {dire, dirs},
}

func main() {
	var srf [][]byte
	var dists [][]int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		srf = append(srf, slices.Clone(scanner.Bytes()))
	}
	srf = gridGrow(srf)
	r, c := len(srf), len(srf[0])
	for i := range srf {
		for j := range srf[i] {
			if srf[i][j] != 0 {
				continue
			}
			for p, os := range pipeEndsStr {
				f := false
				for on := range os {
					ip, jp := i+os[on][0], j+os[on][1]
					if ip < 0 || ip >= r || jp < 0 || jp >= c {
						continue
					}
					for _, po := range pipeEnds[srf[ip][jp]] {
						ir, jr := ip+po[0], jp+po[1]
						if ir == i && jr == j {
							f = true
						}
						if f {
							break
						}
					}
					if f {
						break
					}
				}
				if f {
					srf[i][j] = p
					break
				}
			}
		}
	}
	dists = make([][]int, len(srf))
	for i := range srf {
		dists[i] = make([]int, len(srf[i]))
	}
	pos := gridFind(srf, func(b byte) bool { return b == 'S' })[0]
	dists[pos[0]][pos[1]] = 1
	paths := gridFindInRadiusCoords(srf, dirsList[:], pos[0], pos[1], func(b byte, i, j int) bool {
		f := false
		for _, conv := range pipeEnds[b] {
			if pos == [2]int{i + conv[0], j + conv[1]} {
				f = true
				break
			}
		}
		return f
	})
	for _, v := range paths {
		prev, cur := pos, v
		dir := [2]int{prev[0] - cur[0], prev[1] - cur[1]}
		steps := 1
		for cur != pos {
			dists[cur[0]][cur[1]] = 1
			piece := srf[cur[0]][cur[1]]
			m := -1
			for n, end := range pipeEnds[piece] {
				if end == dir {
					m = n
				}
			}
			if m >= 0 {
				// choose the other end
				oend := pipeEnds[piece][1-m]
				prev[0], prev[1] = cur[0], cur[1]
				cur[0], cur[1] = cur[0]+oend[0], cur[1]+oend[1]
				dir[0], dir[1] = prev[0]-cur[0], prev[1]-cur[1]
				steps++
			} else {
				// dead end
				break
			}
		}
	}
	gridCrawlFromEdges(dists, dirsList[:], func(grd [][]int, i, j int) bool {
		if grd[i][j] == 1 || grd[i][j] == 2 {
			return true
		}
		grd[i][j] = 2
		return false
	})
	dists = gridShrink(dists)
	dmax := gridFind(dists, func(n int) bool { return n == 0 })
	fmt.Println(len(dmax))
}
