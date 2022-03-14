// diskv - dag 数据块存储
//
// merkle-dag 分两类： 1 包含实际数据，体量较大 2 不包含数据数据，只有链接，体量较小
//   大体量 dag 直接落磁盘，在 leveldb 里只保存大小等信息
// 	 小体量 dag 保持到 leveldb
//
//

package diskv
