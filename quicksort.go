package main


// host$ go build quicksort.go
// host$ ./quicksort


import "fmt"
import "os"


type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~uintptr | ~float32 | ~float64 | ~string
}


func Compare[A Ordered](x, y A) int {
    booltoint := func (b bool) int {
        if b {
            return 1
        } else {
            return 0
        }
    }

    return booltoint(x > y) - booltoint(x < y)
}


func QuickSort[A Ordered](compare func (x, y A) int, l []A) []A {
    if len(l) == 0 || len(l) == 1 {
        return l
    }

    pivot := l[0]
    left := make([]A, 0, len(l) / 2)
    right := make([]A, 0, len(l) / 2)

    for _, e := range(l[1:]) {
        if compare(e, pivot) <= 0 {
            left = append(left, e)
        } else {
            right = append(right, e)
        }
    }

    l1 := make([]A, 0, len(l))
    l1 = append(l1, QuickSort(compare, left)...)
    l1 = append(l1, pivot)
    l1 = append(l1, QuickSort(compare, right)...)
    return l1
}


func main() {
    l1 := []int{ 3, 1, 4, 1, 5, 9, 2, 6, 5, 4 }
    fmt.Printf("%v\n", QuickSort(Compare[int], l1))
    // => [1 1 2 3 4 4 5 5 6 9]


    l2 := []int{ 3 }
    fmt.Printf("%v\n", QuickSort(Compare[int], l2))
    // => [3]


    l3 := []int{}
    fmt.Printf("%v\n", QuickSort(Compare[int], l3))
    // => []


    os.Exit(0)
}
