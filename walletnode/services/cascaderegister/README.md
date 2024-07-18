# API Documentation

## 1. Upload Large Files

For cascade registration, files bigger than 250 MB will not be supported by the upload endpoint. For these files, use the walletnode’s new `/v2/upload` endpoint.

### Endpoint
POST `/openapi/cascade/v2/upload`

### Request
The request structure is the same as the `/v1/upload` endpoint.

### Response
The response structure has changed as detailed below.

#### Sample Response
```json
{
    "file_id": "NwWDvdIg",
    "total_estimated_fee": 6290,
    "required_preburn_transaction_amounts": [
        204,
        204,
        204,
        34
    ]
}

```
The client would need to perform PreBurn Transactions of these amounts i.e., 3 transactions of amount 204 and then 1 transaction of amount 34 to call the next endpoint.

### 2. Start Task
Once preburn transactions are done, the next step is to call the start task endpoint which has been slightly modified.

#### Endpoint

POST `/openapi/cascade/start/<<file_id>>`
#### Sample Request Payloads
You can either pass a single burn_txid or an array of burn_txids as shown:

```
{
  "burn_txids": [
    "a4a52e3c67e2d3825d046f6e75145f104c2e3b11c320bee58ac3032b4862d5e4",
    "0786b7388d557b63ad2cb859288dd69012f858b5c8b58fd83bb56a00ffe3c3d3",
    "47917dd82864aa92949ddff83dee71c61795a13f903b25f503737ade708ddc86",
    "b60283f595a813f90acf15179249c61b5b1a4c480eac8624a0a42a11b56d06f4"
  ],
  "app_pastelid": "jXY6jFwJPhqu5nAh63Wie52SSUtvBAkrg8LWD6GCMra7rcbQrLvo4L3jRoutSzEfudsSKUGyNULP7KMn57nPwB",
  "make_publicly_accessible": true
}

```

OR

```
{
  "burn_txid": "a4a52e3c67e2d3825d046f6e75145f104c2e3b11c320bee58ac3032b4862d5e4",
  "app_pastelid": "jXY6jFwJPhqu5nAh63Wie52SSUtvBAkrg8LWD6GCMra7rcbQrLvo4L3jRoutSzEfudsSKUGyNULP7KMn57nPwB",
  "make_publicly_accessible": true
}

```
In the case of multi-volume, this endpoint would return a string which is comma-separated for Task IDs.

### 3. Check Multi-Volume Registration Status
Next, we can check the status of the multi-volume registration task using file_id.

#### Endpoint:
GET `/openapi/cascade/registration_details/<file_id>`



#### Sample Response
Sample response of this endpoint when the registration is in progress (Example: 7 volumes of the file, hence file_index ranges from 0-6):
```


{
    "files": [
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.001",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "0",
            "base_file_id": "NwWDvdIg",
            "task_id": "HskCThhY",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "a4a52e3c67e2d3825d046f6e75145f104c2e3b11c320bee58ac3032b4862d5e4",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "ddb4c8e7-39d8-45e0-afc7-4a4753b01af8",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 27,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.001",
                    "reg_started_at": "2024-07-16 16:08:03.363832595 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.002",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "1",
            "base_file_id": "NwWDvdIg",
            "task_id": "6cRf8Zff",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "0786b7388d557b63ad2cb859288dd69012f858b5c8b58fd83bb56a00ffe3c3d3",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "89dc2c94-4623-4010-a58d-b2d70631bbb4",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 28,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.002",
                    "reg_started_at": "2024-07-16 16:08:03.41780498 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.003",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "2",
            "base_file_id": "NwWDvdIg",
            "task_id": "szcAqL9V",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "47917dd82864aa92949ddff83dee71c61795a13f903b25f503737ade708ddc86",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "dbb98261-75a5-4f89-9b12-d1d9de157d7d",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 29,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.003",
                    "reg_started_at": "2024-07-16 16:08:03.444456683 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.004",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "3",
            "base_file_id": "NwWDvdIg",
            "task_id": "XkD8T37l",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "b60283f595a813f90acf15179249c61b5b1a4c480eac8624a0a42a11b56d06f4",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "75638e9b-84da-45ba-88b7-9aca2e2ee9b0",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 30,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.004",
                    "reg_started_at": "2024-07-16 16:08:03.474462947 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.005",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "4",
            "base_file_id": "NwWDvdIg",
            "task_id": "z26N4I05",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "ae07f7f55312a4e74e8011249603bd52b7c46f60a69a2ca4215dc4f1763acd36",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "105fb7ac-533f-44db-a68b-6143ae635d42",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 31,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.005",
                    "reg_started_at": "2024-07-16 16:08:03.503680527 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.006",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "5",
            "base_file_id": "NwWDvdIg",
            "task_id": "7qAn2Xf7",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "2dea791cb8a57148e1d702fe68affeb9ba929b28e96806883c121c85903c4469",
            "req_amount": 1020,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "08a8b64e-4015-4156-9c65-e0acbf718025",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 32,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.006",
                    "reg_started_at": "2024-07-16 16:08:03.534410197 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.007",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "6",
            "base_file_id": "NwWDvdIg",
            "task_id": "CXbK5zVw",
            "reg_txid": "",
            "activation_txid": "",
            "req_burn_txn_amount": 34,
            "burn_txn_id": "6e33a3e055203c0d0fd7847aedffac50cf45622927700dac5f843cfe655977ec",
            "req_amount": 170,
            "is_concluded": false,
            "cascade_metadata_ticket_id": "",
            "uuid_key": "40aeb7fd-165e-42eb-a127-a7c8628764a1",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 0,
            "registration_attempts": [
                {
                    "id": 33,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.007",
                    "reg_started_at": "2024-07-16 16:08:03.568282649 +0000 UTC",
                    "processor_sns": "",
                    "finished_at": "0001-01-01 00:00:00 +0000 UTC",
                    "is_successful": false,
                    "error_message": ""
                }
            ],
            "activation_attempts": []
        }
    ]
}

```

