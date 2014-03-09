var buf [1024]byte
addr, err := net.ResolveUDPAddr("udp", ":69")
sock, err := net.ListenUDP("udp", addr)
for {
    rlen, remote, err := sock.ReadFromUDP(buf)