package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type RegexType int

const (
	epsilon               = 0x00
	Alternation RegexType = iota
	Concatenation
	KleeneStar
	Constant
)

func main() {
	// regexRules := map[string]string{
	// 	"test": "a(t|w)*",
	// }

    t := &RegexASTNode{
        rType: Constant,
        data: 't',
    }
    w := &RegexASTNode{
        rType: Constant,
        data: 'w',
    }
    a := &RegexASTNode{
        rType: Constant,
        data: 'a',
    }
    taltw := &RegexASTNode{
        rType: Alternation,
        leftChild: t,
        rightChild: w,
    }
    kstaltw := &RegexASTNode{
        rType: KleeneStar,
        leftChild: taltw,
    }
    root := &RegexASTNode{
        rType: Concatenation,
        leftChild: a,
        rightChild: kstaltw,
    }
	graph, _, m := convertRegexToNfa(root)
    // c := epsilonClosure(graph.outputStates[97][0])
    // print(c)
    print(m)
    println(makeMermaid(graph))
}

type NFAState struct {
    id           uint
	outputStates map[byte][]*NFAState //input:nfaState
	isFinal      bool
}

type RegexASTNode struct {
	leftChild  *RegexASTNode
	rightChild *RegexASTNode
	rType      RegexType
	data       byte //TODO: type saftey? -> do generics
}

func convertRegexToNfaRecursion(regexASTRootNode *RegexASTNode, idToState map[uint]*NFAState, id uint) (*NFAState, *NFAState, map[uint]*NFAState, uint) {
    switch regexASTRootNode.rType {
	case Concatenation:
		lNFAState, lNFALastState, idToState, id := convertRegexToNfaRecursion(regexASTRootNode.leftChild, idToState, id)
		rNFAState, rNFALastState, idToState, id := convertRegexToNfaRecursion(regexASTRootNode.rightChild, idToState, id)
		lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], rNFAState)
		return lNFAState, rNFALastState, idToState, id
	case Alternation:
		lNFAState, lNFALastState, idToState, id := convertRegexToNfaRecursion(regexASTRootNode.leftChild, idToState, id)
		rNFAState, rNFALastState, idToState, id := convertRegexToNfaRecursion(regexASTRootNode.rightChild, idToState, id)
		
        start := &NFAState{
			outputStates: map[byte][]*NFAState{
				epsilon: {lNFAState, rNFAState},
			},
            id: id,
		}
        idToState[id] = start
        id++
		
        end := &NFAState{
            id: id,
            outputStates: make(map[byte][]*NFAState),
        }
        idToState[id] = end
        id++

        lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], end)
		rNFALastState.outputStates[epsilon] = append(rNFALastState.outputStates[epsilon], end)

		return start, end, idToState, id
	case KleeneStar:
		lNFAState, lNFALastState, idToState, id := convertRegexToNfaRecursion(regexASTRootNode.leftChild, idToState, id)
		
        start := &NFAState{
			outputStates: map[byte][]*NFAState{
				epsilon: {lNFAState},
			},
            id: id,
		}
        idToState[id] = start
        id++

		end := &NFAState{
			outputStates: map[byte][]*NFAState{
				epsilon: {start},
			},
            id: id,
		}
        idToState[id] = end
        id++

		lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], end)
		start.outputStates[epsilon] = append(start.outputStates[epsilon], end)
		return start, end, idToState, id
	case Constant:
        start := &NFAState{
			outputStates: make(map[byte][]*NFAState),
            id: id,
		}
        idToState[id] = start
        id++

		end := &NFAState{
            id: id,
            outputStates: make(map[byte][]*NFAState),
        }
        idToState[id] = end
        id++
		
        start.outputStates[regexASTRootNode.data] = []*NFAState{ end }
        return start, end, idToState, id
	default:
		fmt.Errorf("BAD")
		return nil, nil, idToState, id
	}
}

func convertRegexToNfa(regexASTRootNode *RegexASTNode) (*NFAState, *NFAState, map[uint]*NFAState) {
    start, end, idMap, _ := convertRegexToNfaRecursion(regexASTRootNode, make(map[uint]*NFAState), 0)
    return start, end, idMap
}

func makeMermaidRecursion(rootState *NFAState, edges []string, closed []uint) ([]string, []uint){
    if slices.Contains(closed, rootState.id){
        return edges, closed
    }
    closed = append(closed, rootState.id)
    for key, nodes := range rootState.outputStates {
        for _, node := range nodes{
            edges, closed = makeMermaidRecursion(node, edges, closed)
            edges = append(edges, fmt.Sprintf("%d -- %d --> %d", rootState.id, key, node.id))
        }
    }
    return edges, closed
}

func makeMermaid(rootState *NFAState) string {
    edges, _ := makeMermaidRecursion(rootState, make([]string, 0), make([]uint, 0))
    return strings.Join(edges, "\n")
}

type DFAState struct {
	outputState map[byte]DFAState
	isFinal     bool
}

func epsilonClosureRecursion(initialState *NFAState, states []*NFAState, closed map[uint]struct{}) ([]*NFAState, map[uint]struct{}) {
    if _,ok := closed[initialState.id]; ok {
        return states, closed
    }
    closed[initialState.id] = struct{}{}
    states = append(states, initialState)
	for _, s := range initialState.outputStates[epsilon] {
        states, closed = epsilonClosureRecursion(s, states, closed)
	}
	return states, closed
}

func epsilonClosure(states ...*NFAState) ([]*NFAState, map[uint]struct{}) {
    allStates := make([]*NFAState, 0)
    closed := make(map[uint]struct{})
    for _, state := range states {
        allStates, closed = epsilonClosureRecursion(state, allStates, closed)
    }
    return allStates, closed
}

func idsToString(ids map[uint]struct{}) string{
    b, _ := json.Marshal(ids)
    return string(b)
}

type OpenListEntry struct {
    NFAstates []*NFAState
    state *DFAState
}

func convertNFAtoDFA(initialNFAState *NFAState, NFAStateMap map[uint]*NFAState) *DFAState {
	closed := make(map[string]*OpenListEntry)
	open := make([]OpenListEntry, 0)
    nodes, ids := epsilonClosure(initialNFAState)
    initialEntry := OpenListEntry{
        NFAstates: nodes,
        state: &DFAState{},
    }
	open = append(open, initialEntry)
    closed[idsToString(ids)] = &initialEntry

	for len(open) > 0 {
        d := open[0]
        open = open[1:]
        var i byte
        for i=1; i <= 255; i++ {
            T := make([]*NFAState, 0)
            for _, state := range d.NFAstates {
                T = append(T, state.outputStates[i]...)
            }
            Tclosure, idMap := epsilonClosure(T...)
            
        }
	}
    return initialEntry.state
}