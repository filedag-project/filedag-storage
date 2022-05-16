package node

//func (d DagNode) recoveryDisk(path, newPath string) error {
//	keyCodeMap, err := d.db.ReadAll("")
//	if err != nil {
//		return err
//	}
//	sliceNode := new(SliceNode)
//	for _, node := range d.nodes {
//		if node.cfg.Path == path {
//			sliceNode = node
//		}
//	}
//	sliceNewNode, err := NewSliceNode(config.CaskNumConf(int(sliceNode.cfg.CaskNum)), config.PathConf(newPath))
//	for key, value := range keyCodeMap {
//		keyCode := sha256String(key)
//		id := d.fileID(keyCode)
//		cask, has := sliceNode.caskMap.Get(id)
//		if !has {
//			return kv.ErrNotFound
//		}
//		length, err := cask.Size(keyCode)
//		if err != nil {
//			return err
//		}
//		if length > 0 {
//			continue
//		}
//		merged := make([][]byte, 0)
//		index := -1
//		for i, node := range d.nodes {
//			cask, has = node.caskMap.Get(id)
//			if !has {
//				fmt.Println("********")
//				return kv.ErrNotFound
//			}
//			bytes, err := cask.Read(keyCode)
//			if err != nil || len(bytes) == 0 {
//				index = i
//				return err
//			}
//			merged = append(merged, bytes)
//		}
//		if index == -1 {
//			continue
//		}
//		i64, err := strconv.ParseInt(value, 10, 64)
//		if err == nil {
//			log.Errorf("strconv fail :%v", err)
//		}
//		enc, err := NewErasure(d.dataBlocks, d.parityBlocks, i64)
//		enc.DecodeDataBlocks(merged)
//		dataByte := merged[index]
//		cask, has = sliceNewNode.caskMap.Get(id)
//		if !has {
//			done := make(chan error)
//			sliceNewNode.createCaskChan <- &createCaskRequst{
//				id:   id,
//				done: done,
//			}
//			if err := <-done; err != ErrNone {
//				return err
//			}
//			cask, _ = sliceNewNode.caskMap.Get(id)
//		}
//		err = cask.Put(keyCode, dataByte)
//		if err != nil {
//			break
//		}
//	}
//	return err
//}
//
//func (d *DagNode) modifyConfig(oldPath, newPath string, newCaskNum uint32) error {
//	defer errors.New("modifyConfig error")
//	for _, node := range d.nodes {
//		if node.cfg.Path == oldPath {
//			node.cfg.Path = newPath
//			node.cfg.CaskNum = newCaskNum
//		}
//	}
//	return nil
//}
//
//func (d DagNode) recoveryHost(newPath string) error {
//	keyCodeMap, err := d.db.ReadAll("")
//	if err != nil {
//		return err
//	}
//	sliceNewNode := new(SliceNode)
//	for _, node := range d.nodes {
//		if node.cfg.Path == newPath {
//			sliceNewNode = node
//		}
//	}
//	for key, value := range keyCodeMap {
//		keyCode := sha256String(key)
//		id := d.fileID(keyCode)
//		cask, has := sliceNewNode.caskMap.Get(id)
//		if !has {
//			return kv.ErrNotFound
//		}
//		merged := make([][]byte, 0)
//		index := -1
//		for i, node := range d.nodes {
//			cask, has = node.caskMap.Get(id)
//			if !has {
//				fmt.Println("********")
//				return kv.ErrNotFound
//			}
//			bytes, err := cask.Read(keyCode)
//			if err != nil || len(bytes) == 0 {
//				index = i
//				return err
//			}
//			merged = append(merged, bytes)
//		}
//		if index == -1 {
//			continue
//		}
//		i64, err := strconv.ParseInt(value, 10, 64)
//		if err == nil {
//			log.Errorf("strconv fail :%v", err)
//		}
//		enc, err := NewErasure(d.dataBlocks, d.parityBlocks, i64)
//		enc.DecodeDataBlocks(merged)
//		dataByte := merged[index]
//		cask, has = sliceNewNode.caskMap.Get(id)
//		if !has {
//			done := make(chan error)
//			sliceNewNode.createCaskChan <- &createCaskRequst{
//				id:   id,
//				done: done,
//			}
//			if err := <-done; err != ErrNone {
//				return err
//			}
//			cask, _ = sliceNewNode.caskMap.Get(id)
//		}
//		err = cask.Put(keyCode, dataByte)
//		if err != nil {
//			break
//		}
//	}
//	return err
//}
