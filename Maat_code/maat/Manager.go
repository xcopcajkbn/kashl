package maat

import (
	"fmt"
	"os"
	"time"
)

var StorePath = "/root/maattest/" + time.Now().Format("2006-01-02 15:04:05") + "/"

type SManager struct {
	Fetcher                   *SFetcher
	Logger                    *SLogger
	Analyzer                  *SAnalyzer
	ChClientMetric            chan SClientMetric
	ChBlockMetric             chan SBlockMetric
	ChNoticeUnderlyingDBMeter chan int64 // fetch notice pebbledb thread
	ChUnderlyingDBMetric      chan SUnderlyingDBMetric
	ChMPTMetric               chan SMPTMetric
	ChOpcodeMetric            chan SOpcodeMetric
	ChOverheadMetric          chan SOverheadMetric
}

func NewSManager() *SManager {
	s := &SManager{
		ChClientMetric:            make(chan SClientMetric, 100),
		ChMPTMetric:               make(chan SMPTMetric, 10000),
		ChBlockMetric:             make(chan SBlockMetric, 17000),
		ChNoticeUnderlyingDBMeter: make(chan int64, 100),
		ChUnderlyingDBMetric:      make(chan SUnderlyingDBMetric, 100),
		ChOpcodeMetric:            make(chan SOpcodeMetric, 100005),
		ChOverheadMetric:          make(chan SOverheadMetric, 10005),
		Fetcher:                   newFetcher(),
		Logger:                    newLogger(),
		Analyzer:                  newAnalyzer(),
	}
	if err := os.MkdirAll(StorePath, os.ModePerm); err != nil {
		fmt.Println("can't init dir", err, StorePath)
		os.Exit(7)
	}
	s.Fetcher.SetFetchChan(s.ChClientMetric, s.ChBlockMetric, s.ChNoticeUnderlyingDBMeter, s.ChUnderlyingDBMetric, s.ChOpcodeMetric, s.ChOverheadMetric)
	s.Logger.SetLoggerChan(s.ChClientMetric, s.ChBlockMetric, s.ChUnderlyingDBMetric, s.ChMPTMetric, s.ChOpcodeMetric, s.ChOverheadMetric)
	s.Analyzer.SetAnalyzeChan(s.ChMPTMetric)
	go s.Logger.OpcodeInfoWorker()
	go s.Logger.OverheadInfoWorker()
	return s
}

func (manager *SManager) Exit() {
	close(manager.ChBlockMetric)
	close(manager.ChClientMetric)
	close(manager.ChNoticeUnderlyingDBMeter)
	close(manager.ChUnderlyingDBMetric)
	close(manager.ChMPTMetric)
	close(manager.ChOpcodeMetric)
	close(manager.ChOverheadMetric)
}

type SClientMetric struct {
	CurBlockNum                 int64
	BlockSyncTime               time.Duration
	InsertChainTime             time.Duration
	BlockReplayTime             time.Duration
	BlockProcTime               time.Duration
	BlockTxExecTime             time.Duration
	ContractTxExecTime          time.Duration
	NormalTxExecTime            time.Duration
	CreateTxExecTime            time.Duration
	BlockAuthTime               time.Duration
	BlockValideTime             time.Duration
	BlockWriteTime              time.Duration
	BlockCommitTime             time.Duration
	BlockstatecommitTime        time.Duration
	BlockStateAccountCommitTime time.Duration
	BlockStateStorageCommitTime time.Duration
	BlockFinaliseTime           time.Duration
	BlockAccountMPTUpdateTime   time.Duration //include storage trie hash time
	BlockStorageMPTUpdataTime   time.Duration
	DereferenceTime             time.Duration
	SASprocTime                 time.Duration
	UpdateInsertTime            time.Duration
	CodeReadTime                time.Duration
	AccountMPTReadTime          time.Duration
	StorageMPTReadTime          time.Duration
	AccountSASReadTime          time.Duration
	StorageSASReadTime          time.Duration
	AccountHashTime             time.Duration
	StorageHashTime             time.Duration
	HeaderWriteTime             time.Duration
	TrieWriteTime               time.Duration
	SASWriteTime                time.Duration
	CodeWriteTime               time.Duration
	ProcessInitTime             time.Duration
	BlockhashInsCounter         int64
	AccountHashRawCounter       int64
	AccountHashCounter          int64
	StorageHashRawCounter       int64
	StorageHashCounter          int64
	ContractModifyCounter       int64
}

