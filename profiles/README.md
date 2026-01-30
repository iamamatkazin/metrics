File: agent
Type: inuse_space
Time: 2026-01-23 17:56:47 MSK
Duration: 50.01s, Total samples = 1040.21kB 
Showing nodes accounting for -1040.21kB, 100% of 1040.21kB total
      flat  flat%   sum%        cum   cum%
 -528.17kB 50.77% 50.77%  -528.17kB 50.77%  compress/flate.(*dictDecoder).init (inline)
 -512.05kB 49.23%   100%  -512.05kB 49.23%  net/http.(*persistConn).writeLoop
         0     0%   100%  -528.17kB 50.77%  compress/flate.NewReader
         0     0%   100%  -528.17kB 50.77%  compress/gzip.(*Reader).Reset
         0     0%   100%  -528.17kB 50.77%  compress/gzip.(*Reader).readHeader
         0     0%   100%  -528.17kB 50.77%  compress/gzip.NewReader (inline)
         0     0%   100%  -528.17kB 50.77%  github.com/iamamatkazin/metrics.git/internal/agent.(*Agent).Worker (inline)
         0     0%   100%  -528.17kB 50.77%  github.com/iamamatkazin/metrics.git/pkg/http.(*Client).Post
         0     0%   100%  -528.17kB 50.77%  io.Copy (inline)
         0     0%   100%  -528.17kB 50.77%  io.copyBuffer
         0     0%   100%  -528.17kB 50.77%  io.discard.ReadFrom
         0     0%   100%  -528.17kB 50.77%  main.main.func3
         0     0%   100%  -528.17kB 50.77%  main.main.func3.(*Agent).Worker.1
         0     0%   100%  -528.17kB 50.77%  net/http.(*cancelTimerBody).Read
         0     0%   100%  -528.17kB 50.77%  net/http.(*gzipReader).Read