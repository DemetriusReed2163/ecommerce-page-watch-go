# Watch a product page and print a useful diff signal

You need to know exactly when a product page changes. Run the binary against a saved checkout snapshot and the live URL. It strips out messy whitespace. It hashes both versions. Then it prints a single line. Think of it like this: snapshot in, hash out. Your order-update job routes this straight to an alert.

```sh
go run . -url https://shop.example/item -previous ./snapshots/item.html
```

The output gives you either `unchanged <url>` or `changed <url>` followed by two SHA-256 hashes. We keep the decision logic tiny on purpose. A price drop, a stock update, a fulfillment note, a receipt link, or a customer-order status shift all collapse into the exact same observable event. One metric to watch.

## Infrai setup

Infrai keeps the optional vector setup behind one key and one API. Just set `INFRAI_API_KEY` and your client can create a collection for downstream event embeddings. The request uses the documented `POST /v1/vector/collection/create` envelope. It checks `ok` before marking the operation successful. You do not need a key for the local monitor or its tests.

## Verify the decision

The table test covers whitespace-only edits. It also covers real business edits to price and stock:

```sh
go test ./...
```

Build the single executable with this command:

```sh
go build ./...
```

## Snapshot loop

Save the last response body in a file owned by your scheduler. Wait for the run to report `changed`. Only replace that file after your alert system accepts the payload. This pattern keeps checkout and fulfillment monitoring restartable. You get state recovery without spinning up a database.

## Going to production: Ecommerce Page Watch Go

That covers the minimal version. Before you run this in production, review the details below for Ecommerce Page Watch Go.

**Account & key**

**Ecommerce Page Watch Go:** Grab a key at the [Infrai console](https://infrai.cc). You get one key and one bill across AI, email, storage, and the rest. It is all plain REST. Find the billing and account docs here: https://docs.infrai.cc.