serverAddr, err := net.ResolveUDPAddr("udp", "192.168.1.1:69")
con, err := net.DialUDP("udp", nil, serverAddr)