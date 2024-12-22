#!/bin/bash

# write #1
echo ""
echo "write #1.1"
curl http://localhost:3001/write
curl http://localhost:3001/read
curl http://localhost:3002/read

# write #2
echo ""
echo "write #2.1"
curl http://localhost:3002/write
curl http://localhost:3001/read
curl http://localhost:3002/read

# write #2
echo ""
echo "write #2.2"
curl http://localhost:3002/write
curl http://localhost:3001/read
curl http://localhost:3002/read

# write #1
echo ""
echo "write #1.2"
curl http://localhost:3001/write
curl http://localhost:3001/read
curl http://localhost:3002/read