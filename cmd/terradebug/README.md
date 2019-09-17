# Terradebug

Simple tool for simple debugging.

We try to accept both hex and base64 formats and provide a useful response.

Note we often encode bytes as hex in the logs, but as base64 in the JSON.

## Pubkeys

The following give the same result:

```
terradebug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
terradebug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
```

## Txs

Pass in a hex/base64 tx and get back the full JSON

```
terradebug tx <hex or base64 transaction>
```
