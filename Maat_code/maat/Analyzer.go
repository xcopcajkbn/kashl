package maat

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
)

type SAnalyzer struct {
	MPTMetric     *SMPTMetric
	ChMPTMetric   chan SMPTMetric
	TrieStorePath string
}

func newAnalyzer() *SAnalyzer {
	a := &SAnalyzer{
		MPTMetric:     new(SMPTMetric),
		TrieStorePath: StorePath + "triedata/",
	}
	a.MPTMetric.TrieDepthCounters = make(map[int64]int64)
	a.MPTMetric.ExtendPrefixDisturbution = make(map[string]int64)
	return a
}

func (analyzer *SAnalyzer) InitAnalyzer(blockheight int64, ContractAddr string) {
	analyzer.MPTMetric = new(SMPTMetric)
	analyzer.TrieStorePath = StorePath + "triedata/trieinfo-" + strconv.FormatInt(blockheight, 10)
	analyzer.MPTMetric.CurBlockNum = blockheight
	analyzer.MPTMetric.ContractAddr = ContractAddr
	analyzer.MPTMetric.TrieDepthCounters = make(map[int64]int64)
	analyzer.MPTMetric.ExtendPrefixDisturbution = make(map[string]int64)
	if err := os.MkdirAll(analyzer.TrieStorePath, os.ModePerm); err != nil {
		fmt.Println("can't init trie dir", err, StorePath)
		os.Exit(7)
	}
}

func (analyzer *SAnalyzer) SetAnalyzeChan(chMPTMetric chan SMPTMetric) {
	analyzer.ChMPTMetric = chMPTMetric
}

func (analyzer *SAnalyzer) SendMptInfo() {
	MPTInfo := analyzer.MPTMetric
	var updateMPTINFO = *MPTInfo //deep copy
	analyzer.ChMPTMetric <- updateMPTINFO
}

// traverse
func (analyzer *SAnalyzer) AnalyzerCounter(Type int64, Counter0 int64, Counter1 int64) {
	MPTInfo := analyzer.MPTMetric
	if Type == 1 {
		MPTInfo.BranchNodeNum += Counter0
	} else if Type == 2 {
		MPTInfo.ExtendNodeNum += Counter0
	} else if Type == 3 {
		MPTInfo.ValueNodeNum += Counter0
	} else if Type == 4 {
		MPTInfo.TrieDepthCounters[Counter0] += 1
	} else if Type == 5 {
		s0 := strconv.FormatInt(Counter0, 10)
		s1 := strconv.FormatInt(Counter1, 10)
		MPTInfo.ExtendPrefixDisturbution[s0+"#"+s1] += 1
	} else {
		log.Error("hzyserr", "critical bug", Type, Counter0)
		os.Exit(0)
	}
}
