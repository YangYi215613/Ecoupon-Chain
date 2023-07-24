package main

import (
	"fmt"
	"strconv"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)


const protocol = "tcp"  // 通讯协议
var networkAddr string
var nodeNumber int   // 系统中总结点数(通过.env传入)
var nodeName string  // 节点名 node1
var allNodeName []string // 所有节点名
var committeeNum int // 委员会数目
var nodeCommitteeID int  // 委员会ID
var isCommitteeLeader bool
const maxConn = 15                  // 节点最大连接数
var genesisTimeStamp uint64         // 创世区块时间戳
var FindNodes []string = make([]string, maxConn)         // 节点发现，存放该节点发现的节点 [ip:port, ip:port, ...]


// ******************** 统计 **********************
var currentTransNum uint64       // 每轮交易数
var allTransNum uint64       // 总交易数
const maxTransNumPerRound = 10 // 每轮最大处理交易数目

// ***************** 内存池设置 ********************
// 交易内存池:存放节点收到的交易
var transMemPool []*NorTransaction
var allTransMemInfo []*NorTransaction  // 存放所有交易的信息


// 微块内存池: 存放各节点收到的微块数目(如果满足条件，节点打包生成完整区块)
var nodeAcceptMiniBlock []*MiniBlock

// ***************** 数据库设置 ********************
const blockChainDB = "blockChain.db"  // 区块链数据库
const blockBucket = "bucket_node"     // 数据库中表名
const lastBlockHash = "LastBlockHash" // 区块链中最后一个区块hash的索引

// 区块链对象
var bc *BlockChain


// ***************** 日志设置 ********************
// Level 日志级别。建议从服务配置读取。
var LogConf = struct {
	Dir     string `yaml:"dir"`
	Name    string `yaml:"name"`
	Level   string `yaml:"level"`
	MaxSize int    `yaml:"max_size"`
}{
	Dir:     "./logs",
	Name:    "logs.log",
	Level:   "trace",
	MaxSize: 100,
}

type MyFormatter struct {
	logrus.Formatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 自定义日志格式
	function := entry.Caller.Function
	file := entry.Caller.File
	line := entry.Caller.Line

	// 自定义日志格式
	return []byte(entry.Time.Format("2006-01-02 15:04:05") + " [" + entry.Level.String() + "] " + function + " " + file + ":" + strconv.Itoa(line) + " [" + entry.Message + "]\n"), nil
}

// Init logrus logger.
func InitLogger() error {
	// 设置日志格式。
	logrus.SetFormatter(&MyFormatter{})

	logrus.SetLevel(logrus.InfoLevel)

	// 实现日志滚动。
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v", LogConf.Dir, LogConf.Name), // 日志输出文件路径。
		MaxSize:    LogConf.MaxSize,                                 // 日志文件最大 size(MB)，缺省 100MB。
		MaxBackups: 10,                                              // 最大过期日志保留的个数。
		MaxAge:     30,                                              // 保留过期文件的最大时间间隔，单位是天。
		LocalTime:  true,                                            // 是否使用本地时间来命名备份的日志。
	}
	logrus.SetReportCaller(true)
	logrus.SetOutput(io.MultiWriter(os.Stdout, logger))
	return nil
}