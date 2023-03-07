package core

import (
	"github.com/PlainDAG/go-PlainDAG/config"
	"github.com/PlainDAG/go-PlainDAG/p2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func (n *Node) constructpubkeyMap() error {
	peeridSlice := n.network.H.Peerstore().Peers()

	stringpubkeymap := make(map[string]crypto.PubKey, len(peeridSlice))
	stringidmap := make(map[string]int, len(peeridSlice))
	for _, peerid := range peeridSlice {
		pubkey := n.network.H.Peerstore().PubKey(peerid)
		str := peerid.Pretty()
		pubkeybytes, err := crypto.MarshalPublicKey(pubkey)
		if err != nil {
			return err
		}
		stringpubkeymap[string(pubkeybytes)] = pubkey
		stringidmap[string(pubkeybytes)] = n.cfg.PubkeyIdMap[str]
	}
	n.cfg.StringpubkeyMap = stringpubkeymap
	n.cfg.StringIdMap = stringidmap
	return nil
}

func NewNode(filepath string) (*Node, error) {
	if plainMsgSize == 20 {
		messageconst = []byte{1, 2, 3, 5, 4, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	} else {

		messageconst = []byte{154, 80, 2, 82, 229, 242, 220, 255, 179, 87, 61, 154, 8, 88, 8, 107, 15, 130, 189, 156, 210, 66, 119, 158, 22, 164, 100, 166, 125, 222, 189, 140, 149, 138, 224, 105, 95, 112, 255, 126, 180, 47, 154, 161, 172, 224, 168, 68, 205, 1, 82, 65, 119, 220, 239, 199, 105, 36, 211, 130, 17, 17, 103, 221, 81, 251, 40, 174, 56, 32, 146, 64, 32, 181, 80, 209, 211, 86, 83, 153, 68, 131, 145, 200, 112, 162, 8, 165, 245, 11, 186, 213, 79, 2, 56, 69, 144, 62, 66, 63, 226, 226, 183, 23, 230, 176, 191, 150, 200, 66, 1, 221, 85, 140, 19, 251, 66, 183, 61, 235, 12, 47, 212, 153, 66, 125, 132, 214, 184, 218, 185, 125, 118, 61, 102, 15, 180, 44, 230, 134, 105, 13, 127, 44, 250, 1, 224, 47, 241, 108, 120, 95, 125, 49, 191, 125, 135, 222, 211, 120, 82, 31, 103, 199, 193, 217, 50, 34, 78, 214, 131, 99, 95, 18, 235, 235, 180, 40, 33, 188, 178, 39, 143, 147, 167, 96, 78, 150, 248, 165, 91, 10, 138, 102, 214, 206, 176, 200, 85, 185, 53, 121, 76, 116, 151, 119, 155, 76, 16, 211, 193, 184, 250, 202, 83, 91, 147, 87, 31, 234, 191, 1, 114, 192, 255, 105, 110, 14, 98, 57, 110, 87, 100, 154, 188, 38, 84, 87, 137, 200, 60, 51, 84, 216, 201, 82, 11, 170, 16, 52}
	}
	c := config.Loadconfig(filepath)
	n, err := p2p.Startpeer(c.Port, c.Prvkey, ReflectedTypesMap)
	if err != nil {
		return nil, err
	}
	chain, err := NewBlockchain()
	if err != nil {
		return nil, err
	}
	node := Node{
		bc:      chain,
		network: n,
	}
	c.Pubkey = n.H.Peerstore().PubKey(n.H.ID())
	Pubkeyraw, err := crypto.MarshalPublicKey(c.Pubkey)
	if err != nil {
		return nil, err
	}
	c.Pubkeyraw = Pubkeyraw

	node.cfg = c
	node.handler = NewStatichandler(&node)
	node.ls = NewLeaderSelector(&node)

	node.currentround.Store(0)
	return &node, err
}
