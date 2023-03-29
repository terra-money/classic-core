<!--
order: 5
-->

# Ante

The fees module uses the ante handler to distribute fees between developers and the community.

## Handling

An [Ante Decorator](/x/feeshare/ante/ante.go) executes custom logic after each successful WasmExecuteMsg transaction. All fees paid by a user for transaction execution are sent to the `FeeCollector` module account during the `AnteHandler` execution before being redistributed to the registered contract developers.

If the `x/feeshare` module is disabled or the Wasm Execute Msg transaction targets an unregistered contract, the handler returns `nil`, without performing any actions. In this case, 100% of the transaction fees remain in the `FeeCollector` module, to be distributed elsewhere.

If the `x/feeshare` module is enabled and a Wasm Execute Msg transaction targets a registered contract, the handler sends a percentage of the transaction fees (paid by the user) to the withdraw address set for that contract.

1. The user submits an Execute transaction (`MsgExecuteContract`) to a smart contract and the transaction is executed successfully
2. Check if
   * fees module is enabled
   * the smart contract is registered to receive fee split
3. Calculate developer fees according to the `DeveloperShares` parameter.
4. Check what fees governance allows to be paid in
5. Check which contracts the user executed that also have been registered.
6. Calculate the total amount of fees to be paid to the developer(s). If multiple, split the 50% between all registered withdrawal addresses.
7. Distribute the remaining amount in the `FeeCollector` to validators according to the [SDK  Distribution Scheme](https://docs.cosmos.network/main/modules/distribution/03_begin_block.html#the-distribution-scheme).
