package themis

import (
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type SFetcher struct {
	CurrentBlockNum           int64
	CurTxIndex                int
	ClientMetric              *SClientMetric
	BlockMetric               *SBlockMetric
	OpcodeMetric              *SOpcodeMetric
	ChClientMetric            chan SClientMetric
	ChBlockMetric             chan SBlockMetric
	ChNoticeUnderlyingDBMeter chan int64
	ChUnderlyingDBMetric      chan SUnderlyingDBMetric
	ChOpcodeMetric            chan SOpcodeMetric
	ChOverheadMetric          chan SOverheadMetric
	GenesisTime               time.Time
}

func newFetcher() *SFetcher {
	f := &SFetcher{
		//SFetcher: newFetcher(false, 1e5),
		ClientMetric: new(SClientMetric),
		BlockMetric:  new(SBlockMetric),
	}
	return f
}

func (fetcher *SFetcher) SetFetchChan(chClientMetric chan SClientMetric, chBlockMetric chan SBlockMetric, chNoticeUnderlyingDBMeter chan int64, chUnderlyingDBMetric chan SUnderlyingDBMetric, chOpcodeMetric chan SOpcodeMetric, chOverheadMetric chan SOverheadMetric) {
	fetcher.ChClientMetric = chClientMetric
	fetcher.ChBlockMetric = chBlockMetric
	fetcher.ChNoticeUnderlyingDBMeter = chNoticeUnderlyingDBMeter
	fetcher.ChUnderlyingDBMetric = chUnderlyingDBMetric
	fetcher.ChOpcodeMetric = chOpcodeMetric
	fetcher.ChOverheadMetric = chOverheadMetric
}

func (fetcher *SFetcher) SetCurTxIndex(i int) {
	fetcher.CurTxIndex = i
}

func (fetcher *SFetcher) GetCurTxindex() int {
	return fetcher.CurTxIndex
}

func (fetcher *SFetcher) SetBlockData(height int64, gasused int64, txsnum int64) {
	fetcher.CurrentBlockNum = height
	if height == 1 {
		fetcher.GenesisTime = time.Now()
	}
	fetcher.BlockMetric = new(SBlockMetric)
	fetcher.BlockMetric.Gasused = gasused
	fetcher.BlockMetric.TxsCounter = txsnum
	fetcher.BlockMetric.TxInfo = make([]STxMetric, txsnum+1)
	fetcher.BlockMetric.AccountDepthDistribution = make(map[int64]int64)
	fetcher.BlockMetric.StorageDepthDistribution = make(map[int64]int64)
	fetcher.CurTxIndex = 0
}

func (fetcher *SFetcher) DispatchClientInfo() {
	//fetcher.ClientMetric.CurBlockNum = fetcher.CurrentBlockNum
	var ClientInfo SClientMetric = *fetcher.ClientMetric //deepcopy
	ClientInfo.CurBlockNum = fetcher.CurrentBlockNum
	ClientInfo.BlockSyncTime = time.Since(fetcher.GenesisTime)
	fetcher.ChClientMetric <- ClientInfo
}

func (fetcher *SFetcher) DispatchBlockInfo() {
	var BlockInfo SBlockMetric = *fetcher.BlockMetric //deepcopy
	BlockInfo.CurBlockNum = fetcher.CurrentBlockNum
	fetcher.ChBlockMetric <- BlockInfo
}

func (fetcher *SFetcher) FetchClientMetric(Type int, metric time.Duration) {
	clientmetric := fetcher.ClientMetric
	blockmetric := fetcher.BlockMetric
	txindex := fetcher.GetCurTxindex()
	//fmt.Println(txindex)
	txinfo := &fetcher.BlockMetric.TxInfo[txindex]
	if Type == 1 {
		clientmetric.BlockReplayTime += metric
		blockmetric.BlockReplayTime = metric
	} else if Type == 2 {
		clientmetric.BlockTxExecTime += metric
		blockmetric.BlockTxExecTime += metric
		txinfo.TxExecTime += metric
	} else if Type == 3 {
		clientmetric.ContractTxExecTime += metric
		blockmetric.ContractTxExecTime += metric
	} else if Type == 4 {
		clientmetric.BlockAuthTime += metric
		blockmetric.BlockAuthTime += metric
		txinfo.TxAuthTime += metric
	} else if Type == 5 {
		clientmetric.BlockCommitTime += metric
		blockmetric.BlockCommitTime += metric
	} else if Type == 6 {
		clientmetric.BlockWriteTime += metric
		blockmetric.BlockWriteTime += metric
	} else if Type == 7 {
		clientmetric.AccountMPTReadTime += metric
		blockmetric.AccountMPTReadTime += metric
		txinfo.AccountMPTReadTime += metric
	} else if Type == 8 {
		clientmetric.StorageMPTReadTime += metric
		blockmetric.StorageMPTReadTime += metric
		txinfo.StorageMPTReadTime += metric
	} else if Type == 9 {
		clientmetric.AccountHashTime += metric
		blockmetric.AccountHashTime += metric
		txinfo.AccountHashTime += metric
	} else if Type == 10 {
		clientmetric.StorageHashTime += metric
		blockmetric.StorageHashTime += metric
		txinfo.StorageHashTime += metric
	} else if Type == 12 {
		clientmetric.BlockValideTime += metric
		blockmetric.BlockValideTime += metric
	} else if Type == 13 {
		clientmetric.BlockProcTime += metric
	} else if Type == 14 {
		clientmetric.BlockAccountMPTUpdateTime += metric //include storage trie hash time
		blockmetric.BlockAccountMPTUpdateTime += metric
	} else if Type == 15 {
		clientmetric.BlockStorageMPTUpdataTime += metric
		blockmetric.BlockStorageMPTUpdataTime += metric
	} else if Type == 16 {
		clientmetric.BlockStateAccountCommitTime += metric
		blockmetric.BlockStateAccountCommitTime += metric
	} else if Type == 17 {
		clientmetric.BlockStateStorageCommitTime += metric
		blockmetric.BlockStateStorageCommitTime += metric
	} else if Type == 18 {
		clientmetric.BlockstatecommitTime += metric
		blockmetric.BlockstatecommitTime += metric
	} else if Type == 19 {
		clientmetric.BlockFinaliseTime += metric
		blockmetric.BlockFinaliseTime += metric
	} else if Type == 20 {
		clientmetric.NormalTxExecTime += metric
		blockmetric.NormalTxExecTime += metric
	} else if Type == 21 {
		clientmetric.CreateTxExecTime += metric
		blockmetric.CreateTxExecTime += metric
	} else if Type == 22 {
		clientmetric.InsertChainTime += metric
	} else if Type == 24 {
		clientmetric.DereferenceTime += metric
		blockmetric.DereferenceTime += metric
	} else if Type == 25 {
		clientmetric.SASprocTime += metric
		blockmetric.SASprocTime += metric
	} else if Type == 26 {
		clientmetric.UpdateInsertTime += metric
		blockmetric.UpdateInsertTime += metric
	} else if Type == 27 {
		clientmetric.CodeReadTime += metric
		blockmetric.CodeReadTime += metric
		txinfo.CodeReadTime += metric
	} else if Type == 28 {
		clientmetric.AccountSASReadTime += metric
		blockmetric.AccountSASReadTime += metric
		txinfo.AccountSASReadTime += metric
	} else if Type == 29 {
		clientmetric.StorageSASReadTime += metric
		blockmetric.StorageSASReadTime += metric
		txinfo.StorageSASReadTime += metric
	} else if Type == 30 {
		clientmetric.HeaderWriteTime += metric
		blockmetric.HeaderWriteTime += metric
	} else if Type == 31 {
		clientmetric.TrieWriteTime += metric
		blockmetric.TrieWriteTime += metric
	} else if Type == 32 {
		clientmetric.SASWriteTime += metric
		blockmetric.SASWriteTime += metric
	} else if Type == 33 {
		clientmetric.CodeWriteTime += metric
		blockmetric.CodeWriteTime += metric
	} else if Type == 34 {
		clientmetric.ProcessInitTime += metric
		blockmetric.ProcessInitTime += metric
	} else {
		log.Error("hzyserr", "critical bug", Type, metric)
		os.Exit(0)
	}
}

func (fetcher *SFetcher) FetchBlockMetric(Type int, Counter int64) {
	clientmetric := fetcher.ClientMetric
	blockmetric := fetcher.BlockMetric
	txindex := fetcher.GetCurTxindex()
	txinfo := &fetcher.BlockMetric.TxInfo[txindex]
	if Type == 1 {
		clientmetric.AccountHashCounter += Counter
		blockmetric.AccountHashCounter += Counter
		txinfo.AccountHashCounter += Counter
	} else if Type == 2 {
		clientmetric.AccountHashRawCounter += Counter
		blockmetric.AccountHashRawCounter += Counter
		txinfo.AccountHashRawCounter += Counter
	} else if Type == 3 {
		clientmetric.StorageHashCounter += Counter //node change
		blockmetric.StorageHashCounter += Counter
		txinfo.StorageHashCounter += Counter
	} else if Type == 4 {
		clientmetric.StorageHashRawCounter += Counter //trie change ,update delete
		blockmetric.StorageHashRawCounter += Counter
		txinfo.StorageHashRawCounter += Counter
	} else if Type == 8 {
		clientmetric.BlockhashInsCounter += Counter
		blockmetric.BlockhashInsCounter += Counter
		txinfo.BlockhashInsCounter += Counter
	} else if Type == 11 {
		clientmetric.ContractModifyCounter += Counter
		blockmetric.ContractModifyCounter += Counter
	} else if Type == 12 {
		blockmetric.AccountDepthDistribution[Counter] += 1
	} else if Type == 13 {
		blockmetric.StorageDepthDistribution[Counter] += 1
	} else {
		log.Error("hzyserr", "critical bug", Type, Counter)
		os.Exit(0)
	}
}

func (fetcher *SFetcher) FetchTxMetric(Type int, Counter int64) {
	txindex := fetcher.GetCurTxindex()
	txinfo := &fetcher.BlockMetric.TxInfo[txindex]
	if Type == 1 {
		txinfo.Gasused = Counter
	} else if Type == 2 {
		txinfo.TxType = Counter
	} else {
		log.Error("hzyserr", "critical bug", Type, Counter)
		os.Exit(0)
	}
}

func (fetcher *SFetcher) FetchUnderlyingDbMetric(BlockNum int64, writevolume float64, readvolume float64) {
	var underlyingDbInfo SUnderlyingDBMetric
	underlyingDbInfo.CurBlockNum = BlockNum
	underlyingDbInfo.DiskWriteMeter = writevolume
	underlyingDbInfo.DiskReadMeter = readvolume
	fetcher.ChUnderlyingDBMetric <- underlyingDbInfo
}

func (fetcher *SFetcher) FetchOpcodeMetric(
	opcodecommoninfo SOpcodecommon,
	BytecodeExsit bool,
	transfer bool,
	accountsas bool,
	storagesas bool,
) {
	var OpcodeInfo SOpcodeMetric
	OpcodeInfo.Opcodeinfo = opcodecommoninfo
	OpcodeInfo.BytecodeExsit = BytecodeExsit
	OpcodeInfo.Transfer = transfer
	OpcodeInfo.AccountSAS = accountsas
	OpcodeInfo.StorageSAS = storagesas
	fetcher.ChOpcodeMetric <- OpcodeInfo
}

func (fetcher *SFetcher) FetchOverheadMetric(BlockNum int64, processtime time.Duration, alloc uint64) {
	var overheadinfo SOverheadMetric
	overheadinfo.CurBlockNum = BlockNum
	overheadinfo.ProcessTime = processtime
	overheadinfo.Alloc = alloc
	fetcher.ChOverheadMetric <- overheadinfo
}
