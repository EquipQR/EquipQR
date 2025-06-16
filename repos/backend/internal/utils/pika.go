package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var PikaGenerator *Generator

const (
	epoch        int64 = 1685577600000 // Custom epoch: June 1, 2023 (ms)
	nodeIDBits         = 10
	sequenceBits       = 12
	maxNodeID          = -1 ^ (-1 << nodeIDBits)
	maxSequence        = -1 ^ (-1 << sequenceBits)
	timeShift          = nodeIDBits + sequenceBits
	nodeIDShift        = sequenceBits
)

type Generator struct {
	nodeID   int64
	lastTS   int64
	sequence int64
	mutex    sync.Mutex
}

func InitPikaGenerator(nodeID int64) error {
	generatorInstance, err := NewGenerator(nodeID)

	if err != nil {
		return err
	}

	PikaGenerator = generatorInstance
	return nil
}

func NewGenerator(nodeID int64) (*Generator, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, errors.New("node ID out of range")
	}
	return &Generator{nodeID: nodeID}, nil
}

// Generate a string ID like "user_18237182317123"
func (g *Generator) NextID(prefix string) string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := currentMillis()
	if now == g.lastTS {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastTS {
				now = currentMillis()
			}
		}
	} else {
		g.sequence = 0
	}
	g.lastTS = now

	id := ((now - epoch) << timeShift) | (g.nodeID << nodeIDShift) | g.sequence
	return fmt.Sprintf("%s_%d", prefix, id)
}

func currentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
