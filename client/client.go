package client

import (
	"bufio"
	"fmt"
	"github.com/willgnicks/eden/errors"
	"log"
	"net"
	"os"
	"sync"
)

// 地址和名字
type client struct {
	serverAddr string
	username   string
}

// new对象
func New(host string, port int, username string) *client {
	return &client{
		serverAddr: fmt.Sprintf("%s:%d", host, port),
		username:   username,
	}
}

// 开始运行
func (cli *client) Spin() {
	conn, err := net.Dial("tcp", cli.serverAddr)
	if err != nil {
		log.Fatalf("fail to connect server with err: %s\n", err.Error())
	}
	reader := negotiateProtocol(err, conn, cli)
	// 信号
	done := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(2) // 等两个goroutine
	go cli.read(reader, done, &wg)
	go cli.write(conn, done, &wg)
	// 注销用户
	wg.Wait() // 结束等待
	// conn.Write([]byte(fmt.Sprintf("%s\n", cli.protocol())))
	_ = conn.Close()
}

func negotiateProtocol(err error, conn net.Conn, cli *client) *bufio.Reader {
	_, err = conn.Write(cli.protocol())
	if err != nil {
		_ = conn.Close()
		log.Fatalf("send protocol failed with err: %s\n", err.Error())
	}
	reader := bufio.NewReader(conn)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		_ = conn.Close()
		log.Fatalf("read response from the server failed with err: %s\n", err.Error())
	}

	res := string(bytes)
	if res == string(errors.BytesProtocolErr) || res == string(errors.BytesUsernameExists) {
		_ = conn.Close()
		log.Fatalf(res)
	}
	fmt.Print(res)
	return reader
}

func (cli *client) protocol() []byte {
	return []byte(fmt.Sprintf("%s\n", cli.username))
}

func (cli *client) read(reader *bufio.Reader, done <-chan interface{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			wg.Done()
			break
		default:
			bytes, err := reader.ReadBytes('\n')
			if err != nil {
				continue
			}
			fmt.Print(string(bytes))
		}
	}
}

func (cli *client) write(conn net.Conn, done chan<- interface{}, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Printf("write text `%s` failed with err: %s\n", text, err.Error())
			continue
		}
		if text == "exit" {
			log.Println("closing the connection with the chat server...")
			done <- struct{}{}
			log.Println("exiting the program...")
			break
		}
	}
	wg.Done()
}
