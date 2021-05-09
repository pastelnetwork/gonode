package dupedetection

import (
	"errors"
	"math"
	"sort"

	onlinestats "github.com/dgryski/go-onlinestats"
	"github.com/montanaflynn/stats"
)

type rank struct {
	X     float64
	Y     float64
	Xrank float64
	Yrank float64
}

type Float64Data []float64

func (f Float64Data) Len() int { return len(f) }

func (f Float64Data) Get(i int) float64 { return f[i] }

func Spearman2(data1, data2 Float64Data) (float64, error) {
	//defer Measure(time.Now())
	r, _ := onlinestats.Spearman(data1, data2)
	return r, nil
}

func Spearman(data1, data2 Float64Data) (float64, error) {
	//defer Measure(time.Now())
	if data1.Len() < 3 || data2.Len() != data1.Len() {
		return math.NaN(), errors.New("invalid size of data")
	}

	ranks := []rank{}

	for index := 0; index < data1.Len(); index++ {
		x := data1.Get(index)
		y := data2.Get(index)
		ranks = append(ranks, rank{
			X: x,
			Y: y,
		})
	}

	sort.Slice(ranks, func(i int, j int) bool {
		return ranks[i].X < ranks[j].X
	})

	for position := 0; position < len(ranks); position++ {
		ranks[position].Xrank = float64(position) + 1

		duplicateValues := []int{position}
		for nested, p := range ranks {
			if ranks[position].X == p.X {
				if position != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}
		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}

		avg := float64((sum + len(duplicateValues))) / float64(len(duplicateValues))
		ranks[position].Xrank = avg

		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].Xrank = avg
		}

		position += len(duplicateValues) - 1
	}

	sort.Slice(ranks, func(i int, j int) bool {
		return ranks[i].Y < ranks[j].Y
	})

	for position := 0; position < len(ranks); position++ {
		ranks[position].Yrank = float64(position) + 1

		duplicateValues := []int{position}
		for nested, p := range ranks {
			if ranks[position].Y == p.Y {
				if position != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}
		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}
		// fmt.Println(sum + len(duplicateValues))
		avg := float64((sum + len(duplicateValues))) / float64(len(duplicateValues))
		ranks[position].Yrank = avg

		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].Yrank = avg
		}

		position += len(duplicateValues) - 1
	}

	xRanked := []float64{}
	yRanked := []float64{}

	for _, rank := range ranks {
		xRanked = append(xRanked, rank.Xrank)
		yRanked = append(yRanked, rank.Yrank)
	}

	return stats.Pearson(xRanked, yRanked)
}
