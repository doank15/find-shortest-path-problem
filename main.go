package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Edge struct {
	To     int
	Weight int64 // Using int64 to handle large weights
}

type Graph struct {
	adj      [][]Edge
	n        int // Number of vertices
	directed bool
}

func NewGraph(n int, directed bool) *Graph {
	adj := make([][]Edge, n+1)
	for i := range adj {
		adj[i] = make([]Edge, 0)
	}
	return &Graph{
		adj:      adj,
		n:        n,
		directed: directed,
	}
}

func (g *Graph) AddEdge(from, to int, weight int64) {
	g.adj[from] = append(g.adj[from], Edge{To: to, Weight: weight})
	if !g.directed {
		g.adj[to] = append(g.adj[to], Edge{To: from, Weight: weight})
	}
}

func (g *Graph) HasNegativeEdge() bool {
	for i := 1; i <= g.n; i++ {
		for _, e := range g.adj[i] {
			if e.Weight < 0 {
				return true
			}
		}
	}
	return false
}

type Item struct {
	vertex   int
	distance int64
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// SPFA (Shortest Path Faster Algorithm)
func (g *Graph) SPFA(start, end int) ([]int, int64, bool) {
	// Initialize distances
	const INF int64 = 1e18
	dist := make([]int64, g.n+1)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	parent := make([]int, g.n+1)
	for i := range parent {
		parent[i] = -1
	}

	inQueue := make([]bool, g.n+1)
	queue := make([]int, 0, g.n)
	queue = append(queue, start)
	inQueue[start] = true

	// Track number of times a vertex is pushed to detect negative cycles
	pushCount := make([]int, g.n+1)
	pushCount[start] = 1

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false

		for _, edge := range g.adj[u] {
			v := edge.To
			if dist[u] != INF && dist[u]+edge.Weight < dist[v] {
				dist[v] = dist[u] + edge.Weight
				parent[v] = u

				// If v is not in queue, add it
				if !inQueue[v] {
					queue = append(queue, v)
					inQueue[v] = true
					pushCount[v]++

					// Check for negative cycle
					if pushCount[v] > g.n {
						fmt.Println("Negative cycle detected")
						return nil, -1, false
					}
				}
			}
		}
	}

	if dist[end] == INF {
		fmt.Println("No path exists")
		return nil, -1, false
	}

	path := []int{}
	for at := end; at != -1; at = parent[at] {
		path = append(path, at)
		if at == start {
			break
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, dist[end], true
}

// Dijkstra algorithm for graphs with non-negative weights
func (g *Graph) Dijkstra(start, end int) ([]int, int64, bool) {
	const INF int64 = 1e18
	dist := make([]int64, g.n+1)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	parent := make([]int, g.n+1)
	for i := range parent {
		parent[i] = -1
	}

	// Initialize priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{vertex: start, distance: 0})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item)

		if u.vertex == end {
			break
		}

		if u.distance > dist[u.vertex] {
			continue
		}

		for _, edge := range g.adj[u.vertex] {
			v := edge.To
			newDist := dist[u.vertex] + edge.Weight

			if newDist < dist[v] {
				dist[v] = newDist
				parent[v] = u.vertex
				heap.Push(&pq, &Item{vertex: v, distance: newDist})
			}
		}
	}

	if dist[end] == INF {
		fmt.Println("No path exists")
		return nil, -1, false
	}

	path := []int{}
	for at := end; at != -1; at = parent[at] {
		path = append(path, at)
		if at == start {
			break
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, dist[end], true
}

// Johnson algorithm for all pairs shortest paths
func (g *Graph) Johnson() ([][]int64, bool) {
	n := g.n
	tempGraph := NewGraph(n+1, true)

	for u := 1; u <= n; u++ {
		for _, edge := range g.adj[u] {
			tempGraph.AddEdge(u, edge.To, edge.Weight)
		}
	}

	newSource := n + 1
	for i := 1; i <= n; i++ {
		tempGraph.AddEdge(newSource, i, 0)
	}

	h := make([]int64, n+2)
	for i := range h {
		h[i] = math.MaxInt64
	}
	h[newSource] = 0

	for i := 1; i <= n+1; i++ {
		for u := 1; u <= n+1; u++ {
			for _, edge := range tempGraph.adj[u] {
				v := edge.To
				w := edge.Weight
				if h[u] != math.MaxInt64 && h[u]+w < h[v] {
					h[v] = h[u] + w
					if i == n+1 {
						fmt.Println("Negative cycle detected")
						return nil, false
					}
				}
			}
		}
	}

	reweightedGraph := NewGraph(n, true)
	for u := 1; u <= n; u++ {
		for _, edge := range g.adj[u] {
			v := edge.To
			newWeight := edge.Weight + h[u] - h[v]
			reweightedGraph.AddEdge(u, v, newWeight)
		}
	}

	// Run Dijkstra from each vertex
	dist := make([][]int64, n+1)
	for i := range dist {
		dist[i] = make([]int64, n+1)
		for j := range dist[i] {
			dist[i][j] = math.MaxInt64
		}
	}

	for u := 1; u <= n; u++ {
		// Single-source Dijkstra
		dijkstraDist := make([]int64, n+1)
		for i := range dijkstraDist {
			dijkstraDist[i] = math.MaxInt64
		}
		dijkstraDist[u] = 0

		pq := make(PriorityQueue, 0)
		heap.Init(&pq)
		heap.Push(&pq, &Item{vertex: u, distance: 0})

		for pq.Len() > 0 {
			item := heap.Pop(&pq).(*Item)
			v := item.vertex

			if item.distance > dijkstraDist[v] {
				continue
			}

			for _, edge := range reweightedGraph.adj[v] {
				to := edge.To
				weight := edge.Weight

				if dijkstraDist[v]+weight < dijkstraDist[to] {
					dijkstraDist[to] = dijkstraDist[v] + weight
					heap.Push(&pq, &Item{vertex: to, distance: dijkstraDist[to]})
				}
			}
		}

		for v := 1; v <= n; v++ {
			if dijkstraDist[v] != math.MaxInt64 {
				dist[u][v] = dijkstraDist[v] - h[u] + h[v]
			} else {
				dist[u][v] = math.MaxInt64
			}
		}
	}

	return dist, true
}

// FindShortestPath chooses the appropriate algorithm based on graph properties
func (g *Graph) FindShortestPath(start, end int) ([]int, int64, bool) {
	if start <= 0 || start > g.n || end <= 0 || end > g.n {
		fmt.Println("Invalid vertices")
		return nil, -1, false
	}

	if g.HasNegativeEdge() {
		if !g.directed {
			fmt.Println("Using SPFA to find shortest path or detect negative cycles")
		} else {
			fmt.Println("Graph has negative edges. Using SPFA algorithm")
		}
		return g.SPFA(start, end)
	} else {
		fmt.Println("Graph has only non-negative edges. Using Dijkstra's algorithm")
		return g.Dijkstra(start, end)
	}
}

func main() {
	n := 5

	// Create a directed graph
	g := NewGraph(n, true)

	// Add edges
	g.AddEdge(1, 2, 10)
	g.AddEdge(1, 3, 3)
	g.AddEdge(2, 3, 1)
	g.AddEdge(2, 4, 2)
	g.AddEdge(3, 2, 4)
	g.AddEdge(3, 4, 8)
	g.AddEdge(3, 5, 2)
	g.AddEdge(4, 5, 7)

	start, end := 1, 5
	path, distance, ok := g.FindShortestPath(start, end)

	if ok {
		fmt.Printf("Shortest path from %d to %d: %v\n", start, end, path)
		fmt.Printf("Total distance: %d\n", distance)
	}

	fmt.Println("\nExample with negative weights:")
	g2 := NewGraph(4, true)
	g2.AddEdge(1, 2, 3)
	g2.AddEdge(2, 3, -8)
	g2.AddEdge(1, 3, 5)
	g2.AddEdge(3, 4, 2)

	path, distance, ok = g2.FindShortestPath(1, 4)
	if ok {
		fmt.Printf("Shortest path from 1 to 4: %v\n", path)
		fmt.Printf("Total distance: %d\n", distance)
	}

	fmt.Println("\nExample with negative cycle:")
	g3 := NewGraph(3, true)
	g3.AddEdge(1, 2, 2)
	g3.AddEdge(2, 3, 3)
	g3.AddEdge(3, 1, -6)

	path, distance, ok = g3.FindShortestPath(1, 3)
	if !ok {
		fmt.Println("As expected, negative cycle detected")
	}
}
