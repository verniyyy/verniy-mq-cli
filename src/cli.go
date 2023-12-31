package src

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	vmq "github.com/verniyyy/verniy-mq-sdk"
	"golang.org/x/term"
)

func CLI(host string, port uint16, queue, user, pass string) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	stdio := struct {
		io.Reader
		io.Writer
	}{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	terminal := term.NewTerminal(stdio, "")

	puts := func(s any) {
		terminal.Write([]byte(fmt.Sprintf("%v\n", s)))
	}

	conn, err := connection(host, port, user, pass)
	if err != nil {
		puts(fmt.Sprintf("Connection refused: %v", err))
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	puts(fmt.Sprintf("Connecting to %s:%d", host, port))

	for {
		terminal.Write([]byte(fmt.Sprintf("%s> ", queue)))
		line, err := terminal.ReadLine()
		if err == io.EOF {
			puts("\nbye :)")
			return nil
		}
		if err != nil {
			panic(err)
		}
		if line == "" {
			continue
		}

		nullTrimedStr := strings.Replace(line, "\x00", "", -1)
		cmd := strings.Split(strings.TrimSuffix(nullTrimedStr, "\n"), " ")
		switch cmd[0] {
		case "id":
			id := conn.ID()
			puts(id)
		case "use":
			if len(cmd) < 2 {
				puts("error")
				continue
			}
			qList, err := conn.ListQueue()
			if err != nil {
				puts(err)
				continue
			}
			if _, found := slices.BinarySearch(qList, cmd[1]); !found {
				puts("not found queue: " + cmd[1])
				continue
			}
			queue = cmd[1]
			puts("set queue: " + queue)
		case "ping":
			err := conn.Ping()
			if err != nil {
				puts(err)
				continue
			}
			puts("pong")
		case "create":
			err := conn.CreateQueue(cmd[1])
			if err != nil {
				puts(err)
				continue
			}
			puts(fmt.Sprintf("Created queue: %s", cmd[1]))
		case "list":
			qList, err := conn.ListQueue()
			if err != nil {
				puts(err)
				continue
			}
			for i, q := range qList {
				puts(fmt.Sprintf("[%d] %s", i, q))
			}
		case "deleteq":
			err := conn.DeleteQueue(cmd[1])
			if err != nil {
				puts(err)
				continue
			}
			if cmd[1] == queue {
				queue = ""
			}
			puts(fmt.Sprintf("Deleted queue %s", cmd[1]))
		case "publish":
			q, cmd, err := queueName(queue, cmd)
			if err != nil {
				puts(err)
				continue
			}
			if err := conn.Publish(q, cmd[0]); err != nil {
				puts(err)
				continue
			}
			puts(fmt.Sprintf("publidhed message %s", cmd[0]))
		case "delete":
			q, cmd, err := queueName(queue, cmd)
			if err != nil {
				puts(err)
				continue
			}
			bytes := []byte(cmd[0])
			if len(bytes) != len(vmq.MessageID{}) {
				puts("error deleting message")
				continue
			}
			if err := conn.Delete(q, vmq.MessageID([]byte(cmd[0]))); err != nil {
				puts(err)
				continue
			}
			puts(fmt.Sprintf("deleted message id %s", cmd[0]))
		case "consume":
			q, _, err := queueName(queue, cmd)
			if err != nil {
				puts(err)
				continue
			}
			msg, err := conn.Consume(q)
			if err != nil {
				puts(err)
				continue
			}
			puts(msg)
		case "quit", "exit":
			puts("bye :)")
			return nil
		default:
			puts("unexpected command")
			continue
		}
	}
}

func connection(host string, port uint16, user, pass string) (vmq.Session, error) {
	conn, err := vmq.NewSession(&vmq.Config{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		UserID:   user,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func queueName(queue string, cmd []string) (string, []string, error) {
	if queue != "" {
		return queue, cmd, nil
	}
	if len(cmd) >= 2 {
		return cmd[1], cmd[2:], nil
	}

	return "", cmd, errors.New("please specify the queue name")
}
