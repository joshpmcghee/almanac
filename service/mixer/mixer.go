package mixer

import (
	"fmt"
	"sort"

	pb_almanac "dinowernli.me/almanac/proto"
	"dinowernli.me/almanac/service/discovery"
	"dinowernli.me/almanac/storage"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Mixer is an implementation of the mixer rpc service. It provides global
// search functionality across the entire system.
type Mixer struct {
	storage   *storage.Storage
	discovery *discovery.Discovery
}

// New returns a new mixer backed by the supplied storage.
func New(storage *storage.Storage, discovery *discovery.Discovery) *Mixer {
	return &Mixer{storage: storage, discovery: discovery}
}

func (m *Mixer) Search(ctx context.Context, request *pb_almanac.SearchRequest) (*pb_almanac.SearchResponse, error) {
	// Do some prep for the parallel searches.
	appenders := m.discovery.ListAppenders()
	chunkIds, err := m.storage.ListChunks(request.StartMs, request.EndMs)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "unable to list chunks: %v", err)
	}

	// Execute all the searches in parallel.
	numSubRequests := len(appenders) + len(chunkIds)
	resultChan := make(chan *partialResult, numSubRequests)
	for _, chunkId := range chunkIds {
		go m.searchChunk(ctx, chunkId, request.Query, request.Num, resultChan)
	}
	for _, appender := range appenders {
		go m.searchAppender(ctx, appender, request, resultChan)
	}

	// Drain the channel and collect all the entries.
	allEntries := []*pb_almanac.LogEntry{}
	err = nil
	for i := 0; i < numSubRequests; i++ {
		result := <-resultChan
		if result.chunk != nil {
			result.chunk.Close()
		}

		if result.err == nil {
			for _, e := range result.entries {
				allEntries = append(allEntries, e)
			}
		} else {
			err = result.err
		}
	}

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "error executing search: %v", err)
	}

	// Sort all the entries by timestamp and take the first "num" distinct ones.
	sort.Sort(byTimestamp(allEntries))
	result := []*pb_almanac.LogEntry{}
	seen := map[string]struct{}{}
	for _, entry := range allEntries {
		if _, ok := seen[entry.Id]; ok {
			continue
		}
		seen[entry.Id] = struct{}{}

		result = append(result, entry)
		if len(result) >= int(request.Num) {
			break
		}
	}

	return &pb_almanac.SearchResponse{Entries: result}, nil
}

// searchChunk performs a search on a single chunk and pipes the result into
// the supplied channel.
func (m *Mixer) searchChunk(ctx context.Context, chunkId string, query string, num int32, resultChan chan *partialResult) {
	result := &partialResult{}
	chunk, err := m.storage.LoadChunk(chunkId)
	if err != nil {
		result.err = fmt.Errorf("unable to load chunk %s: %v\n", chunkId, err)
		resultChan <- result
		return
	}
	result.chunk = chunk

	entries, err := chunk.Search(ctx, query, num)
	if err != nil {
		result.err = fmt.Errorf("unable to perform search on chunk %s: %v\n", chunkId, err)
		resultChan <- result
		return
	}

	result.entries = entries
	resultChan <- result
}

// searchApender performs a search on a single appender and pipes the result into
// the supplied channel.
func (m *Mixer) searchAppender(ctx context.Context, appender pb_almanac.AppenderClient, request *pb_almanac.SearchRequest, resultChan chan *partialResult) {
	response, err := appender.Search(ctx, request)
	if err != nil {
		resultChan <- &partialResult{err: fmt.Errorf("unable to search appender: %v", err)}
		return
	}
	resultChan <- &partialResult{entries: response.Entries}
}

type partialResult struct {
	chunk   *storage.Chunk
	entries []*pb_almanac.LogEntry
	err     error
}

type byTimestamp []*pb_almanac.LogEntry

func (h byTimestamp) Len() int           { return len(h) }
func (h byTimestamp) Less(i, j int) bool { return h[i].TimestampMs < h[j].TimestampMs }
func (h byTimestamp) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }