package site

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/zeebo/bencode"
)

type AnnouncerStats struct {
	Status        string  `json:"status"`
	NumRequest    int     `json:"num_request"`
	NumSuccess    int     `json:"num_success"`
	NumError      int     `json:"num_error"`
	TimeRequest   float64 `json:"time_request"`
	TimeLastError float64 `json:"time_last_error"`
	TimeStatus    float64 `json:"time_status"`
	LastError     string  `json:"last_error"`
}

var headers = http.Header{
	"User-Agent":      []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"},
	"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"Accept-Charset":  []string{"ISO-8859-1,utf-8;q=0.7,*;q=0.3"},
	"Accept-Encoding": []string{"none"},
	"Accept-Language": []string{"en-US,en;q=0.8"},
	"Connection":      []string{"keep-alive"},
}

type requestParams struct {
	infoHash []byte
	peerID   string
	port     int
	numWant  int
}

type requestResponse struct {
	Failure     string `bencode:"failure"`
	Interval    int    `bencode:"interval"`
	MinInterval int    `bencode:"min_interval"`
	Complete    int    `bencode:"complete"`
	Incomplete  int    `bencode:"incomplete"`
	Peers       []byte `bencode:"peers"`
}

// https://www.bittorrent.org/beps/bep_0003.html
// announceToTracker announces to `tracker` and returns a list of parsed peers in the swarm
// TODO: add support for other protocols. For now it only supports HTTP trackers
func announceToTracker(ctx context.Context, tracker string, params requestParams) ([]string, error) {
	parsedURL, err := url.Parse(tracker)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("info_hash", string(params.infoHash))
	values.Set("peer_id", params.peerID)
	values.Set("port", strconv.Itoa(params.port))
	values.Set("uploaded", "0")
	values.Set("downloaded", "0")
	// Hack for tracker compatibility
	// https://github.com/HelloZeroNet/ZeroNet/issues/1248#issuecomment-364706403
	values.Set("left", "431102370")
	// TODO: compact is only advisory
	// clients should support bencoded peers dict as well
	values.Set("compact", "1")
	values.Set("numwant", strconv.Itoa(params.numWant))
	values.Set("event", "started")

	requestURL := fmt.Sprintf("%s?%s", parsedURL, values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsed, err := parseTrackerResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	return parsePeers(parsed.Peers)
}

func parseTrackerResponse(reader io.Reader) (*requestResponse, error) {
	parsed := new(requestResponse)
	return parsed, bencode.NewDecoder(reader).Decode(parsed)
}

// https://www.bittorrent.org/beps/bep_0023.html
func parsePeers(peerList []byte) ([]string, error) {
	if len(peerList)%6 != 0 {
		return nil, errors.New("unexpected peer list size")
	}

	var ips []string
	for i := 0; i < len(peerList); i += 6 {
		ip, port := peerList[i:i+4], peerList[i+4:i+6]

		peerID := net.IPv4(ip[0], ip[1], ip[2], ip[3])
		peerPort := binary.BigEndian.Uint16(port)

		ips = append(ips, fmt.Sprintf("%s:%d", peerID, peerPort))
	}

	return ips, nil
}

// Announce announces to trackers the new peer
// TODO: where to get the peer ID?
// TODO: debounce this function
func (s *Site) Announce() {
	h := sha1.New()
	io.WriteString(h, s.addr)
	params := requestParams{
		infoHash: h.Sum(nil),
		peerID:   random.PeerID(),
		port:     config.FileServerPort,
		numWant:  10,
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	defer cancel()

	now := func() float64 {
		// Convert to a representation similar to Python's time.time function:
		// "Return the current time in seconds since the Epoch."
		// "Fractions of a second may be present if the system clock provides them."
		return float64(time.Now().UnixNano()) / 1e9
	}

	var wg sync.WaitGroup
	wg.Add(len(config.Trackers))
	var addedPeers int64
	for _, tracker := range config.Trackers {
		tracker := tracker
		go func() {
			defer wg.Done()

			s.trackersMutex.RLock()
			stats, exists := s.trackers[tracker]
			s.trackersMutex.RUnlock()
			if !exists {
				stats = new(AnnouncerStats)
			}
			stats.NumRequest++
			stats.TimeRequest = now()
			peers, err := announceToTracker(ctx, tracker, params)
			if err != nil {
				stats.Status = "error"
				stats.LastError = err.Error()
				stats.TimeLastError = now()
				stats.NumError++
				s.trackersMutex.Lock()
				s.trackers[tracker] = stats
				s.trackersMutex.Unlock()
				return
			}
			stats.NumSuccess++
			stats.Status = "announced"
			stats.TimeStatus = now()
			s.trackersMutex.Lock()
			s.trackers[tracker] = stats
			s.trackersMutex.Unlock()

			s.peersMutex.Lock()
			for _, peer := range peers {
				s.peers[peer] = struct{}{}
			}
			s.peersMutex.Unlock()

			atomic.AddInt64(&addedPeers, int64(len(peers)))
		}()
	}
	wg.Wait()

	_ = s.broadcastSiteChange("peers_added", addedPeers)

	// TODO: if trackers.json file exists annouce using the trackers defined there
	// If it doesn't exist use bootstrap trackers in config.Trackers
	// In either case, update trackers.json and return the updates stats here
}

func (s *Site) AnnouncerStats() map[string]*AnnouncerStats {
	return s.trackers
}
