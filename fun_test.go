package fun_test

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	fun "github.com/kirilldd2/go-no-fun"
)

func IntToFloat64(n int) float64 { return math.Sqrt(float64(n)) }

func createRandomSlice(n int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	slice := make([]int, n)
	for i := range slice {
		slice[i] = r.Intn(math.MaxInt)
	}

	return slice
}

func TestMap(t *testing.T) {
	t.Run("ints to float64 by sqrt", func(t *testing.T) {
		inp := createRandomSlice(1000)
		want := make([]float64, 1000)
		for i := range inp {
			want[i] = IntToFloat64(inp[i])
		}
		result := fun.Map(IntToFloat64, inp)
		if len(result) != len(inp) {
			t.Errorf("result len = %d, input's len = %d", len(result), len(inp))
		}
		for i := range result {
			if result[i] != want[i] {
				t.Errorf("on index = %d got %f, but wanted %f", i, result[i], want[i])
			}
		}
	})

	t.Run("empty input slice", func(t *testing.T) {
		var inp []int
		result := fun.Map(IntToFloat64, inp)
		if len(result) != 0 {
			t.Error("result for input empty slice is not empty slice")
		}
	})
}

func BenchmarkMap(b *testing.B) {
	inp := createRandomSlice(100000)

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fun.Map(IntToFloat64, inp)
		}
	})

	b.Run("For loop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := make([]float64, len(inp))
			for j := range inp {
				out[j] = IntToFloat64(inp[j])
			}
		}
	})

	// copy of map for concrete types
	MapIntToFloat64 := func(fn func(int) float64, data []int) []float64 {
		res := make([]float64, len(data))

		for i, item := range data {
			res[i] = fn(item)
		}

		return res
	}

	b.Run("Concrete Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MapIntToFloat64(IntToFloat64, inp)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("concat", func(t *testing.T) {
		inp := []string{"No", "Fun", " !"}
		want := "NoFun !"
		res := fun.Reduce(func(acc string, item string) string {
			return acc + item
		}, inp, "")
		if res != want {
			t.Errorf("got %v, wanted %v", res, want)
		}
	})

	t.Run("squared sum", func(t *testing.T) {
		inp := []string{"1", "2", "3", "4"}
		want := 30
		res := fun.Reduce(func(acc int, item string) int {
			num, _ := strconv.ParseInt(item, 10, 32)
			numInt := int(num)

			return numInt*numInt + acc
		}, inp, 0)
		if res != want {
			t.Errorf("got %v, wanted %v", res, want)
		}
	})
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		inp, out []int
		fn       func(int) bool
	}{
		{"some", []int{1, 2, 3, 4}, []int{3, 4}, func(item int) bool { return item >= 3 }},
		{"all", []int{1, 2, 3, 4}, []int{}, func(item int) bool { return item >= 6 }},
		{"none", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}, func(item int) bool { return item >= 0 }},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if result := fun.Filter(test.fn, test.inp); !fun.Equal(test.out, result) {
				t.Errorf("wanted %v, got %v", test.out, result)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name string
		inp  []int
		out  bool
	}{
		{"true", []int{1, 0, 0, 4}, true},
		{"false", []int{0, 0, 0, 0}, false},
		{"empty", []int{}, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if result := fun.Any(test.inp); test.out != result {
				t.Errorf("wanted %v, got %v", test.out, result)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name string
		inp  []int
		out  bool
	}{
		{"false", []int{1, 0, 0, 4}, false},
		{"true", []int{1, 2, 3, 4}, true},
		{"empty", []int{}, true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if result := fun.All(test.inp); test.out != result {
				t.Errorf("wanted %v, got %v", test.out, result)
			}
		})
	}
}
