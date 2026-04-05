Fields
- IntField: 4 bytes, big endian unsigned
- StringField: first 4 bytes = actual string length as a big endian int, then
  the string bytes, then null bytes padding up to maxLength. Total bytes on disk
  = 4 + maxLength

TupleDesc
- Just metadata, nothing on disk. Describes the schema: ordered list of (field
  name, type) pairs
- For string types it also needs to carry the maxLength

Tuple
- On disk: fields serialized back to back, no separators, fixed total size per
  schema
- In memory: slice of Field values + a RecordID

RecordID
- (tableID, pageID, slotIndex) — uniquely identifies where a tuple lives

HeapPage (4096 bytes)
- Byte 0 to headerSize-1: bitmap header, 1 bit per slot, LSB first within each
  byte
- Byte headerSize onwards: tuple slots, each exactly tupleSize bytes, fixed
  positions regardless of occupancy
- tupleSize = sum of all field sizes
- numSlots = floor(4096 * 8 / (tupleSize * 8 + 1))
- headerSize = ceil(numSlots / 8)

HeapFile
- Flat binary file, pages concatenated sequentially
- Page n is at byte offset n * 4096
- Number of pages = fileSize / 4096

BufferPool
- Fixed capacity (number of pages, not bytes)
- On GetPage(tableID, pageID): return cached page if present, otherwise read  
  from disk via the catalog, cache it, return it
- No eviction needed for Lab 1

Catalog
- Singleton, lives in memory only
- Maps tableID → (DBFile, TupleDesc, name)
- Maps name → tableID
- tableID can just be a hash of the file path or an incrementing int

SeqScan
- Takes a tableID and opens an iterator over the HeapFile
- Walks pages 0 → N, within each page walks occupied slots only
- Implements the same TupleIterator interface as HeapPage