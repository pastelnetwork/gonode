# Register Volumes

## Context

Currently, we can register files as cascade tickets exceeding 500 MB in size. However, this capability has its limits. As we attempt to register larger files, we encounter significant roadblocks, including extended registration times and overall system slowdowns.

The opportunity here is to develop a scalable solution that allows us to store arbitrarily large files efficiently.

## Solution

The most feasible solution, which requires minimal changes to the cnode, involves breaking down large files into volumes, each approximately 350 MB in size. We can register these volumes as separate cascade files. We need to encode in the tickets that a particular ticket is one of the volumes of a larger file, enabling the downloading agent (walletnode) to identify, compile, and return the complete file.

For instance, the action ticket structure Currently looks like this:

```
{
  "height": 48062,
  "ticket": {
    "action_ticket": {
      "action_ticket_version": 1,
      "caller": "jXY6jFwJPhqu5nAh63Wie52SSUtvBAkrg8LWD6GCMra7rcbQrLvo4L3jRoutSzEfudsSKUGyNULP7KMn57nPwB",
      "blocknum": 48060,
      "block_hash": "001a1f01cc4f796eb449e74b9718e4d31d0702af6259572bfc3a249a31253665",
      "action_type": "cascade",
      "api_ticket": {
        "data_hash": "BsWta7mTVAGC7gMJMkONTbsSRUaRgpufMP03yq3RtOw=",
        "file_name": "linux-modules-6.6.1-060601-generic_6.6.1-060601.202311151749_amd64.deb",
        "rq_ic": 3067274846,
        "rq_max": 50,
        "rq_ids": [
          "AhmhFAc3W6Fbv7qP3huATvT2K3LjMsTD81YtuAJWhEdX",
          "8HtPp8CEWkVbMb71Pe7FqhNKDxpvQvqcUrik55jiqtqx",
          ...
        ],
        "rq_oti": "AAo6IMAAw1ABABEI",
        "original_file_size_in_bytes": 171581632,
        "file_type": "application/vnd.debian.binary-package",
        "make_publicly_accessible": false
      },
      "called_at": 48060,
      "key": "zxyizxsfyrvwmemtyhvzm6ojtfcybwwffbcwtawp33dnwxwdsjua====",
      "label": "9dc06426e648fb5485b7adff232dc6e42bab69dc653f20581066a978c9e44588",
      "signatures": {
          ...
      },
      "storage_fee": 821000,
      "type": "action-reg",
      "version": 1
    }
  },
  "tx_info": {
    "is_compressed": false,
    "multisig_outputs_count": 0,
    "multisig_tx_total_fee": 0,
    "uncompressed_size": 6160
  },
  "txid": "933127993338f728543ae9bc43fbbb93c259dbb4361f40d6a6008fb9a04fe216"
}
```

As we are aware, the 'api_ticket' is only parsed by the gonode, so no changes are required on the cnode if we include the keys of all volumes so that the walletnode can seamlessly register all volumes as separate files and later on, compile and regenerate the original file whenever requested.

### Registration

Ideally, we want to register all volumes in parallel because serial registration would be too time-consuming, defeating the purpose of this solution.

#### Faster Registration

To register all tickets in parallel while ensuring that volumes link to each other, we cannot use the txid of the registration ticket, as it is only available after successful registration. Using the txid would force serial registration, which is undesirable.

Instead, we can use the 'key' field inside the ticket, which gonode assigns a random UUID at the moment. What we can do is pre-generating a set of UUIDs for all volumes and linking them together inside the api_ticket.

```
api_ticket": {
  "data_hash": "BsWta7mTVAGC7gMJMkONTbsSRUaRgpufMP03yq3RtOw=",
  "file_name": "linux-modules-6.6.1-060601-generic_6.6.1-060601.202311151749_amd64.deb",
  "rq_ic": 3067274846,
  "rq_max": 50,
  "rq_ids": [
    "AhmhFAc3W6Fbv7qP3huATvT2K3LjMsTD81YtuAJWhEdX",
    "8HtPp8CEWkVbMb71Pe7FqhNKDxpvQvqcUrik55jiqtqx",
    ...
  ],
  "rq_oti": "AAo6IMAAw1ABABEI",
  "original_file_size_in_bytes": 171581632,
  "file_type": "application/vnd.debian.binary-package",
  "make_publicly_accessible": false,
  "volumes": [
    "uuid_001",
    "uuid_002",
    "uuid_003",
    "uuid_004",
  ]
}
```

This approach allows simultaneous registration of all volumes, significantly reducing the registration time.

#### The Pre-Burn TxID Challenge

Another challenge is the Pre-Burn TxID, which serves as a unique identifier and must be valid, as cnode verifies both elements.

The solution is to modify the upload file endpoint. Currently, it looks like this:

