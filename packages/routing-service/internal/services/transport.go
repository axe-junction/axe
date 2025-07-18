// transport_graph.go
package services

import (
	"container/heap"
	"math"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

type Node struct {
	ID        uuid.UUID
	StationID uuid.UUID
	Latitude  float64
	Longitude float64
}

type Edge struct {
	From     uuid.UUID
	To       uuid.UUID
	LineID   uuid.UUID
	Duration float64
	Distance float64
	Mode     string // "walking", "bus", "tram", etc.
}

type TransportGraph struct {
	Nodes    map[uuid.UUID]*Node
	Stations map[uuid.UUID]models.Station
	Edges    map[uuid.UUID][]Edge
}

func NewTransportGraph() *TransportGraph {
	return &TransportGraph{
		Nodes:    make(map[uuid.UUID]*Node),
		Stations: make(map[uuid.UUID]models.Station),
		Edges:    make(map[uuid.UUID][]Edge),
	}
}

func (g *TransportGraph) AddStation(station models.Station) {
	nodeID := uuid.New()
	g.Nodes[nodeID] = &Node{
		ID:        nodeID,
		StationID: station.ID,
		Latitude:  station.Latitude,
		Longitude: station.Longitude,
	}
	g.Stations[station.ID] = station
}

func (g *TransportGraph) AddConnection(from, to, lineID uuid.UUID, duration, distance float64, mode string) {
	g.Edges[from] = append(g.Edges[from], Edge{
		From:     from,
		To:       to,
		LineID:   lineID,
		Duration: duration,
		Distance: distance,
		Mode:     mode,
	})
}

func (g *TransportGraph) AddTransfer(from, to uuid.UUID, duration, distance float64) {
	g.Edges[from] = append(g.Edges[from], Edge{
		From:     from,
		To:       to,
		LineID:   uuid.Nil,
		Duration: duration,
		Distance: distance,
		Mode:     "walking",
	})
}
func (g *TransportGraph) EdgeCount() int {
	count := 0
	for _, edges := range g.Edges {
		count += len(edges)
	}
	return count
}
func (g *TransportGraph) ShortestPath(start, end uuid.UUID) []Edge {
	// Dijkstra's algorithm implementation
	dist := make(map[uuid.UUID]float64)
	prev := make(map[uuid.UUID]Edge)
	pq := make(PriorityQueue, 0)

	for node := range g.Nodes {
		dist[node] = math.Inf(1)
	}
	dist[start] = 0

	heap.Push(&pq, &Item{
		node:     start,
		priority: 0,
	})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		u := item.node

		if u == end {
			break
		}

		for _, edge := range g.Edges[u] {
			v := edge.To
			alt := dist[u] + edge.Duration
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = edge
				heap.Push(&pq, &Item{
					node:     v,
					priority: alt,
				})
			}
		}
	}

	// Reconstruct path
	path := []Edge{}
	u := end
	for u != start {
		edge, exists := prev[u]
		if !exists {
			return nil
		}
		path = append([]Edge{edge}, path...)
		u = edge.From
	}

	return path
}

// Priority queue implementation
type Item struct {
	node     uuid.UUID
	priority float64
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
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
