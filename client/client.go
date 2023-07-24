package main

import (
	"fmt"
	"net"
	"os"
	"github.com/google/uuid"
)

const Info = `
==================================== 操作方法 ===========================================
钱包操作:
 getPubKey      --nodeAddr 容器                                      "获取节点的公钥"
 getPrivateKey  --nodeAddr 容器                                      "获取节点的私钥"
 getAllPubKey   --nodeAddr 容器                                      "获取所有节点的公钥"
 getUserInfo    --nodeAddr 容器  --nodeID 节点编号                   "查看区块链上的账户信息"
交易操作:
 sendNorTrans   --nodeAddr 容器  --from 发起者 --to 接受者 --value 金额   "普通交易"
 sendGeneTrans  --nodeAddr 容器  --from 发起者 --to 接受者 --data  数据   "分发交易"(待完善)
 sendTransTrans --nodeAddr 容器  --from 发起者 --to 接受者 --data  数据   "转移交易"(待完善)
查看区块链:
 printBlockInfo   --nodeAddr 容器  --round 轮次                       "查看区块数据"
 printTransPool   --nodeAddr 容器                                     "查看交易池"
 printTransNum    --nodeAddr 容器                                     "查看每轮已发生交易数"
 printTransNumAll --nodeAddr 容器                                     "查看系统中发生的交易数"
 printMiniPool    --nodeAddr 容器                                     "查看微块池"
 printHeight      --nodeAddr 容器                                     "查看区块高度"
==========================================================================================
`

func main() {
	// 获取客户端参数
	args := os.Args // go run client.go getPubKey --nodeID 8000

	// 判断
	if len(args) < 2 {
		fmt.Print(Info)
		return
	}

	// 分析命令
	cmd := args[1]

	switch cmd {
	case "getPubKey":
		// 标识: d
		getPubKey(args)
	case "getPrivateKey":
		// 标识: d
		getPrivateKey(args)
	case "getAllPubKey":
		getAllPubKey(args)
	case "sendNorTrans":
		// 标识: a
		sendNorTrans(args)
	case "sendGeneTrans":
		// 标识: b
	case "sendTransTrans":
		// 标识: c
	case "printBlockInfo":
		// 标识: d
		printBlockInfo(args)
	case "printTransPool":
		printTransPool(args)
	case "printTransNum":
		printTransNum(args)
	case "printTransNumAll":
		printTransNumAll(args)
	case "printMiniPool":
		printMiniPool(args)
	case "printHeight":
		printHeight(args)
	case "getUserInfo":
		getUserInfo(args)
	default:
		fmt.Print("无效的命令，请重新输入...")
		fmt.Print(Info)
	}
}

// 1 获取节点对应的公钥
// go run client.go getPubKey --nodeID 8000
func getPubKey(args []string) {
	// 组织数据，按照规则进行拆分
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" +"getPubKey" + "_" + args[3] // 拼接字符串 d_f677689a-9631-4897-8156-3053f63d573b_getPubKey_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 2 获取节点对应的私钥
func getPrivateKey(args []string) {
	// 组织数据，按照规则进行拆分
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "getPrivateKey" + "_" + args[3] // 拼接字符串 d_getPubKey_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 3 获取所有节点的公钥
func getAllPubKey(args []string) {
	// 组织数据，按照规则进行拆分
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "getAllPubKey" + "_" + args[3] // 拼接字符串 d_getPubKey_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 4 打印区块链信息
// printBlockInfo --nodeID 节点 --round 轮次
func printBlockInfo(args []string) {
	// 组织数据，按照规则进行拆分
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	// 拼接字符串 d_printBlockInfo_8000_0
	data := "d" + "_" + uuid.String() + "_" + "printBlockInfo" + "_" + args[3] + "_" + args[5]

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 正常交易
func sendNorTrans(args []string) {
	// sendNorTrans   --nodeID 节点 --from 发起人 --to 接收人 --value 金额
	// 拼接字符串
	uuid := uuid.New()
	data := "a" + "_" + uuid.String() + "_" + args[5] + "_" + args[7] + "_" + args[9]
	//      "a"          "from"          "to"          "value"    

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看交易池
func printTransPool(args []string) {
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "printTransPool" + "_" + args[3] // 拼接字符串 d_printTransPool_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看全局交易数目
func printTransNum(args []string) {
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "printTransNum" + "_" + args[3] // 拼接字符串 d_printTransNum_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看系统中发生的总交易数目
func printTransNumAll(args []string) {
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "printTransNumAll" + "_" + args[3] // 拼接字符串 d_printTransNumAll_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看微块池
func printMiniPool(args []string) {
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "printMiniPool" + "_" + args[3] // 拼接字符串 d_printTransPool_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看区块高度
func printHeight(args []string) {
	// 组织数据，按照规则进行拆分
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "printHeight" + "_" + args[3] // 拼接字符串 d_printHeight_8000

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 查看用户信息
func getUserInfo(args []string) {
	// 字符串拼接: 使用 _ 进行划分
	uuid := uuid.New()
	data := "d" + "_" + uuid.String() + "_" + "getUserInfo" + "_" + args[3] + "_" + args[5] // 拼接字符串 d_getUserInfo_8000_公钥

	// 主动发送连接请求
	conn, err := net.Dial("tcp", args[3] + ":8989")

	if err != nil {
		fmt.Println("Dial err", err)
	}

	defer conn.Close() // 客户端终止时，关闭于服务器通讯的socket

	// 将命令行中的数据，发送给客户端
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}