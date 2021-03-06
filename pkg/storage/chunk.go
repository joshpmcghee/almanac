package storage

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"dinowernli.me/almanac/pkg/index"
	"dinowernli.me/almanac/pkg/util"
	pb_almanac "dinowernli.me/almanac/proto"

	"golang.org/x/net/context"
)

const (
	chunkIdFormat    = "%3s-%d-%d-%s"
	chunkIdSeparator = "-"
	chunkUidLength   = 5
)

var (
	chunkTypeString = map[pb_almanac.ChunkId_Type]string{pb_almanac.ChunkId_SMALL: "sml", pb_almanac.ChunkId_BIG: "big"}
	chunkTypeValue  = map[string]pb_almanac.ChunkId_Type{"sml": pb_almanac.ChunkId_SMALL, "big": pb_almanac.ChunkId_BIG}
)

// ChunkId returns the string representation of the supplied chunk id proto.
func ChunkId(idProto *pb_almanac.ChunkId) (string, error) {
	if idProto.Uid == "" {
		return "", fmt.Errorf("cannot create chunk id with empty uid")
	}

	if strings.Contains(idProto.Uid, chunkIdSeparator) {
		return "", fmt.Errorf("chunk uid cannot contain '-', but got: %s", idProto.Uid)
	}

	if idProto.StartMs > idProto.EndMs {
		return "", fmt.Errorf("invalid start and end times: start=%d, end=%d", idProto.StartMs, idProto.EndMs)
	}

	typeString, ok := chunkTypeString[idProto.Type]
	if !ok {
		return "", fmt.Errorf("unknown chunk type: %v", idProto.Type)
	}

	return fmt.Sprintf(chunkIdFormat, typeString, idProto.StartMs, idProto.EndMs, idProto.Uid), nil
}

// ChunkIdProto returns the structured representation of the supplied chunk id.
func ChunkIdProto(chunkId string) (*pb_almanac.ChunkId, error) {
	var uid string
	var startMs int64
	var endMs int64
	var chunkType string
	_, err := fmt.Sscanf(chunkId, chunkIdFormat, &chunkType, &startMs, &endMs, &uid)
	if err != nil {
		return nil, fmt.Errorf("unable to parse id [%s] with format [%s]: %v", chunkId, chunkIdFormat, err)
	}

	typeEnum, ok := chunkTypeValue[chunkType]
	if !ok {
		return nil, fmt.Errorf("unknown chunk type string: %s", chunkType)
	}

	return &pb_almanac.ChunkId{
		StartMs: startMs,
		EndMs:   endMs,
		Uid:     uid,
		Type:    typeEnum,
	}, nil
}

// ChunkProto is a one-stop-shop for creating a chunk proto from a set of entries.
func ChunkProto(entries []*pb_almanac.LogEntry, chunkType pb_almanac.ChunkId_Type) (*pb_almanac.Chunk, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("cannot create chunk proto for zero entries")
	}

	idx, err := index.NewIndex()
	if err != nil {
		return nil, fmt.Errorf("unable to create index: %v", err)
	}

	var minMs int64 = math.MaxInt64
	var maxMs int64 = math.MinInt64
	for _, e := range entries {
		if e.TimestampMs < minMs {
			minMs = e.TimestampMs
		}
		if e.TimestampMs > maxMs {
			maxMs = e.TimestampMs
		}

		err := idx.Index(e.Id, e.EntryJson)
		if err != nil {
			return nil, fmt.Errorf("unable to index entry: %v", err)
		}
	}

	chunkId := &pb_almanac.ChunkId{
		StartMs: minMs,
		EndMs:   maxMs,
		Uid:     NewChunkUid(),
		Type:    chunkType,
	}

	idxProto, err := index.Serialize(idx)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize index: %v", err)
	}

	return &pb_almanac.Chunk{Entries: entries, Id: chunkId, Index: idxProto}, nil
}

// NewChunkId returns a string which can be used as the "uid" part of a new chunk.
func NewChunkUid() string {
	return util.RandomString(chunkUidLength)
}

// Chunk is an in-memory, immutable representation of a stored chunk. A chunk
// must be closed by calling Close() once it is no longer in use.
type Chunk struct {
	id       *pb_almanac.ChunkId
	index    *index.Index
	entryMap map[string]*pb_almanac.LogEntry
	entries  []*pb_almanac.LogEntry
	closed   bool
}

// openChunk returns a chunk instance for the supplied proto. The caller is
// responsible for calling Close() on the returned instance.
func openChunk(chunkProto *pb_almanac.Chunk) (*Chunk, error) {
	chunkId, err := ChunkId(chunkProto.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to extract string id for chunk: %v", err)
	}

	idx, err := index.Deserialize(chunkProto.Index)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize index from chunk %s: %v", chunkId, err)
	}

	entryMap := map[string]*pb_almanac.LogEntry{}
	for _, entry := range chunkProto.Entries {
		entryMap[entry.Id] = entry
	}
	return &Chunk{id: chunkProto.Id, index: idx, entryMap: entryMap, entries: chunkProto.Entries}, nil
}

// Search returns all log entries in the chunk matching the supplied query, in ascending order by timestamp.
func (c *Chunk) Search(ctx context.Context, query string, num int32, startMs int64, endMs int64) ([]*pb_almanac.LogEntry, error) {
	if c.closed {
		return nil, fmt.Errorf("cannot execute search on closed chunk")
	}
	return Search(ctx, c.index, c.entryMap, query, num, startMs, endMs)
}

// Entries returns all the entries in this chunk. Callers must not modify the return value.
func (c *Chunk) Entries() []*pb_almanac.LogEntry {
	return c.entries
}

// Id returns the id of this chunk. Callers must not modify the return value.
func (c *Chunk) Id() *pb_almanac.ChunkId {
	return c.id
}

// Close releases any resources associated with this chunk.
func (c *Chunk) Close() error {
	err := c.index.Close()
	if err != nil {
		return fmt.Errorf("unable to close index: %v", err)
	}
	c.closed = true
	return nil
}

// Search executes a search on a given index and entry map. Results are returned in ascending order by timestamp.
func Search(ctx context.Context, idx *index.Index, entries map[string]*pb_almanac.LogEntry, query string, num int32, startMs int64, endMs int64) ([]*pb_almanac.LogEntry, error) {
	ids, err := idx.Search(ctx, query, num)
	if err != nil {
		return nil, fmt.Errorf("unable to search index: %v", err)
	}

	result := []*pb_almanac.LogEntry{}
	for _, id := range ids {
		entry, ok := entries[id]
		if !ok {
			return nil, fmt.Errorf("could not locate hit %s", id)
		}

		if startMs != 0 && entry.TimestampMs < startMs {
			continue
		}
		if endMs != 0 && entry.TimestampMs > endMs {
			continue
		}
		result = append(result, entry)
		if int32(len(result)) >= num {
			break
		}
	}

	// TODO(dino): Figure out if there is a way to get bleve to return these in sorted order.
	if !sort.IsSorted(OldestEntryFirst(result)) {
		sort.Sort(OldestEntryFirst(result))
	}

	return result, nil
}

func (c *Chunk) fetchEntry(id string) (*pb_almanac.LogEntry, error) {
	result, ok := c.entryMap[id]
	if !ok {
		return nil, fmt.Errorf("entry %s not found", id)
	}
	return result, nil
}
