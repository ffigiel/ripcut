#!/usr/bin/env bash

set -e

# 1. find app source
# pacmd list-sources | egrep '^\s+name:.*\.monitor'
# 2. map sink to source (remember to change output device)
# pacmd load-module module-null-sink sink_name=recording sink_properties=device.description=recording
# pacmd load-module module-combine-sink sink_name=combined sink_properties=device.description=combined \
#   slaves=recording,bluez_sink.FC_F1_52_79_2C_81.a2dp_sink

go build
parecord --device=recording.monitor --format=s16le --rate=44100 --file-format=raw --channels=2 | ./ripcut
