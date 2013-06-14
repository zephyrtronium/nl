package nl

import "io"

type nlreader struct {
	r io.Reader
	f func(*nlreader, []byte) (int, error)
	b [1]byte
	e error
}

func New(r io.Reader) io.Reader {
	return &nlreader{r, unknown, [1]byte{0}}
}

func (n *nlreader) Read(b []byte) (int, error) {
	return n.f(n.r, b)
}

func unknown(n *nlreader, b []byte) (int, error) {
	n, err := n.r.Read(b)
	if n <= 0 {
		return
	}
	for i, c := range b {
		if c == '\r' {
			if i == len(b)-1 {
				var n2 int
				if n2, n.e = n.r.Read(n.b[:]); n2 > 0 {
					if n.b[0] == '\n' {
						n.f = ignoreb
					} else {
						n.f = replaceb
						doRepl(b)
					}
				}
				return
			} else {
				if b[i+1] == '\n' {
					n.f = ignore
				} else {
					n.f = replace
					doRepl(b)
				}
			}
			return
		} else if c == '\n' {
			n.f = ignore
			return
		}
	}
	return
}

func ignoreb(n *nlreader, b []byte) (int, error) {
	b[0] = n.b[0]
	n.f = ignore
	n, err := n.r.Read(b[1:])
	if n.e != nil {
		err = n.e
	}
	return n + 1, err
}

func replaceb(n *nlreader, b []byte) (int, error) {
	b[0] = n.b[0]
	n.f = replace
	n, err := n.r.Read(b[1:])
	doRepl(b)
	if n.e != nil {
		err = n.e
	}
	return n + 1, err
}

func ignore(n *nlreader, b []byte) (int, error) {
	return n.r.Read(b)
}

func replace(n *nlreader, b []byte) (int, error) {
	n, err := n.r.Read(b)
	doRepl(b)
	return n, err
}

func doRepl(b []byte) {
	for i, c := range b {
		if c == '\r' {
			b[i] = '\n'
		}
	}
}
