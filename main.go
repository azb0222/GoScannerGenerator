package main

import (
	"fmt"
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
	regexRules := map[string]string{
		"test": "a(cat|cow)*",
	}
	fmt.Println(regexRules)
}

type NFAState struct {
	outputStates map[byte][]NFAState //input:nfaState
	isFinal      bool
}

type RegexASTNode struct {
	leftChild  *RegexASTNode
	rightChild *RegexASTNode
	rType      RegexType
	data       byte //TODO: type saftey? -> do generics
}

func convertRegexToAST(regex string) RegexASTNode {
	return RegexASTNode{}
}

//TODO: create helper print function?

func convertRegexToNfa(regexASTRootNode *RegexASTNode) (NFAState, NFAState) {
	switch regexASTRootNode.rType {
	case Concatenation:
		lNFAState, lNFALastState := convertRegexToNfa(regexASTRootNode.leftChild)
		rNFAState, rNFALastState := convertRegexToNfa(regexASTRootNode.rightChild)
		lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], rNFAState)
		return lNFAState, rNFALastState
	case Alternation:
		lNFAState, lNFALastState := convertRegexToNfa(regexASTRootNode.leftChild)
		rNFAState, rNFALastState := convertRegexToNfa(regexASTRootNode.rightChild)
		start := NFAState{
			outputStates: map[byte][]NFAState{
				epsilon: {lNFAState, rNFAState},
			},
		}
		end := NFAState{}
		lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], end)
		rNFALastState.outputStates[epsilon] = append(rNFALastState.outputStates[epsilon], end)

		return start, end
	case KleeneStar:
		lNFAState, lNFALastState := convertRegexToNfa(regexASTRootNode.leftChild)
		start := NFAState{
			outputStates: map[byte][]NFAState{
				epsilon: {lNFAState},
			},
		}
		end := NFAState{
			outputStates: map[byte][]NFAState{
				epsilon: {start},
			},
		}
		lNFALastState.outputStates[epsilon] = append(lNFALastState.outputStates[epsilon], end)
		start.outputStates[epsilon] = append(start.outputStates[epsilon], end)
		return start, end
	case Constant:
		end := NFAState{}
		start := NFAState{
			outputStates: map[byte][]NFAState{
				regexASTRootNode.data: {end},
			},
		}
		return start, end
	default:
		fmt.Errorf("BAD")
		return NFAState{}, NFAState{}
	}
}

type DFAState struct {
	outputState map[any]DFAState
	isFinal     bool
}

//
//func makeDFA() DFAState {
//	//TODO: do validation in initalizer function
//	// TODO: validate that a dfa is a proper dfa
//}
