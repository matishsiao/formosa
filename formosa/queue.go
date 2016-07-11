package formosa

import (
	"fmt"
	"log"
)

func (dm *DBManager) QueuePush(queue string, value string) error {
	sizeKey := EncodeQueueEnd(queue)
	log.Println("QueuePush", queue, sizeKey, value)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return err
	}
	var size int64 = 0
	exists, _ := transaction.Has([]byte(sizeKey), nil)
	if !exists {
		transaction.Put([]byte(sizeKey), []byte("0"), nil)
	} else {
		data, _ := transaction.Get([]byte(sizeKey), nil)
		size = ToInt64(string(data))
	}

	if size >= QUEUE_SIZE {
		transaction.Commit()
		return fmt.Errorf("Queue size has full.")
	}
	var rear int64 = 0
	rearKey := EncodeQueueRear(queue)
	exists, _ = transaction.Has([]byte(rearKey), nil)
	if exists {
		data, _ := transaction.Get([]byte(rearKey), nil)
		rear = ToInt64(string(data))
	}
	rear++
	if rear >= QUEUE_SIZE {
		rear = 0
	}
	//Update rearSeq to DB
	transaction.Put([]byte(rearKey), []byte(fmt.Sprintf("%d", rear)), nil)
	//Push data to queue
	rearSeq := EncodeQueue(queue, PaddingLeft(fmt.Sprintf("%d", rear), 9))
	transaction.Put([]byte(rearSeq), []byte(value), nil)
	//Update queueSize to DB
	size++
	transaction.Put([]byte(sizeKey), []byte(fmt.Sprintf("%d", size)), nil)
	log.Println(size, rear, rearKey, rearSeq)
	return transaction.Commit()
}

func (dm *DBManager) QueuePop(queue string) ([]byte, error) {
	sizeKey := EncodeQueueEnd(queue)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return []byte{}, err
	}
	var size int64 = 0
	exists, _ := transaction.Has([]byte(sizeKey), nil)
	if !exists {
		transaction.Put([]byte(sizeKey), []byte("0"), nil)
	} else {
		data, _ := transaction.Get([]byte(sizeKey), nil)
		size = ToInt64(string(data))
	}

	if size == 0 {
		transaction.Commit()
		return []byte{}, fmt.Errorf("Queue size has empty.")
	}

	var front int64 = 0
	frontKey := EncodeQueueFront(queue)
	exists, _ = transaction.Has([]byte(frontKey), nil)
	if exists {
		data, _ := transaction.Get([]byte(frontKey), nil)
		front = ToInt64(string(data))
	}
	front++
	if front >= QUEUE_SIZE {
		front = 0
	}
	//Update frontSeq to DB
	transaction.Put([]byte(frontKey), []byte(fmt.Sprintf("%d", front)), nil)
	//Get data from queue
	frontSeq := EncodeQueue(queue, PaddingLeft(fmt.Sprintf("%d", front), 9))
	value, _ := transaction.Get([]byte(frontSeq), nil)
	//Remove data from queue
	transaction.Delete([]byte(frontSeq), nil)
	//Update size to DB
	size--
	transaction.Put([]byte(sizeKey), []byte(fmt.Sprintf("%d", size)), nil)
	transaction.Commit()
	return value, nil
}

func (dm *DBManager) QueueSize(queue string) (int64, error) {
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return 0, err
	}
	size, err := transaction.Get([]byte(EncodeQueueEnd(queue)), nil)
	transaction.Commit()
	return ToInt64(string(size)), err
}
