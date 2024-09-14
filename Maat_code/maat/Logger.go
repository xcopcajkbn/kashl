package maat

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
)

type SLogger struct {
	ChClientMetric       chan SClientMetric
	ChBlockMetric        chan SBlockMetric
	ChOpcodeMetric       chan SOpcodeMetric
	ChOverheadMetric     chan SOverheadMetric
	ClientRecords        []SClientMetric
	BlockRecords         []SBlockMetric
	UnderlyingDBReords   []SUnderlyingDBMetric
	OpcodeRecords        []SOpcodeMetric
	OverheadRecords      []SOverheadMetric
	ChUnderlyingDBMetric chan SUnderlyingDBMetric
	ChMPTMetric          chan SMPTMetric
}

func newLogger() *SLogger {
	f := &SLogger{
		ClientRecords:      make([]SClientMetric, 0, 11),
		BlockRecords:       make([]SBlockMetric, 0, 1e4+10),
		OpcodeRecords:      make([]SOpcodeMetric, 0, 100005),
		OverheadRecords:    make([]SOverheadMetric, 0, 10005),
		UnderlyingDBReords: make([]SUnderlyingDBMetric, 0, 11),
	}
	return f
}

func (logger *SLogger) SetLoggerChan(chClientMetric chan SClientMetric, chBlockMetric chan SBlockMetric, chChUnderlyingDBMetric chan SUnderlyingDBMetric, chMPTMetric chan SMPTMetric, chOpcodeMetric chan SOpcodeMetric, chOverheadMetric chan SOverheadMetric) {
	logger.ChClientMetric = chClientMetric
	logger.ChBlockMetric = chBlockMetric
	logger.ChUnderlyingDBMetric = chChUnderlyingDBMetric
	logger.ChMPTMetric = chMPTMetric
	logger.ChOpcodeMetric = chOpcodeMetric
	logger.ChOverheadMetric = chOverheadMetric
}

func (logger *SLogger) ClientInfoWorker() {
	counter := 0
	for {
		Clientinfo, ok := <-logger.ChClientMetric
		logger.ClientRecords = append(logger.ClientRecords, Clientinfo)
		if !ok {
			strBlocknum := strconv.FormatInt(Clientinfo.CurBlockNum, 10)
			logger.LoggerOutput(0, StorePath+"clientinfo-0-"+strBlocknum)
			//prior = strBlocknum
			break
		}
		counter++
		if counter == 10 {
			strBlocknum := strconv.FormatInt(Clientinfo.CurBlockNum, 10)
			logger.LoggerOutput(0, StorePath+"clientinfo-0-"+strBlocknum)
			//prior = strBlocknum
			logger.ClientRecords = make([]SClientMetric, 0, 11)
			counter = 0
		}
	}
}

func (logger *SLogger) BlockInfoWorker() {
	var prior string = "0"
	counter := 0
	for {
		Blockinfo, ok := <-logger.ChBlockMetric
		if !ok {
			strBlocknum := strconv.FormatInt(Blockinfo.CurBlockNum, 10)
			logger.LoggerOutput(1, StorePath+"blockinfo-"+prior+"-"+strBlocknum)
			prior = strBlocknum
			break
		}
		logger.BlockRecords = append(logger.BlockRecords, Blockinfo)
		counter++
		if counter == 1e4 {
			strBlocknum := strconv.FormatInt(Blockinfo.CurBlockNum, 10)
			logger.LoggerOutput(1, StorePath+"blockinfo-"+prior+"-"+strBlocknum)
			prior = strBlocknum
			logger.BlockRecords = make([]SBlockMetric, 0, 1e4+10)
			counter = 0
		}
	}
}

// measure for the bottom
func (logger *SLogger) UnderlyingDBnfoWorker() {
	counter := 0
	for {
		UnderlyingDBinfo, ok := <-logger.ChUnderlyingDBMetric
		logger.UnderlyingDBReords = append(logger.UnderlyingDBReords, UnderlyingDBinfo)
		if !ok {
			strBlocknum := strconv.FormatInt(UnderlyingDBinfo.CurBlockNum, 10)
			logger.LoggerOutput(2, StorePath+"UnderlyingDB-0-"+strBlocknum)
			//prior = strBlocknum
			break
		}
		counter++
		if counter == 10 {
			strBlocknum := strconv.FormatInt(UnderlyingDBinfo.CurBlockNum, 10)
			logger.LoggerOutput(2, StorePath+"UnderlyingDB-0-"+strBlocknum)
			//prior = strBlocknum
			logger.UnderlyingDBReords = make([]SUnderlyingDBMetric, 0, 11)
			counter = 0
		}
	}
}

func (logger *SLogger) OpcodeInfoWorker() {
	prior := "0"
	for {
		OpcodeInfo, ok := <-logger.ChOpcodeMetric
		strBlocknum := strconv.FormatInt(OpcodeInfo.Opcodeinfo.CurBlockNum, 10)
		if !ok {
			strBlocknum = strconv.FormatInt(OpcodeInfo.Opcodeinfo.CurBlockNum, 10)
			logger.LoggerOutput(3, StorePath+"opcodeinfo-0-"+strBlocknum)
			prior = strBlocknum
			break
		}
		if prior != strBlocknum {
			logger.LoggerOutput(3, StorePath+"opcodeinfo-0-"+prior)
			prior = strBlocknum
			logger.OpcodeRecords = make([]SOpcodeMetric, 0, 100005)
		}
		logger.OpcodeRecords = append(logger.OpcodeRecords, OpcodeInfo)
	}
}

func (logger *SLogger) OverheadInfoWorker() {
	for {
		OverheadInfo, ok := <-logger.ChOverheadMetric
		logger.OverheadRecords = append(logger.OverheadRecords, OverheadInfo)
		if !ok {
			strBlocknum := strconv.FormatInt(OverheadInfo.CurBlockNum, 10)
			logger.LoggerOutput(6, StorePath+"overheadinfo-0-"+strBlocknum)
			break
		}
		if len(logger.OverheadRecords) == 10000 {
			strBlocknum := strconv.FormatInt(OverheadInfo.CurBlockNum, 10)
			logger.LoggerOutput(6, StorePath+"overheadinfo-0-"+strBlocknum)
			logger.OverheadRecords = make([]SOverheadMetric, 0, 10005)
		}
	}

}

func (logger *SLogger) LoggerOutput(Type int, filepath string) {
	//prior := 0
	dataFile, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//os.Exit(0)
	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)

	if Type == 0 {
		dataEncoder.Encode(logger.ClientRecords)
	} else if Type == 1 {
		dataEncoder.Encode(logger.BlockRecords)
	} else if Type == 2 {
		dataEncoder.Encode(logger.UnderlyingDBReords)
	} else if Type == 3 {
		dataEncoder.Encode(logger.OpcodeRecords)
	} else if Type == 6 {
		dataEncoder.Encode(logger.OverheadRecords)
	}
	dataFile.Close()
}
