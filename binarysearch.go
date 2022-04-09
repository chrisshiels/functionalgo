package main


// host$ go build binarysearch.go
// host$ ./binarysearch


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


func binarysearch[A Ordered](compare func (x, y A) int, l []A, v A,
                             i int, j int) int {
    if i == j {
        return -1
    } else {
        k := (i + j) / 2
        switch (compare(v, l[k])) {
            case -1:
                return binarysearch(compare, l, v, i, k)
            case 0:
                return k
            default:
                return binarysearch(compare, l, v, k + 1, j)
        }
    }
}


func BinarySearch[A Ordered](compare func (x, y A) int, l []A, v A) int {
    return binarysearch(compare, l, v, 0, len(l))
}


func main() {
    fmt.Printf("%v\n",
               BinarySearch(Compare[int],
                            []int{ 10, 20, 30, 40, 50, 60, 70, 80, 90, 100 },
                            70))
    // => 6


    os.Exit(0)
}
