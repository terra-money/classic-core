# I. Why we need to setup an automated testing framework?
In many cases, we have observed that unit test isn't enough to guarantee that the whole system will not break. We simply got lucky when Ed found a problem through his manual system testing that may not always happen.

In some cases, not even the unit test is correct due to writer misunderstanding of required logic.

# II. Framework
1. Standarize writing unit test (we can even go further with a ci to check). A PR should include:
* feature specification (in written markdown file in a folder called feature-request) (Gherkin language)
* unit test write

2. a full - fledged e2e test flow through ci that does following flow
* upgrade testing:
    * for chain (if exists, else use current terrad): do an upgrade gov test to newer chain binary
    * for wasm (if there are changes to wasm, else use current wasm files): change to newer wasm files
* regression testing (change to system should not break existing one)
    * backend: 
        * chain should run normal 
        * ibc should work normal
        * wasm should run normal
        * state sync should run normal
    * frontend: api testing of most important transaction types and query
* fuzz testing:
    * some nice background reading: https://github.com/osmosis-labs/osmosis/blob/main/simulation/ADR.md
    * we need to care only about events that invoke a state change, we then create randomized request to these events:
        * transaction: create randomized request and broadcast multiple times
        * begin block: randomized environment condition that invokes a begin block logic (Ex: system time, number of validators, ...)
        * end block: randomized environment condition that invokes an end block logic (Ex: system time, number of validators, ...)
        * init genesis: create randomized genesis state

It is hard to invoke a begin, end block logic without first understanding what environment condition it relies on.