and here’s how it looks like when registration is successful and multi-volume contract ticket has also been created

```
{
    "files": [
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.001",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "0",
            "base_file_id": "NwWDvdIg",
            "task_id": "HskCThhY",
            "reg_txid": "7195a4b67d82f242f7f3979407b072ac0ca7e8bca8c157ad0cfdcbc34176f358",
            "activation_txid": "e448cb9fbd78f580d0feaeaee34bba8937406fdff9cb534165978bf1cebcc47e",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "a4a52e3c67e2d3825d046f6e75145f104c2e3b11c320bee58ac3032b4862d5e4",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "ddb4c8e7-39d8-45e0-afc7-4a4753b01af8",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97050,
            "registration_attempts": [
                {
                    "id": 27,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.001",
                    "reg_started_at": "2024-07-16 16:08:03.363832595 +0000 UTC",
                    "processor_sns": "154.12.243.33:24444,154.12.246.128:24444,154.12.246.131:24444,",
                    "finished_at": "2024-07-16 16:13:34.896512035 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 25,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.001",
                    "activation_attempt_at": "2024-07-16 16:45:03.09246019 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.002",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "1",
            "base_file_id": "NwWDvdIg",
            "task_id": "6cRf8Zff",
            "reg_txid": "1a567609281db81f5f119d6c2b5e6e8df2976cc1b450ddd1d239712534084d01",
            "activation_txid": "eb986726b716fbdd09ed399347d82aa0397c50ff2d92f57abf255dadd59676ec",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "0786b7388d557b63ad2cb859288dd69012f858b5c8b58fd83bb56a00ffe3c3d3",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "89dc2c94-4623-4010-a58d-b2d70631bbb4",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97048,
            "registration_attempts": [
                {
                    "id": 28,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.002",
                    "reg_started_at": "2024-07-16 16:08:03.41780498 +0000 UTC",
                    "processor_sns": "154.38.164.99:24444,154.12.246.129:24444,154.12.246.131:24444,",
                    "finished_at": "2024-07-16 16:10:19.899360674 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 23,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.002",
                    "activation_attempt_at": "2024-07-16 16:36:45.992604839 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.003",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "2",
            "base_file_id": "NwWDvdIg",
            "task_id": "szcAqL9V",
            "reg_txid": "4afb36a8d91d619341dfb932ff5ca0079c8bddb5e3bd35f1f3ac365855a249c8",
            "activation_txid": "66da2108c5c3805eb7ca0b96ab197da8a09494b748110defc7deaacc984d6566",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "47917dd82864aa92949ddff83dee71c61795a13f903b25f503737ade708ddc86",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "dbb98261-75a5-4f89-9b12-d1d9de157d7d",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97050,
            "registration_attempts": [
                {
                    "id": 29,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.003",
                    "reg_started_at": "2024-07-16 16:08:03.444456683 +0000 UTC",
                    "processor_sns": "207.244.236.251:24444,154.12.246.130:24444,154.38.164.104:24444,",
                    "finished_at": "2024-07-16 16:11:03.871768764 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 27,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.003",
                    "activation_attempt_at": "2024-07-16 16:45:04.003916154 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.004",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "3",
            "base_file_id": "NwWDvdIg",
            "task_id": "XkD8T37l",
            "reg_txid": "a777884586342015a47fc9d223fa9c3ee866990acd770795eba56ea32a980a2e",
            "activation_txid": "bff4a79c5d671cce7ecbc22404d6ecbc850e6bfa4b69a55fbf6d342ae8a03ee2",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "b60283f595a813f90acf15179249c61b5b1a4c480eac8624a0a42a11b56d06f4",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "75638e9b-84da-45ba-88b7-9aca2e2ee9b0",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97050,
            "registration_attempts": [
                {
                    "id": 30,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.004",
                    "reg_started_at": "2024-07-16 16:08:03.474462947 +0000 UTC",
                    "processor_sns": "154.12.246.128:24444,154.12.246.130:24444,207.244.236.251:24444,",
                    "finished_at": "2024-07-16 16:11:00.604286877 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 26,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.004",
                    "activation_attempt_at": "2024-07-16 16:45:03.157739528 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.005",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "4",
            "base_file_id": "NwWDvdIg",
            "task_id": "z26N4I05",
            "reg_txid": "047fa18ed022e8a36dcb4b08b1cf16c167aa3e47c21aa20369111b602520cc0b",
            "activation_txid": "109f73144b1e665b6bf2ca4a2ce2e9fd0dd9d59821e421c0b5c9cf2d77912c42",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "ae07f7f55312a4e74e8011249603bd52b7c46f60a69a2ca4215dc4f1763acd36",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "105fb7ac-533f-44db-a68b-6143ae635d42",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97050,
            "registration_attempts": [
                {
                    "id": 31,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.005",
                    "reg_started_at": "2024-07-16 16:08:03.503680527 +0000 UTC",
                    "processor_sns": "154.12.246.131:24444,154.38.164.99:24444,154.12.246.129:24444,",
                    "finished_at": "2024-07-16 16:10:44.896195133 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 24,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.005",
                    "activation_attempt_at": "2024-07-16 16:45:01.521751888 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.006",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "5",
            "base_file_id": "NwWDvdIg",
            "task_id": "7qAn2Xf7",
            "reg_txid": "1d87d071cc83b7917b371e8f040e09adadbd07d92fccc85df0e4115369b0e988",
            "activation_txid": "53abacb80b32f88dfee8cec55e1b9372f0cc0708f2d667e4dbfa69f3e1cb1bbe",
            "req_burn_txn_amount": 204,
            "burn_txn_id": "2dea791cb8a57148e1d702fe68affeb9ba929b28e96806883c121c85903c4469",
            "req_amount": 1020,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "08a8b64e-4015-4156-9c65-e0acbf718025",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97051,
            "registration_attempts": [
                {
                    "id": 32,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.006",
                    "reg_started_at": "2024-07-16 16:08:03.534410197 +0000 UTC",
                    "processor_sns": "154.12.246.129:24444,154.12.246.131:24444,154.12.246.128:24444,",
                    "finished_at": "2024-07-16 16:12:23.515338118 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 29,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.006",
                    "activation_attempt_at": "2024-07-16 16:45:44.123257433 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        },
        {
            "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.007",
            "upload_timestamp": "2024-07-16 15:50:40.395810811 +0000 UTC",
            "file_index": "6",
            "base_file_id": "NwWDvdIg",
            "task_id": "CXbK5zVw",
            "reg_txid": "18480f0941684181b8c227b8a48aa7c36751c80218be6b435b03ad2c4d7ab209",
            "activation_txid": "1f817f66fb2bec073a4b4db720a539633fa808e970b1e8e67bf048760ae27428",
            "req_burn_txn_amount": 34,
            "burn_txn_id": "6e33a3e055203c0d0fd7847aedffac50cf45622927700dac5f843cfe655977ec",
            "req_amount": 170,
            "is_concluded": true,
            "cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
            "uuid_key": "40aeb7fd-165e-42eb-a127-a7c8628764a1",
            "hash_of_original_big_file": "dd5e035fb0ce7afeede98081ae808b9cfd0d3cc94279d202599b67b29ff2c386",
            "name_of_original_big_file_with_ext": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz",
            "size_of_original_big_file": 129965900,
            "start_block": 97036,
            "done_block": 97050,
            "registration_attempts": [
                {
                    "id": 33,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.007",
                    "reg_started_at": "2024-07-16 16:08:03.568282649 +0000 UTC",
                    "processor_sns": "154.12.246.128:24444,154.38.164.104:24444,154.12.246.130:24444,",
                    "finished_at": "2024-07-16 16:12:49.95335445 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ],
            "activation_attempts": [
                {
                    "id": 28,
                    "file_id": "google-cloud-cli-475.0.0-linux-x86_64.tar.gz.7z.007",
                    "activation_attempt_at": "2024-07-16 16:45:41.925499054 +0000 UTC",
                    "is_successful": true,
                    "error_message": ""
                }
            ]
        }
    ]
}
```

