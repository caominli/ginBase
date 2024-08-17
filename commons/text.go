// 这个文件存放处理文本的函数
package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// 字符串转大写，用于验证码可以不区分大小写
// 参数：str ，返回：str
func DaXie(code string) (code2 string) {
	return strings.ToUpper(code)
}

// 生成随机码,参数:长度
func Captcha(length int) string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZ123456789"
	// const charset = "ABC"
	var sb strings.Builder
	sb.Grow(length)

	// 创建一个伪随机数生成器的种子源
	source := rand.NewSource(time.Now().UnixNano())

	// 使用种子源创建一个伪随机数生成器
	random := rand.New(source)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}

// 生成订单号
var counter int64

func GetOrderID() uint64 {
	// 使用时间戳
	timestamp := time.Now().Unix() // 获取当前时间戳
	// 使用原子计数器
	uniqueNum := atomic.AddInt64(&counter, 1) % 100
	jieguo := fmt.Sprintf("%d%03d", timestamp, uniqueNum) //拼接字符串
	num, _ := strconv.ParseUint(jieguo, 10, 64)           //字符串转uint64
	return num
}
