package conn

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ecadlabs/rosgw/conn/pool"
	"golang.org/x/crypto/ssh"
)

type KeyFunc func() ([]byte, error)

type Config struct {
	KeyFunc  KeyFunc
	Username string
	Password string
}

type Client struct {
	*ssh.Client
	conn net.Conn
}

func (c *Client) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func dial(ctx context.Context, address string, c *Config) (client *Client, err error) {
	var auth []ssh.AuthMethod

	if c.KeyFunc != nil {
		key, err := c.KeyFunc()
		if err != nil {
			return nil, fmt.Errorf("key func: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}

		auth = append(auth, ssh.PublicKeys(signer))
	}

	if c.Password != "" {
		auth = append(auth, ssh.Password(c.Password))
	}

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	config := ssh.ClientConfig{
		User:            c.Username,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var (
		sshConn ssh.Conn
		chans   <-chan ssh.NewChannel
		reqs    <-chan *ssh.Request
	)

	if d, ok := ctx.Deadline(); ok {
		conn.SetDeadline(d)
	}

	ch := make(chan struct{})

	go func() {
		sshConn, chans, reqs, err = ssh.NewClientConn(conn, address, &config)
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case <-ch:
		if err != nil {
			return nil, err
		}
	}

	conn.SetDeadline(time.Time{})

	client = &Client{
		Client: ssh.NewClient(sshConn, chans, reqs),
		conn:   conn,
	}

	return client, nil
}

// Connection pool wrapper below

type poolItem struct {
	client *Client
	conf   *Config
	addr   string
}

func poolKey(username, addr string) string { return username + "@" + addr }

func (p *poolItem) Key() interface{} { return poolKey(p.conf.Username, p.addr) }

func (p *poolItem) Finalize() {
	p.client.Close()
}

type newPoolItemArg struct {
	ctx     context.Context
	address string
	conf    *Config
}

func newPoolItem(arg interface{}) (pool.Item, error) {
	a := arg.(*newPoolItemArg)

	client, err := dial(a.ctx, a.address, a.conf)
	if err != nil {
		return nil, err
	}

	return &poolItem{
		client: client,
		conf:   a.conf,
		addr:   a.address,
	}, nil
}

type Pool pool.Pool

func NewConnPool(maxconn int) *Pool {
	return (*Pool)(&pool.Pool{
		MaxLength: maxconn,
		New:       newPoolItem,
	})
}

func (p *Pool) Dial(ctx context.Context, address string, c *Config) (client *Client, err error) {
	v, err := (*pool.Pool)(p).Get(poolKey(c.Username, address), &newPoolItemArg{
		ctx:     ctx,
		address: address,
		conf:    c,
	})

	if err != nil {
		return nil, err
	}

	return v.(*poolItem).client, nil
}
