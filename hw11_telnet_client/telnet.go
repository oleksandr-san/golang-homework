package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrNotConnected = errors.New("client is not connected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return ErrNotConnected
	}

	if err := t.conn.Close(); err != nil {
		return err
	}
	return t.in.Close()
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return ErrNotConnected
	}

	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return ErrNotConnected
	}

	_, err := io.Copy(t.conn, t.in)
	return err
}
