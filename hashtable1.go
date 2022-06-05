package main


// https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/HashIntro.html


import "fmt"
import "os"


type HashTableEntry[K, V comparable] struct {
    k K
    v V
}


type HashTable[K, V comparable] struct {
    nentries int
    nslots int
    slots [][]HashTableEntry[K, V]
    hash func(k K, m int) int
}


func HashTableNew[K, V comparable](nslots int,
                                   hash func(k K, m int) int) *HashTable[K, V] {
    h := new(HashTable[K, V])
    h.nentries = 0
    h.nslots = nslots
    h.slots = make([][]HashTableEntry[K, V], nslots)
    for i := range h.slots {
        h.slots[i] = make([]HashTableEntry[K, V], 0)
    }
    h.hash = hash
    return h
}


func HashTableGet[K, V comparable](h *HashTable[K, V], k K) (V, bool) {
    slot := h.hash(k, h.nslots)
    for _, e := range h.slots[slot] {
        if e.k == k {
            return e.v, true
        }
    }
    return *new(V), false
}


func HashTableKeys[K, V comparable](h *HashTable[K, V]) []K {
    keys := make([]K, 0, h.nentries)
    for i, _ := range h.slots {
        for _, e := range h.slots[i] {
            keys = append(keys, e.k)
        }
    }
    return keys
}


func HashTableLoadFactor[K, V comparable](h *HashTable[K, V]) float32 {
    return float32(h.nentries) / float32(h.nslots)
}


func remove[T any](v []T, i int) []T {
    v[i] = v[len(v) - 1]
    v = v[:len(v) - 1]
    return v
}


func hashtableremove[K, V comparable](h *HashTable[K, V],
                                      k K) *HashTable[K, V] {
    slot := h.hash(k, h.nslots)
    for i, e := range h.slots[slot] {
        if e.k == k {
            h.slots[slot] = remove(h.slots[slot], i)
            h.nentries--
            return h
        }
    }
    return h
}


func HashTableRemove[K, V comparable](h *HashTable[K, V],
                                      k K) *HashTable[K, V] {
    h = hashtableremove(h, k)
    if HashTableLoadFactor(h) < 0.3 {
        h = HashTableResize(h, h.nslots / 2)
    }
    return h
}


func HashTableResize[K, V comparable](h *HashTable[K, V],
                                      nslots int) *HashTable[K, V] {
    h2 := HashTableNew[K, V](nslots, h.hash)
    for i, _ := range h.slots {
        for _, e := range h.slots[i] {
            h2 = hashtableset(h2, e.k, e.v)
        }
    }
    return h2
}


func hashtableset[K, V comparable](h *HashTable[K, V],
                                   k K, v V) *HashTable[K, V] {
    slot := h.hash(k, h.nslots)
    for i, e := range h.slots[slot] {
        if e.k == k {
            h.slots[slot][i].v = v
            return h
        }
    }
    h.slots[slot] = append(h.slots[slot], HashTableEntry[K, V]{ k, v })
    h.nentries++
    return h
}


func HashTableSet[K, V comparable](h *HashTable[K, V],
                                   k K, v V) *HashTable[K, V] {
    h = hashtableset(h, k, v)
    if HashTableLoadFactor(h) > 0.7 {
        h = HashTableResize(h, h.nslots * 2)
    }
    return h
}


func hashsumchars(s string, m int) int {
    sum := 0
    for _, e := range s {
        sum += int(e)
    }
    return sum % m
}


func main() {
    elements := []string{ "hydrogen",
                          "helium",
                          "lithium",
                          "beryllium",
                          "boron",
                          "carbon",
                          "nitrogen",
                          "oxygen",
                          "fluorine",
                          "neon",
                          "sodium",
                          "magnesium",
                          "aluminium",
                          "silicon",
                          "phosphorus",
                          "sulfur",
                          "chlorine",
                          "argon",
                          "potassium",
                          "calcium" }
    hashtable := HashTableNew[string, int](1, hashsumchars)
    fmt.Printf("%v load %v\n", hashtable, HashTableLoadFactor(hashtable))


    for _, e := range elements {
        hashtable = HashTableSet(hashtable, e, len(e))
        fmt.Printf("%v load %v\n", hashtable, HashTableLoadFactor(hashtable))
    }

    v, ok := HashTableGet(hashtable, "oxygen")
    fmt.Printf("%v %v\n", v, ok)


    for _, e := range HashTableKeys(hashtable) {
        hashtable = HashTableRemove(hashtable, e)
        fmt.Printf("%v load %v\n", hashtable, HashTableLoadFactor(hashtable))
    }

    v, ok = HashTableGet(hashtable, "oxygen")
    fmt.Printf("%v %v\n", v, ok)


    vs := []int{ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 }
    vs = remove(vs, 0)
    vs = remove(vs, 1)
    fmt.Printf("%v\n", vs)


    h := HashTableNew[string, int](1, hashsumchars)
    h2 := HashTableSet(h, "tom", 1)
    h3 := HashTableSet(h2, "dick", 2)
    h4 := HashTableSet(h3, "harry", 3)
    h5 := HashTableSet(h4, "harry", 4)
    h6 := HashTableRemove(h5, "harry")
    h7 := HashTableRemove(h6, "dick")
    h8 := HashTableRemove(h7, "tom")
    fmt.Printf("%p %v load %v\n", h, h, HashTableLoadFactor(h))
    fmt.Printf("%p %v load %v\n", h2, h2, HashTableLoadFactor(h2))
    fmt.Printf("%p %v load %v\n", h3, h3, HashTableLoadFactor(h3))
    fmt.Printf("%p %v load %v\n", h4, h4, HashTableLoadFactor(h4))
    fmt.Printf("%p %v load %v\n", h5, h5, HashTableLoadFactor(h5))
    fmt.Printf("%p %v load %v\n", h6, h6, HashTableLoadFactor(h6))
    fmt.Printf("%p %v load %v\n", h7, h7, HashTableLoadFactor(h7))
    fmt.Printf("%p %v load %v\n", h8, h8, HashTableLoadFactor(h8))


    os.Exit(0)
}
