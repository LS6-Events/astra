# Snapshot testing
This directory contains snapshot tests for Astra. These tests are used to
ensure that the output of Astra is consistent and to prevent regressions. The tests are fairly slow right now, so the consideration is to run these in `-short` mode in the future but for now we can use them as appopriate.

## Running the tests
To run the tests, simply run `go test ./...` in this directory. The tests will run and compare the output of Astra to the expected output. If there are any differences, the test will fail and you will need to update the expected output. To do this, simply run `GENERATE_SNAPSHOTS=true go test ./...` and the tests will update the expected output.

## Adding new tests
To add new tests, copy the `0-template` directory, replace the name with the next available number (to track versioning) like `1-basic`, `2-doc-comments` etc. and modify the handler code accordingly. Be sure to double check the package imports. Feel free to modify/add/remove and handler endpoints to emphasise your testing needs. Once you have added your test, run the tests as described above to ensure that the output is as expected.

You shouldn't need to modify the `snapshot_test.go` file unless you want to change the handler paths etc.

**Be sure to create a new snapshot for your test whenever you create it!**