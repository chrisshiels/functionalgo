package main


// host$ go1.18beta2 build either.go


import "fmt"
import "io/ioutil"
import "os"
import "regexp"
import "sort"
import "strings"


type Either[A, B any] struct {
    a *A
    b *B
}


func Left[A, B any](a A) Either[A, B] {
    return Either[A, B] { &a, nil }
}


func Right[A, B any](b B) Either[A, B] {
    return Either[A, B] { nil, &b }
}


func Bind[A, B, C any](e Either[A, B], f func(b B) Either[A, C]) Either[A, C] {
    if e.a != nil {
        return Left[A, C](*e.a)
    }
    return f(*e.b)
}


func IsLeft[A, B any](e Either[A, B]) bool {
    return e.a != nil
}


func IsRight[A, B any](e Either[A, B]) bool {
    return e.a == nil
}


func FromLeft[A, B any](e Either[A, B]) A {
    return *e.a
}


func FromRight[A, B any](e Either[A, B]) B {
    return *e.b
}


func Readfile(filename string) Either[error, string] {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return Left[error, string](err)
    }
    return Right[error, string](string(bytes))
}


func Lowercase(s string) Either[error, string] {
    return Right[error, string](strings.ToLower(s))
}


func RemovePossessives(s string) Either[error, string] {
    return Right[error, string](regexp.MustCompile("'s").
                                       ReplaceAllString(s, ""))
}


func RemoveNonAlphanumerics(s string) Either[error, string] {
    return Right[error, string](regexp.MustCompile("\\W").
                                       ReplaceAllString(s, " "))
}


func Words(s string) Either[error, []string] {
    return Right[error, []string](strings.Fields(s))
}


func MapFrequencies(l []string) Either[error, map[string]int] {
    m := make(map[string]int)
    for _, s := range l {
        m[s] = m[s] + 1
    }
    return Right[error, map[string]int](m)
}


type Freq struct {
    s string
    count int
}


func ListFrequencies(m map[string]int) Either[error, []Freq] {
    l := make([]Freq, len(m))
    i := 0
    for k, v := range m {
        l[i] = Freq{ k, v }
        i += 1
    }
    return Right[error, []Freq](l)
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


func SortFrequencies(l []Freq) Either[error, []Freq] {
    sort.Sort(Freqs(l))
    return Right[error, []Freq](l)
}


func Output(l []Freq) Either[error, []Freq] {
    for _, v := range l {
        fmt.Printf("%v %v\n", v.s, v.count)
    }
    return Right[error, []Freq](l)
}


func main() {
    for _, a := range os.Args[1:] {
        e := Bind(Bind(Bind(Bind(Bind(Bind(Bind(Bind(Bind(
                 Right[error, string](a),
                 Readfile),
                 Lowercase),
                 RemovePossessives),
                 RemoveNonAlphanumerics),
                 Words),
                 MapFrequencies),
                 ListFrequencies),
                 SortFrequencies),
                 Output)
        if (IsLeft(e)) {
            err := FromLeft(e)
            fmt.Printf("Error:  %s\n", err)
            os.Exit(1)
        }
    }
    os.Exit(0)
}
