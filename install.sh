#!/bin/bash
esbuild clientSideRouting/src.js --bundle --outfile="clientSideRouting/out.js"
go build
sudo mv ./melte-framework /usr/bin/melte