you see,  following 2 fields are important to look out for! 

```
"is_concluded": true,
"cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
```

If any of the volume registration fails after 3 attempts, the entire registration would obviously fail. currently, max_attempts is set to 3, meaning, for each volume, walletnode will attempt to register it 3 times, if it fails 3 times, its considered as a failure. Similar is the case for activation of registered tickets. Activation could fail for a number of reasons like user’s wallet being out of balance etc. 

Once all volumes get registered and activated, walletnode creates a cascade_multi_volume_metadata contract ticket and he txid is returned in this field as posted above as well. 

```
"cascade_metadata_ticket_id": "4d11ce243a64de54cf918a2f759fe52ff79768be44a66b6157040de3ac3eb15e",
```

This txid should be used to download the file! 
BUT! existing download endpoint would not support cascade_metadata_ticket download.

### 4. Download Multi-Volume Ticket's Artifact
There’s a new endpoint to download multi-volume ticket’s artifact.

#### Endpoint
GET `/openapi/cascade/v2/download?pid=<<pastel_id>>&txid=<<txid>>`
#### Sample Response


```
{
    "file_id": "275852"
}
```
This will kickoff the download process and then immediately reply with this file_id.

### 5. Check Download Status
Use the following endpoint to check the status of the download process using the file_id returned above.

