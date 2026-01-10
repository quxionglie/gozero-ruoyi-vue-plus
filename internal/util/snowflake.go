package util

import (
	"sync"
	"time"
)

const (
	// 起始时间戳 (2024-01-01 00:00:00)
	epoch int64 = 1704067200000

	// 机器ID占用的位数
	workerIDBits int64 = 5

	// 数据标识ID占用的位数
	datacenterIDBits int64 = 5

	// 序列号占用的位数
	sequenceBits int64 = 12

	// 机器ID最大值 (2^5 - 1 = 31)
	maxWorkerID int64 = -1 ^ (-1 << workerIDBits)

	// 数据标识ID最大值 (2^5 - 1 = 31)
	maxDatacenterID int64 = -1 ^ (-1 << datacenterIDBits)

	// 机器ID向左移12位
	workerIDShift int64 = sequenceBits

	// 数据标识ID向左移17位(12+5)
	datacenterIDShift int64 = sequenceBits + workerIDBits

	// 时间戳向左移22位(12+5+5)
	timestampLeftShift int64 = sequenceBits + workerIDBits + datacenterIDBits

	// 生成序列的掩码，这里为4095 (0b111111111111=0xfff=4095)
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)
)

// Snowflake 雪花算法生成器
type Snowflake struct {
	mutex         sync.Mutex // 互斥锁
	workerID      int64      // 机器ID
	datacenterID  int64      // 数据标识ID
	sequence      int64      // 序列号
	lastTimestamp int64      // 上次生成ID的时间戳
}

var (
	// 默认雪花算法实例（单例模式）
	defaultSnowflake *Snowflake
	once             sync.Once
)

// NewSnowflake 创建雪花算法实例
// workerID: 机器ID (0-31)
// datacenterID: 数据标识ID (0-31)
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, ErrInvalidDatacenterID
	}

	return &Snowflake{
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

// InitDefaultSnowflake 初始化默认雪花算法实例
// workerID: 机器ID (0-31)，默认0
// datacenterID: 数据标识ID (0-31)，默认0
func InitDefaultSnowflake(workerID, datacenterID int64) error {
	var err error
	once.Do(func() {
		defaultSnowflake, err = NewSnowflake(workerID, datacenterID)
	})
	return err
}

// GetDefaultSnowflake 获取默认雪花算法实例
func GetDefaultSnowflake() *Snowflake {
	if defaultSnowflake == nil {
		// 如果没有初始化，使用默认值（0, 0）初始化
		_ = InitDefaultSnowflake(0, 0)
	}
	return defaultSnowflake
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixMilli()

	// 如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过，此时应当抛出异常
	if timestamp < s.lastTimestamp {
		return 0, ErrClockBackwards
	}

	// 如果是同一时间生成的，则进行毫秒内序列
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		// 毫秒内序列溢出
		if s.sequence == 0 {
			// 阻塞到下一个毫秒，获得新的时间戳
			timestamp = s.tilNextMillis(s.lastTimestamp)
		}
	} else {
		// 时间戳改变，毫秒内序列重置
		s.sequence = 0
	}

	// 上次生成ID的时间戳
	s.lastTimestamp = timestamp

	// 移位并通过或运算拼到一起组成64位的ID
	id := ((timestamp - epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}

// tilNextMillis 阻塞到下一个毫秒，直到获得新的时间戳
func (s *Snowflake) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixMilli()
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixMilli()
	}
	return timestamp
}

// GenerateID 使用默认雪花算法实例生成ID（便捷函数）
func GenerateID() (int64, error) {
	return GetDefaultSnowflake().NextID()
}

// 错误定义
var (
	ErrInvalidWorkerID     = &SnowflakeError{Message: "worker ID must be between 0 and 31"}
	ErrInvalidDatacenterID = &SnowflakeError{Message: "datacenter ID must be between 0 and 31"}
	ErrClockBackwards      = &SnowflakeError{Message: "clock moved backwards, refusing to generate id"}
)

// SnowflakeError 雪花算法错误
type SnowflakeError struct {
	Message string
}

func (e *SnowflakeError) Error() string {
	return e.Message
}
