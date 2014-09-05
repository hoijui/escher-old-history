// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package understand

import (
	"fmt"

	"github.com/gocircuit/escher/see"
	. "github.com/gocircuit/escher/image"
)

func Understand(s *see.Circuit) *Circuit {
	x := &Circuit{peer: Make()}
	x.genus = []*see.Circuit{s}
	x.name = s.Name

	// Add “this” circuit as the empty-string peer
	sup := &Peer{
		name: "",
		index: 0,
		valve: Make(),
		design: nil,
	}
	x.peer[""] = sup
	x.index = 1

	// Add peers from circuit definition, valves are not added on this pass
	for _, p := range s.Peer {
		x.addPeer(p.Name.(string), x.index, p.Design) // all peer names in a syntactic (can be seen) circuit are strings
		x.index++
	}
	var nsugar int // Counter for generating names of desugared peer definitions
	for _, m := range s.Match {
		var end [2]*Valve // reciprocals
		for i, j := range m.Join {
			switch t := j.(type) {
			case *see.DesignJoin: // unfold sugar
				nsugar++
				p := fmt.Sprintf("sugar#%d", nsugar)
				x.addPeer(p, x.index, t.Design)
				x.index++
				end[i] = x.reserveValve(p, see.DefaultValve, x.index)
				x.index++
			case *see.PeerJoin:
				end[i] = x.reserveValve(t.Peer, t.Valve, x.index)
				x.index++
			case *see.ValveJoin:
				end[i] = x.reserveValve("", t.Valve, x.index)
				x.index++
			case nil: // match other argument to empty-string valve of this circuit
				end[i] = x.reserveValve("", "", x.index)
				x.index++
			default:
				panic(fmt.Sprintf("unknown or missing matching endpoint: %T·%v", j, j))
			}
		}
		// Link two ends
		if end[0].Matching != nil || end[1].Matching != nil {
			panic("reuse of valve")
		}
		end[0].Matching, end[1].Matching = end[1], end[0]
	}

	// Verify no dangling/unmatched valves remain
	for _, pn := range x.PeerNames() {
		p := x.PeerByName(pn)
		for _, vn := range p.ValveNames() {
			if p.ValveByName(vn).Matching == nil {
				panic("unmatched valve")
			}
		}
	}

	return x
}

func (x *Circuit) Merge(y *Circuit) {
	if len(y.genus) != 1 {
		panic("?")
	}
	x.genus = append(x.genus, y.genus[0])
	var m int // track max index
	for _, pn := range y.PeerNames() {
		if x.PeerByName(pn) != nil {
			panicf("collision adding peer %s in circuit %s", pn, x.Name())
		}
		q := y.PeerByName(pn).Copy()
		q.index += x.index
		m = max(m, q.index)
		x.attachPeer(q)
	}
	x.index += y.index
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (x *Circuit) attachPeer(p *Peer) {
	if _, ok := x.peer[p.name]; ok {
		panic("peer name already present")
	}
	if p.design == nil {
		panic("peer is missing design")
	}
	x.peer[p.name] = p
}

func (x *Circuit) addPeer(name interface{}, index int, design interface{}) {
	if _, ok := x.peer[name]; ok {
		panic("peer name already present")
	}
	if design == nil {
		panic("peer is missing design")
	}
	p := &Peer{
		name: name,
		index: index,
		valve: Make(),
		design: design,
	}
	x.peer[name] = p
}

// reserveValve returns the addressed valve, creating it if necessary.
// Creating is prohibited solely for the empty-string peer, corresponding to this circuit.
func (x *Circuit) reserveValve(peer, valve string, index int) *Valve {
	p := x.PeerByName(peer)
	if p == nil {
		panic(fmt.Sprintf("peer %v is missing", peer))
	}
	v := p.ValveByName(valve)
	if v != nil {
		return v
	}
	v = &Valve{
		Of: p,
		Name: valve,
		Index: index,
		Matching: nil,
	}
	p.valve[valve] = v
	return v
}