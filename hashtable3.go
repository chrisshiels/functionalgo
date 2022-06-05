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


func (h *HashTable[K, V]) hashtablenewlayer(slot int) *HashTable[K, V] {
    h2 := new(HashTable[K, V])
    h2.nentries = h.nentries
    h2.nslots = h.nslots
    h2.slots = make([][]HashTableEntry[K, V], h2.nslots)
    copy(h2.slots, h.slots)
    h2.slots[slot] = make([]HashTableEntry[K, V], len(h.slots[slot]))
    copy(h2.slots[slot], h.slots[slot])
    h2.hash = h.hash
    return h2
}


func (h *HashTable[K, V]) String() string {
    return fmt.Sprintf("&{%v %v %v %v} load %v",
                       h.nentries,
		       h.nslots,
		       h.slots,
		       h.hash,
		       h.LoadFactor())
}


func (h *HashTable[K, V]) Get(k K) (V, bool) {
    slot := h.hash(k, h.nslots)
    for _, e := range h.slots[slot] {
        if e.k == k {
            return e.v, true
        }
    }
    return *new(V), false
}


func (h *HashTable[K, V]) Keys() []K {
    keys := make([]K, 0, h.nentries)
    for i, _ := range h.slots {
        for _, e := range h.slots[i] {
            keys = append(keys, e.k)
        }
    }
    return keys
}


func (h *HashTable[K, V]) LoadFactor() float32 {
    return float32(h.nentries) / float32(h.nslots)
}


func remove[T any](v []T, i int) []T {
    v[i] = v[len(v) - 1]
    v = v[:len(v) - 1]
    return v
}


func (h *HashTable[K, V]) remove(k K) *HashTable[K, V] {
    slot := h.hash(k, h.nslots)
    for i, e := range h.slots[slot] {
        if e.k == k {
            h2 := h.hashtablenewlayer(slot)
            h2.slots[slot] = remove(h2.slots[slot], i)
            h2.nentries--
            return h2
        }
    }
    return h
}


func (h *HashTable[K, V]) Remove(k K) *HashTable[K, V] {
    h = h.remove(k)
    if h.LoadFactor() < 0.3 {
        h = h.Resize(h.nslots / 2)
    }
    return h
}


func (h *HashTable[K, V]) Resize(nslots int) *HashTable[K, V] {
    h2 := HashTableNew[K, V](nslots, h.hash)
    for i, _ := range h.slots {
        for _, e := range h.slots[i] {
            h2 = h2.set(e.k, e.v)
        }
    }
    return h2
}


func (h *HashTable[K, V]) set(k K, v V) *HashTable[K, V] {
    slot := h.hash(k, h.nslots)
    h2 := h.hashtablenewlayer(slot)
    for i, e := range h2.slots[slot] {
        if e.k == k {
            h2.slots[slot][i].v = v
            return h2
        }
    }
    h2.slots[slot] = append(h2.slots[slot], HashTableEntry[K, V]{ k, v })
    h2.nentries++
    return h2
}


func (h *HashTable[K, V]) Set(k K, v V) *HashTable[K, V] {
    h = h.set(k, v)
    if h.LoadFactor() > 0.7 {
        h = h.Resize(h.nslots * 2)
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
    fmt.Println(hashtable)


    for _, e := range elements {
        hashtable = hashtable.Set(e, len(e))
        fmt.Println(hashtable)
    }

    v, ok := hashtable.Get("oxygen")
    fmt.Printf("%v %v\n", v, ok)


    for _, e := range hashtable.Keys() {
        hashtable = hashtable.Remove(e)
        fmt.Println(hashtable)
    }

    v, ok = hashtable.Get("oxygen")
    fmt.Printf("%v %v\n", v, ok)


    vs := []int{ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 }
    vs = remove(vs, 0)
    vs = remove(vs, 1)
    fmt.Printf("%v\n", vs)


    h := HashTableNew[string, int](1, hashsumchars)
    h2 := h.Set("tom", 1)
    h3 := h2.Set("dick", 2)
    h4 := h3.Set("harry", 3)
    h5 := h4.Set("harry", 4)
    h6 := h5.Remove("harry")
    h7 := h6.Remove("dick")
    h8 := h7.Remove("tom")
    fmt.Printf("%p %v\n", h, h)
    fmt.Printf("%p %v\n", h2, h2)
    fmt.Printf("%p %v\n", h3, h3)
    fmt.Printf("%p %v\n", h4, h4)
    fmt.Printf("%p %v\n", h5, h5)
    fmt.Printf("%p %v\n", h6, h6)
    fmt.Printf("%p %v\n", h7, h7)
    fmt.Printf("%p %v\n", h8, h8)


    os.Exit(0)
}
