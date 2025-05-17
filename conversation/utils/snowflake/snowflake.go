package snowflake

import (
	"context"
	"fmt"
	"sync"
	"time"

	"conversation/db"
)

// TODO：需要考虑41位时间戳可以用多久哈
// 雪花算法：1 Bit（符号位） | 41 Bit（时间戳） | 10 Bit（机器ID） | 12 Bit（序列号） |
const (
	epoch          = int64(1609459200000)                     // 自定义起始时间（2021-01-01 00:00:00 UTC）
	machineIDBits  = uint(10)                                 // 机器 ID 位数
	sequenceBits   = uint(12)                                 // 序列号位数
	machineIDShift = sequenceBits                             // 机器 ID 左移位数
	timestampShift = sequenceBits + machineIDBits             // 时间戳左移位数
	maxMachineID   = int64(-1) ^ (int64(-1) << machineIDBits) // 最大机器 ID
	maxSequence    = int64(-1) ^ (int64(-1) << sequenceBits)  // 最大序列号
)

type Snowflake struct {
	mu            sync.Mutex
	machineID     int64
	sequence      int64
	lastTimestamp int64
}

func NewSnowflake() (*Snowflake, error) {
	//machineID可以直接从redis incr获取
	ctx := context.Background()
	machineID, err := db.Rdb.Incr(ctx, "machine_id").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %v", err)
	}
	// 检查机器 ID 是否超出范围
	if machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID exceeds max limit: %d", maxMachineID)
	}
	return &Snowflake{machineID: machineID}, nil
}

func (s *Snowflake) Generate() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // 毫秒时间戳

	// 如果当前时间小于上次生成 ID 的时间，说明系统时钟回退了
	// TODO: 如果回退可以考虑重新获得一个machineID？
	if now < s.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards. Refusing to generate id for %d milliseconds", s.lastTimestamp-now)
	}
	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 { // 当前毫秒的序列号用完，等待下一毫秒
			for now <= s.lastTimestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = now
	return (now-epoch)<<timestampShift | (s.machineID << machineIDShift) | s.sequence, nil
}
