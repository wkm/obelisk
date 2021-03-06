
# resp

Generic implementation of Redis's RESP line-based command protocol over TCP, with the difference that out-of-order response is allowed within pipelined responses.

The protocol is constructed by attaching methods to a struct. Public methods are exposed. Methods must have simple parameters (ints, floats, and strings). Commands are derived from method names and are case insensitive. Error return values are automatically translated into `OK` when nil, and `ERROR err.Error()` otherwise.

For example:

    type protocol struct {}
    
    func (p* protocol) Get(key string) string {
        return p.values[key]
    }
    
    func (p *protocol) Put(key, val string) error {
        p.values[key] = val
        return nil
    }

The commands available are derived from the public methods of the backing interface, enabling queries of the form:

    PUT key1 123
    +OK
    PUT key2 345
    +OK
    GET key1
    $3
    123

