package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ecadlabs/rosgw/conn"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type RemoteCommand struct {
	Pool   *conn.Pool
	Logger *logrus.Logger
}

type sshCommandResponse struct {
	io.Reader
	session *ssh.Session
	release func()
}

type readNotifier struct {
	rd         io.Reader
	done       chan<- struct{}
	pendingErr <-chan error
	closed     bool
	err        error
	m          sync.Mutex
}

func (r *readNotifier) close(err error) {
	r.m.Lock()
	defer r.m.Unlock()

	if !r.closed {
		r.err = err
		r.closed = true
		close(r.done)
	}
}

func (r *readNotifier) error() error {
	r.m.Lock()
	defer r.m.Unlock()

	return r.err
}

func (r *readNotifier) Read(p []byte) (int, error) {
	r.m.Lock()
	if r.closed {
		defer r.m.Unlock()
		return 0, r.err
	}
	r.m.Unlock()

	n, err := r.rd.Read(p)
	if err != nil {
		select {
		case err = <-r.pendingErr:
		default:
		}

		r.close(err)
	}
	return n, err
}

func (s *sshCommandResponse) Close() (err error) {
	rn := s.Reader.(*readNotifier)
	defer rn.close(io.EOF)

	if err := s.session.Wait(); err != nil {
		return err
	}

	if err := s.session.Close(); err != nil && err != io.EOF {
		return err
	}

	// Reuse live connection
	if e := rn.error(); e == nil || e == io.EOF {
		s.release()
	}

	return nil
}

func (c *RemoteCommand) Run(ctx context.Context, address string, conf *conn.Config, cmd string) (response io.ReadCloser, err error) {
	l := c.Logger.WithField("address", address)

	var (
		session  *ssh.Session
		client   *conn.Client
		readDone chan struct{}
	)

	pendingErr := make(chan error, 1)

	for {
		cl := c.Pool.Get(address, conf.Username)

		var dialed bool
		if cl == nil {
			l.Info("establishing SSH connection...")

			cl, err = conn.Dial(ctx, address, conf)
			if err != nil {
				return nil, err
			}

			dialed = true
		}

		if d, ok := ctx.Deadline(); ok {
			cl.SetDeadline(d)
		}

		rd := make(chan struct{})

		go func() {
			select {
			case <-ctx.Done():
				pendingErr <- ctx.Err()
				cl.SetDeadline(time.Now())
			case <-rd:
			}
		}()

		session, err = cl.NewSession()
		if err == nil {
			client = cl
			readDone = rd
			break
		}

		close(rd)
		cl.Close()

		if e, ok := err.(net.Error); ok && e.Timeout() || err == context.Canceled || dialed {
			return nil, fmt.Errorf("new session: %v", err)
		}

		l.Warnln(err)
		// Try next
	}

	defer func() {
		if err != nil {
			client.Close()
		}
	}()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	l.Infof("issuing `%s' command...", cmd)

	if err = session.Start(cmd); err != nil {
		return nil, fmt.Errorf("session start: %v", err)
	}

	rn := readNotifier{
		rd:         bufio.NewReader(stdout),
		done:       readDone,
		pendingErr: pendingErr,
	}

	res := sshCommandResponse{
		Reader:  &rn,
		session: session,
		release: func() {
			client.SetDeadline(time.Time{})
			c.Pool.Put(client)
		},
	}

	return &res, nil
}
