package storage

type Storer interface {
	Lister
	Deleter
}

type Lister interface {
	List(ext string) ([]string, error)
}

type Deleter interface {
	Delete(key string) error
	DeleteMulti(keys []string) error
}

// Batcher groups strings into batches of a given size
func Batcher(items []string, size int) [][]string {
	totalBatches := float64(len(items) / size)
	if len(items)%size != 0 {
		totalBatches++
	}
	var batched = make([][]string, int(totalBatches))
	var batch []string
	for i, item := range items {
		batchNum := int(i / size)
		if i%size == 0 {
			batch = nil
		}
		batch = append(batch, item)
		batched[batchNum] = batch
	}
	return batched
}
