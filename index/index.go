package index

import (
	"container/heap"
	"github.com/kanatohodets/carbonsearch/util"
	"sort"
)

type Metric uint64
type MetricSlice []Metric

func (a MetricSlice) Len() int           { return len(a) }
func (a MetricSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MetricSlice) Less(i, j int) bool { return a[i] < a[j] }

type Tag uint64
type TagSlice []Tag

func (a TagSlice) Len() int           { return len(a) }
func (a TagSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TagSlice) Less(i, j int) bool { return a[i] < a[j] }

type Index interface {
	Query([]Tag) ([]Metric, error)
	Name() string
}

func HashTag(tag string) Tag {
	return Tag(util.HashStr64(tag))
}

func HashTags(tags []string) []Tag {
	result := make([]Tag, len(tags))
	for i, tag := range tags {
		result[i] = HashTag(tag)
	}
	return result
}

func HashMetric(metric string) Metric {
	return Metric(util.HashStr64(metric))
}

func HashMetrics(metrics []string) []Metric {
	result := make([]Metric, len(metrics))
	for i, metric := range metrics {
		result[i] = HashMetric(metric)
	}
	return result
}

func SortMetrics(metrics []Metric) {
	sort.Sort(MetricSlice(metrics))
}

type MetricSetsHeap [][]Metric

func (h MetricSetsHeap) Len() int           { return len(h) }
func (h MetricSetsHeap) Less(i, j int) bool { return h[i][0] < h[j][0] }
func (h MetricSetsHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MetricSetsHeap) Push(x interface{}) {
	t := x.([]Metric)
	*h = append(*h, t)
}

func (h *MetricSetsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func UnionMetrics(metricSets [][]Metric) []Metric {
	h := &MetricSetsHeap{}
	heap.Init(h)
	for _, list := range metricSets {
		// empty sets have nothing to add to the union
		if len(list) > 0 {
			heap.Push(h, list)
		}
	}
	set := []Metric{}
	for h.Len() > 0 {
		cur := (*h)[0]
		metric := cur[0]
		if len(set) == 0 || set[len(set)-1] != metric {
			set = append(set, metric)
		}
		if len(cur) == 1 {
			heap.Pop(h)
		} else {
			(*h)[0] = cur[1:]
			heap.Fix(h, 0)
		}
	}
	return set
}

func IntersectMetrics(metricSets [][]Metric) []Metric {
	if len(metricSets) == 0 {
		return []Metric{}
	}

	h := &MetricSetsHeap{}
	heap.Init(h)
	for _, list := range metricSets {
		// any empty set --> empty intersection
		if len(list) == 0 {
			return []Metric{}
		}
		heap.Push(h, list)
	}
	set := []Metric{}
	for {
		cur := (*h)[0]
		smallestMetric := cur[0]
		present := 0
		fixups := make([]bool, h.Len())
		for i, candidate := range *h {
			if candidate[0] <= smallestMetric {
				fixups[i] = true
			}
			if candidate[0] == smallestMetric {
				present++
			}
		}

		// found something in every subset
		if present == len(metricSets) {
			if len(set) == 0 || set[len(set)-1] != smallestMetric {
				set = append(set, smallestMetric)
			}
		}

		for i, fix := range fixups {
			if fix {
				list := (*h)[i]
				if len(list) == 1 {
					return set
				}
				(*h)[i] = list[1:]
				heap.Fix(h, i)
			}
		}
	}
}

func SortTags(tags []Tag) {
	sort.Sort(TagSlice(tags))
}

type TagSetsHeap [][]Tag

func (h TagSetsHeap) Len() int           { return len(h) }
func (h TagSetsHeap) Less(i, j int) bool { return h[i][0] < h[j][0] }
func (h TagSetsHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *TagSetsHeap) Push(x interface{}) {
	t := x.([]Tag)
	*h = append(*h, t)
}

func (h *TagSetsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//TODO(btyler) can we keep the benefits of distinct tag/metric types without the copypasta?
func UnionTags(tagSets [][]Tag) []Tag {
	h := &TagSetsHeap{}
	heap.Init(h)
	for _, list := range tagSets {
		// empty sets have nothing to add to the union
		if len(list) > 0 {
			heap.Push(h, list)
		}
	}
	set := []Tag{}
	for h.Len() > 0 {
		cur := (*h)[0]
		tag := cur[0]
		if len(set) == 0 || set[len(set)-1] != tag {
			set = append(set, tag)
		}
		if len(cur) == 1 {
			heap.Pop(h)
		} else {
			(*h)[0] = cur[1:]
			heap.Fix(h, 0)
		}
	}
	return set
}

func IntersectTags(tagSets [][]Tag) []Tag {
	if len(tagSets) == 0 {
		return []Tag{}
	}

	h := &TagSetsHeap{}
	heap.Init(h)
	for _, list := range tagSets {
		// any empty set --> empty intersection
		if len(list) == 0 {
			return []Tag{}
		}
		heap.Push(h, list)
	}
	set := []Tag{}
	for {
		cur := (*h)[0]
		smallestTag := cur[0]
		present := 0
		fixups := make([]bool, h.Len())
		for i, candidate := range *h {
			if candidate[0] <= smallestTag {
				fixups[i] = true
			}
			if candidate[0] == smallestTag {
				present++
			}
		}

		// found something in every subset
		if present == len(tagSets) {
			if len(set) == 0 || set[len(set)-1] != smallestTag {
				set = append(set, smallestTag)
			}
		}

		for i, fix := range fixups {
			if fix {
				list := (*h)[i]
				if len(list) == 1 {
					return set
				}
				(*h)[i] = list[1:]
				heap.Fix(h, i)
			}
		}
	}
}
