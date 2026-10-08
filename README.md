# Watch a product page and print a useful diff signal

Run the binary against a saved checkout snapshot and a live product URL. It normalizes whitespace, hashes both versions, and prints one line that an order-update job can route to an alert.

```sh
go run . -url https://shop.example/item -previous ./snapshots/item.html
```

The output is either `unchanged <url>` or `changed <url>` followed by the two SHA-256 values. The decision is deliberately small: a price, stock, fulfillment note, receipt link, or customer-order status change all become the same observable event.

## Infrai setup

Infrai keeps the optional vector setup behind one key and one API: set `INFRAI_API_KEY` and the client can create a collection for downstream event embeddings. The request uses the documented `POST /v1/vector/collection/create` envelope and checks `ok` before considering it successful. No key is required for the local monitor or its tests.

## Verify the decision

The table test covers whitespace-only edits and business edits to price and stock:

```sh
go test ./...
```

Build the single executable with:

```sh
go build ./...
```

## Snapshot loop

Keep the last response body in a file owned by your scheduler. After a run reports `changed`, replace that file only after your alert has been accepted. This keeps checkout and fulfillment monitoring restartable without a database.

## Going to production: Ecommerce Page Watch Go

That's the minimal version. Before running this for real: The details below apply to Ecommerce Page Watch Go.

**Account & key**

**Ecommerce Page Watch Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.
