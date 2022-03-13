package main


// host$ go1.18beta2 build higherorder.go


import "fmt"
import "os"
import "strconv"


func Range(m, n int) []int {
    l := make([]int, n - m)
    for i := 0; m < n; i++ {
        l[i] = m
        m++
    }
    return l
}


// (a -> b) -> [a] -> [b].
func Map[A, B any](f func (e A) B, l []A) []B {
    l1 := make([]B, len(l))
    i := 0
    for _, e := range(l) {
        l1[i] = f(e)
        i++
    }
    return l1
}


// (a -> Bool) -> [a] -> [a].
func Filter[A any](f func (e A) bool, l []A) []A {
    l1 := make([]A, len(l))
    i := 0
    for _, e := range(l) {
        if f(e) {
            l1[i] = e
            i++
        }
    }
    return l1[:i]
}


// (a -> b -> a) -> a -> [b] -> a.
func Reduce[A, B any](f func (a A, e B) A, v A, l []B) A {
    a := v
    for _, e := range(l) {
        a = f(a, e)
    }
    return a
}


func main() {
    l := Range(1, 11)
    fmt.Printf("%T %v\n", l, l)
    // => []int [1 2 3 4 5 6 7 8 9 10]


    l1 := Map(func(e int) string {
                  return strconv.Itoa(e)
              },
              l)
    fmt.Printf("%T %v\n", l1, l1)
    // => []string [1 2 3 4 5 6 7 8 9 10]


    l2 := Filter(func(e int) bool {
                     return e % 2 == 0
                 },
                 l)
    fmt.Printf("%T %v\n", l2, l2)
    // => []int [2 4 6 8 10]


    v := Reduce(func(a, e int) int {
                    return a + e
                },
                0,
                l)
    fmt.Printf("%T %v\n", v, v)
    // => int 55


    os.Exit(0)
}
