package main

// Connects a local bot to a remote server, which is a neat demo.
// The server has no awareness that a bot is controlling the player.

import (
	"flag"
	"log"

	"github.com/MustafaNova2/MatchPulse/pkg/backend"
	"github.com/MustafaNova2/MatchPulse/pkg/bot"
	"github.com/MustafaNova2/MatchPulse/pkg/client"
	"github.com/MustafaNova2/MatchPulse/pkg/frontend"
	"github.com/MustafaNova2/MatchPulse/proto"
	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", ":8888", "The server address.")
	flag.Parse()

	game := backend.NewGame()
	game.IsAuthoritative = false
	view := frontend.NewView(game)
	game.Start()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	grpcClient := proto.NewGameClient(conn)
	client := client.NewGameClient(game, view)

	bots := bot.NewBots(game)
	player := bots.AddBot("Bob")

	err = client.Connect(grpcClient, player.ID(), player.Name, "")
	if err != nil {
		log.Fatalf("connect request failed %v", err)
	}
	client.Start()

	view.Start()
	bots.Start()

	err = <-view.Done
	if err != nil {
		log.Fatal(err)
	}
}
