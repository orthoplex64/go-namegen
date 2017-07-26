// Package distr is for randomly choosing from a set according to variably
// weighted probability distributions.
package distr

import (
	"fmt"
	"math"
	"math/rand"
)

// StrDistr is a distribution of strings.
type StrDistr struct {
	p   map[string]float64
	sum float64
}

// NewStrDistr returns a new StrDistr.
func NewStrDistr() *StrDistr {
	d := new(StrDistr)
	d.p = make(map[string]float64)
	return d
}

// Add adds s to d with weight f >= 0, adding to any existing weight.
// In this sense, "Add" means both set inclusion and numeric addition.
func (d *StrDistr) Add(s string, f float64) *StrDistr {
	f = math.Max(f, 0)
	d.p[s] += f
	d.sum += f
	return d
}

// Strings returns the strings in d, in no particular order.
func (d *StrDistr) Strings() []string {
	res := make([]string, len(d.p))
	for s := range d.p {
		res = append(res, s)
	}
	return res
}

// Weight returns the weight of s in d.
func (d *StrDistr) Weight(s string) float64 {
	return d.p[s]
}

// Remove removes s from d, returning the weight of s.
func (d *StrDistr) Remove(s string) float64 {
	p := d.p[s]
	delete(d.p, s)
	d.sum -= p
	return p
}

// Sum returns the sum of the weights in d.
func (d *StrDistr) Sum() float64 {
	return d.sum
}

// Pick randomly chooses a string according to the current distribution.
func (d *StrDistr) Pick() string {
	for {
		r := rand.Float64() * d.sum
		for str, p := range d.p {
			r -= p
			if r < 0 {
				return str
			}
		}
		fmt.Println("distr: end of pick loop reached (this should happen extremely rarely if at all)")
	}
}
