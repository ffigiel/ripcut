# ripcut
Save a PCM audio stream into silence-delimited sub-streams as FLAC.

Currently `ripcut` is hard-coded to work with a specific PCM format.
Most parameters can be adjusted by modifying the `config` variable in `main.go` and running `go build`.

Support for custom sample size / number of channels requires more work.

### Usage

```bash
git clone https://github.com/megapctr/ripcut
cd ripcut
go build

# using a file as input
./ripcut < path/to/your/stream.pcm

# using parecord - see `rec.sh` for details
parecord [options] | ./ripcut
```

Sub-streams will be saved as `stream_*.flac`
