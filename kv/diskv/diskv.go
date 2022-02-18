package diskv

type DisKV struct {
}

func (di *DisKV) Put(key []byte, value []byte) error {
	return nil
}

func (di *DisKV) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (di *DisKV) Delete(key []byte) error {
	return nil
}
