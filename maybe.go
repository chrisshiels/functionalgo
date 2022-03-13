package main


// host$ go1.18beta2 build maybe.go
// host$ ./maybe prideandprejudice.txt


import "fmt"
import "io/ioutil"
import "os"
import "regexp"
import "sort"
import "strings"


type Maybe[A any] struct {
    a *A
}


func Nothing[A any]() Maybe[A] {
    return Maybe[A] { nil }
}


func Just[A any](a A) Maybe[A] {
    return Maybe[A] { &a }
}


func Bind[A, B any](m Maybe[A], f func(a A) Maybe[B]) Maybe[B] {
    if m.a == nil {
        return Nothing[B]()
    }
    return f(*m.a)
}


func IsNothing[A any](m Maybe[A]) bool {
    return m.a == nil
}


func IsJust[A any](m Maybe[A]) bool {
    return m.a != nil
}


func FromJust[A any](m Maybe[A]) A {
    return *m.a
}


func Readfile(filename string) Maybe[string] {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return Nothing[string]()
    }
    return Just(string(bytes))
}


func Lowercase(s string) Maybe[string] {
    return Just(strings.ToLower(s))
}


func RemovePossessives(s string) Maybe[string] {
    return Just(regexp.MustCompile("'s").
                       ReplaceAllString(s, ""))
}


func RemoveNonAlphanumerics(s string) Maybe[string] {
    return Just(regexp.MustCompile("\\W").
                       ReplaceAllString(s, " "))
}


func Words(s string) Maybe[[]string] {
    return Just(strings.Fields(s))
}


func MapFrequencies(l []string) Maybe[map[string]int] {
    m := make(map[string]int)
    for _, s := range l {
        m[s] = m[s] + 1
    }
    return Just(m)
}


type Freq struct {
    s string
    count int
}


func ListFrequencies(m map[string]int) Maybe[[]Freq] {
    l := make([]Freq, len(m))
    i := 0
    for k, v := range m {
        l[i] = Freq{ k, v }
        i += 1
    }
    return Just(l)
}


type Freqs []Freq


func (f Freqs) Len() int {
    return len(f)
}


func (f Freqs) Swap(i, j int) {
    f[i], f[j] = f[j], f[i]
}


func (f Freqs) Less(i, j int) bool {
    if f[i].count != f[j].count {
        return f[i].count > f[j].count
    } else {
        return f[i].s < f[j].s
    }
}


func SortFrequencies(l []Freq) Maybe[[]Freq] {
    sort.Sort(Freqs(l))
    return Just(l)
}


func Output(l []Freq) Maybe[[]Freq] {
    for _, v := range l {
        fmt.Printf("%v %v\n", v.s, v.count)
    }
    return Just(l)
}


func main() {
    for _, a := range os.Args[1:] {
        m := Bind(Bind(Bind(Bind(Bind(Bind(Bind(Bind(Bind(
                 Just(a),
                 Readfile),
                 Lowercase),
                 RemovePossessives),
                 RemoveNonAlphanumerics),
                 Words),
                 MapFrequencies),
                 ListFrequencies),
                 SortFrequencies),
                 Output)
        if IsNothing(m) {
            fmt.Println("Error")
            os.Exit(1)
        }
    }
    os.Exit(0)
}
