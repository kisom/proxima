Running the flight sim:

go test -v -bench '.*Mission.*' -cpuprofile cpu.out -memprofile mem.out -timeout 0 -run None
