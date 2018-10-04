package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/ecadlabs/rosgw/utils"
	"gopkg.in/routeros.v2/proto"
)

func printSentence(w io.Writer, s *proto.Sentence) error {
	tw := tabwriter.NewWriter(w, 4, 0, 1, ' ', 0)

	if s.Tag != "" {
		fmt.Fprintf(tw, "%s @%s\n", s.Word, s.Tag)
	} else {
		fmt.Fprintf(tw, "%s\n", s.Word)
	}

	for _, p := range s.List {
		fmt.Fprintf(tw, "\t%s:\t%s\n", p.Key, p.Value)
	}

	fmt.Fprintf(tw, "\n")

	return tw.Flush()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] URL arguments\n", os.Args[0])
		flag.PrintDefaults()
	}

	username := flag.String("u", "", "User name")
	password := flag.String("p", "", "Password")
	listen := flag.Bool("l", false, "Listening mode")

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	opts := utils.DialOptions{
		URL:       args[0],
		Username:  *username,
		Password:  *password,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client, err := utils.Dial(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	if !*listen {
		reply, err := client.Run(args[1:]...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, s := range reply.Re {
			printSentence(os.Stdout, s)
		}

		printSentence(os.Stdout, reply.Done)
		os.Exit(0)
	}

	reply, err := client.Listen(args[1:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for s := range reply.Chan() {
		printSentence(os.Stdout, s)
	}

	if err := reply.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printSentence(os.Stdout, reply.Done)
}
