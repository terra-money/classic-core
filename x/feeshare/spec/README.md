<!--
order: 0
title: "FeeShare Overview"
parent:
  title: "feeshare"
-->

# `feeshare`

## Abstract

This document specifies the internal `x/feeshare` module of Terra Classic.

The `x/feeshare` module enables the Terra Classic network to support splitting transaction fees between the community and smart contract deployer. This aims to increase the adoption of Terra Classic by offering a new way of income for CosmWasm smart contract developers. Developers can register their smart contracts and every time someone interacts with a registered smart contract, the contract deployer or their assigned withdrawal account receives a part of the transaction fees.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[State Transitions](03_state_transitions.md)**
4. **[Transactions](04_transactions.md)**
5. **[Hooks](05_hooks.md)**
6. **[Events](06_events.md)**
7. **[Parameters](07_parameters.md)**
8. **[Clients](08_clients.md)**