#### Endpoint

GET `/openapi/cascade/downloads/<<file_id>>/status`
#### Sample Response
Sample response in case file is already downloaded or in progress.


```
{
    "total_volumes": 0,
    "downloaded_volumes": 0,
    "volumes_download_in_progress": 0,
    "volumes_pending_download": 0,
    "volumes_download_failed": 0,
    "task_status": "Completed",
    "size_of_the_file_megabytes": 124,
    "data_downloaded_megabytes": 124,
    "details": null
}
```

see, `task_status`
its an enum with possible values:
##### 1- Completed
##### 2- In Progress
##### 3- Failed

if the task is In Progress - you'd see the details populated - here’s the sample response in that case.

```

{
    "total_volumes": 7,
    "downloaded_volumes": 2,
    "volumes_download_in_progress": 2,
    "volumes_pending_download": 3,
    "volumes_download_failed": 0,
    "task_status": "In Progress",
    "size_of_the_file_megabytes": 124,
    "data_downloaded_megabytes": 40,
    "details": {
        "fields": {
            "047fa18ed022e8a36dcb4b08b1cf16c167aa3e47c21aa20369111b602520cc0b": {
                "Txid": "047fa18ed022e8a36dcb4b08b1cf16c167aa3e47c21aa20369111b602520cc0b",
                "TaskID": "",
                "Status": "Downloaded",
                "Error": ""
            },
            "18480f0941684181b8c227b8a48aa7c36751c80218be6b435b03ad2c4d7ab209": {
                "Txid": "18480f0941684181b8c227b8a48aa7c36751c80218be6b435b03ad2c4d7ab209",
                "TaskID": "m2OEE4D1",
                "Status": "In Progress",
                "Error": ""
            },
            "1a567609281db81f5f119d6c2b5e6e8df2976cc1b450ddd1d239712534084d01": {
                "Txid": "1a567609281db81f5f119d6c2b5e6e8df2976cc1b450ddd1d239712534084d01",
                "TaskID": "",
                "Status": "Pending",
                "Error": ""
            },
            "1d87d071cc83b7917b371e8f040e09adadbd07d92fccc85df0e4115369b0e988": {
                "Txid": "1d87d071cc83b7917b371e8f040e09adadbd07d92fccc85df0e4115369b0e988",
                "TaskID": "lO2UWkq7",
                "Status": "In Progress",
                "Error": ""
            },
            "4afb36a8d91d619341dfb932ff5ca0079c8bddb5e3bd35f1f3ac365855a249c8": {
                "Txid": "4afb36a8d91d619341dfb932ff5ca0079c8bddb5e3bd35f1f3ac365855a249c8",
                "TaskID": "",
                "Status": "Pending",
                "Error": ""
            },
            "7195a4b67d82f242f7f3979407b072ac0ca7e8bca8c157ad0cfdcbc34176f358": {
                "Txid": "7195a4b67d82f242f7f3979407b072ac0ca7e8bca8c157ad0cfdcbc34176f358",
                "TaskID": "",
                "Status": "Pending",
                "Error": ""
            },
            "a777884586342015a47fc9d223fa9c3ee866990acd770795eba56ea32a980a2e": {
                "Txid": "a777884586342015a47fc9d223fa9c3ee866990acd770795eba56ea32a980a2e",
                "TaskID": "",
                "Status": "Downloaded",
                "Error": ""
            }
        }
    }
}
```
Details of each volume and individual and overall progress will be returned as shown.