type SBlockMetric struct {
	Gasused                     int64
	CurBlockNum                 int64
	BlockReplayTime             time.Duration
	BlockTxExecTime             time.Duration
	ContractTxExecTime          time.Duration
	NormalTxExecTime            time.Duration
	CreateTxExecTime            time.Duration
	BlockAuthTime               time.Duration
	BlockFinaliseTime           time.Duration
	BlockAccountMPTUpdateTime   time.Duration //include storage trie hash time
	BlockStorageMPTUpdataTime   time.Duration
	BlockValideTime             time.Duration
	BlockWriteTime              time.Duration
	BlockCommitTime             time.Duration
	BlockstatecommitTime        time.Duration
	BlockStateAccountCommitTime time.Duration
	BlockStateStorageCommitTime time.Duration
	DereferenceTime             time.Duration
	SASprocTime                 time.Duration
	UpdateInsertTime            time.Duration
	CodeReadTime                time.Duration
	AccountMPTReadTime          time.Duration
	StorageMPTReadTime          time.Duration
	AccountSASReadTime          time.Duration
	StorageSASReadTime          time.Duration
	AccountHashTime             time.Duration
	HeaderWriteTime             time.Duration
	TrieWriteTime               time.Duration
	SASWriteTime                time.Duration
	CodeWriteTime               time.Duration
	ReceiptProcTime             time.Duration
	ProcessInitTime             time.Duration
	AccountHashRawCounter       int64
	AccountHashCounter          int64
	StorageHashTime             time.Duration
	StorageHashRawCounter       int64
	StorageHashCounter          int64
	TxsCounter                  int64
	TxInfo                      []STxMetric
	BlockhashInsCounter         int64
	ContractModifyCounter       int64
	AccountDepthDistribution    map[int64]int64
	StorageDepthDistribution    map[int64]int64
}

type STxMetric struct {
	TxType                int64 //0 normalTx & 1 call contract &2 contract create
	Gasused               int64
	TxExecTime            time.Duration
	TxAuthTime            time.Duration
	AccountMPTReadTime    time.Duration
	StorageMPTReadTime    time.Duration
	AccountSASReadTime    time.Duration
	StorageSASReadTime    time.Duration
	AccountHashTime       time.Duration
	CodeReadTime          time.Duration
	AccountHashRawCounter int64
	AccountHashCounter    int64
	StorageHashTime       time.Duration
	StorageHashRawCounter int64
	StorageHashCounter    int64
	BlockhashInsCounter   int64
	CreateConstructorTime time.Duration
}

type SUnderlyingDBMetric struct {
	CurBlockNum    int64
	DiskReadMeter  float64 //iostats
	DiskWriteMeter float64
}

type SMPTMetric struct {
	CurBlockNum              int64
	ContractAddr             string //"world means account trie"
	BranchNodeNum            int64
	ExtendNodeNum            int64
	ValueNodeNum             int64
	TrieDepthCounters        map[int64]int64
	ExtendPrefixDisturbution map[string]int64
}

type SOpcodeMetric struct {
	Opcodeinfo    SOpcodecommon
	BytecodeExsit bool
	Transfer      bool
	AccountSAS    bool
	StorageSAS    bool
}

type SOpcodecommon struct {
	CurBlockNum   int64
	Tx            []byte
	Txtype        int64
	Opcode        byte
	Gasused       uint64
	ExecutionTime time.Duration
	Parameter     []byte
}

type SOverheadMetric struct {
	CurBlockNum int64
	ProcessTime time.Duration
	Alloc       uint64
}