```
http://localhost:8080/openapi/cascade/upload

{
  "file_id": "KcEoczz3",
  "expires_in": "2024-06-04T01:13:11Z",
  "total_estimated_fee": 70,
  "required_preburn_amount": 12
}
```

We need to change it to return a list of pre-burn transactions that the user must perform. It would look like this (preferably writing a new version of the endpoint, i.e., v2):

```
http://localhost:8080/openapi/v2/cascade/upload

[
  {
    "file_id": "KcEoczz3",
    "expires_in": "2024-06-04T01:13:11Z",
    "total_estimated_fee": 710,
    "required_preburn_txn_amounts": [
      24,
      24,
      24,
      24,
      24,
      10,
    ]
  }
]
```

Consequently, the start register endpoint would also change to accept a list of Pre-Burn TxIDs instead of a single TxID:

```
http://localhost:8080/openapi/cascade/v2/start/KcEoczz3

{
  "burn_txids": [
    "ac4af899e7ac9f06edb1ab8f2d2c951ef7934821a4fbc3592b55330b836d7415",
    "bc5xf899e7ac9f06edb1ab8f2d2c891ef7934821a4fbc3592b55330b836d9073",
    "xd8xfh76e7ac9f06edb1ab8f2d2c311ef7934821a4fbc3592b55330b836d1921",
    "hyi5xf766e7ac9f06edb1ab8f2p2c951ef7934821a4flk3592b55330b836io90",
    "po9cf861e7ac9f06edb1ab8f2b2c951ef7934821a4fad3592b55330b836u6755",
    "cc0df128e7ac9f06edb1ab8f0a2c951ef7934821a4fnh3592b55330b836e0098",
  ],
  "app_pastelid": "jXai5nZSXLQr7KHQ5H7Q94jRN8NnZzWynKR8fgB9jaJ3YuqNSzD7JdA9JGVPeiS9C9tUsUcftgWFvY2CAeVTNa",
  "make_publicly_accessible": true
}
```

This change enables walletnode to treat each volume as a separate ticket, each with its own Burn TxID.

Additionally, upon calling the upload file API, walletnode should break down the large file into volumes and maintain an index table mapping file IDs to all their volumes (filenames). This allows the walletnode to correctly index all volumes during registration and track registration and activation status. Using a local database like SQLite ensures persistence across shutdowns or restarts and helps handle scenarios where some cascade tickets fail due to insufficient balance.

Maintaining this information on disk rather than in-memory allows the system to recover more easily and provide accurate progress updates to users via a new v2 status endpoint.

### Detailed Process

Here is how the process

 would look:

1. Walletnode receives a registration request for a larger file via the upload file endpoint.
2. Walletnode breaks down the file into chunks/volumes and returns a list of amounts for the Pre-Burn Transactions that the caller needs to perform. Internally, it maintains an index linking the file_id to all its volumes (filenames and required preburn transaction amounts).
3. The registration endpoint then accepts a list of Pre-Burn Transactions.
4. Walletnode validates that the user performed the required number of Pre-Burn Transactions with the correct amounts.
5. Upon successful verification, walletnode starts registering all volumes:
   a. Pre-generate UUID combinations for all chunks and store them in the database records for future reference.

   b. Initiate the registration process for all chunks, registering no more than five per block. For more than five volumes, initiate registration in subsequent blocks, updating the database to mark the start of registration.

   c. Once all tickets are successfully registered, update the database to mark successful registration and then create activation tickets, again updating the database to mark successful activation for each volume.

   d. If any registration ticket fails, retry registering the ticket three times in different blocks, maintaining a counter in the database and storing encountered errors.

   e. Similarly, if any activation ticket creation fails, update the database accordingly.

The purpose of maintaining everything on disk/database is to enable recovery if some volumes fail to get either registered or activated.

6. Once all volumes are registered and activated, create a `cascade_multi_volume_metadata_ticket` that includes the hash of the original file, TxIDs for all registration and activation tickets for all volumes, and the keys.
7. Walletnode can then return either the TxID of the `cascade_multi_volume_metadata_ticket` or the TxID of the first volume.

### Download

Upon receiving a download request, walletnode can do one of the following:

### 1- Using the New Ticket Type

 Check if the TxID is of a `cascade_multi_volume_metadata_ticket`. If so, it fetches all volumes from the network, compiles them in the correct order, generates the hash, verifies it against the ticket, and returns the file.

### 2- Without using the New Ticket Type

 Retrieve the action ticket, decode to see if it has volumes, fetch all volumes from the network, compile them in the correct order, generate the hash, verify it against the ticket, and return the file.

This comprehensive approach ensures efficient handling of large files, maintains system performance, and provides reliable recovery mechanisms in case of failures.