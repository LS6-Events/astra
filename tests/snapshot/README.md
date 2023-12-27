# Snapshot testing
This directory contains snapshot tests for Astra. These tests are used to
ensure that the output of Astra is consistent and to prevent regressions. The tests are fairly slow right now, so the consideration is to run these in `-short` mode in the future but for now we can use them as appropriate.

## Running the tests
To run the tests, simply run `go test ./...` in this directory. The tests will run and compare the output of Astra to the expected output. If there are any differences, the test will fail and you will need to update the expected output. 

## Updating the snapshots
To do this, simply run `GENERATE_SNAPSHOTS=true go test ./...` and the tests will update the expected output.
 
