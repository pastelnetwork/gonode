package kademlia

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClosestContactsWithIncludingNode(t *testing.T) {

	tests := map[string]struct {
		nodes           string
		hashTable       *HashTable
		targetKey       string
		includedNode    *Node
		ignoredNodes    []*Node
		topCount        int
		expectedResults *NodeList
	}{
		"uc1": {
			nodes: `3YjK6T7Cc9P29MdUnmZAxv1qYCGcRsNPh4mGmmXYyfVAtKCEjdBGegAMLiKcW47hi4gpGMWE9peM3DRFj7RWwST3x3EP1BrqKjDcdMJBc7xVAF9h1wDwaY-3.19.48.187:14445,
			3YjK9rsqRn8jXC4T48J1tdeEGC2XxHw4YNdfLX8PvZyw53ykNdXYU7WSEPHSNkDyg3npdxAPG57vFqLs9yZzrzfrCSvt6uFaUTLV6ReXkVWdZ8SY6YfeHu-95.111.235.112:14445,
			3YjKEC59G1hzBWWhGQivGqvJd27Hj7yyrqBjzBKbYHfKj8ycVTRQjGpMzgWZfucknYvcN1uwtrP1Fr7SMSPnmuBaeCejRfzFoeJPKZRLm76xVWcsAH8HZz-18.189.251.53:14445,
			3YjKEip9oXJMm8YFiUkhBgv8DbyAjWhzJMZaWqWs3KWhdvnKxQj77gYy2zohN5Fx8Yibzf3v1UbovnUvBTKq64caWDFLEVLMpLVVGFybNARM694t1Wz3h1-3.18.240.77:14445,
			3YjKAnECn9CpGB3kqKEKiEKoisqpyWVYKzDNdQXzTZ97MH4hpg3nyZknEGb3vszNMccQNvGz51yQPtjs2BwWJpCYwPmzJJbs8s7Vy7qiUewaK5Ce5UGBxL-38.242.151.234:14445,
			3YjKAqro6htwGLENofxUwJ1V1QSsdcPGJ2v4A4zb8AphwVACizg8wFTq4jsJfCGhRznPYb86HqpYFVmGvgDiq6EASrHxFkbkX37JzfDGntvVGLinvZyqz4-154.12.235.19:14445,
			3YjKjB1PaGGdsWGoASNhBMs3XxGfFRV5qzqUFD31r8AwtAUcxYZXtnG7AGicRcbKS1LwD79dtRn1f9euS41wNeZNgcUM7GyHKeyU1RmNJbygethBimm9ei-3.130.254.15:14445,
			3YjKAtueoP1xMZBtH98xGKxDRmqTds58VHS9ev4uZSeTjRUeWyd1qfD99FNkJoKsQxZonjSaSdTNRPWTNqkGhXMxLZjhjzxTHryuJhMQRKJSs97mb2U4im-154.38.161.44:14445,
			3YjKEDhanRTxuP1WvEg4PFsq1ejnQuS1NhHRBZWWzxFq3dcUEwKeRdjRBrVU4b9BobvX6LesQRRNYpELhXccYzPmeNiLi69eFCGh4MZ1jK5zU24RSJo2DJ-149.102.147.113:14445,
			3YjK6mqfZDNgbjSaTqPYDtkqsPm4XQ1vezUn7NYqSekpV3c9q9Qjndxaq3UcEmkBEih8iajd7iboPRbAjABqz7Ri6fCdXvETrhU1tFfUZHeJK8vNmZju6D-154.38.160.234:14445,
			3YjK6FkgW1eZvnVpcbeR5LyazNCmRNKg59yHHdpV52sAHp7c2Np9TVQeoCvZmRciQndXXbSZhVW81mEUunXd6apUbv8Df9Lk7sLRvMMdVd1nWQEYRGi1tu-95.111.253.200:14445,
			3YjKA64Z5iLhRN7NEU1P444nMyjAsGtHR3RgdoXMV14evXh6iox3qWDudkbDnwFztk5h8NHQZPZxyFwyv6USyFyoxG8YT3qB1NLU4V1XFFEDkve4ie5gcE-154.12.235.12:14445,
			3YjKDzhn3GVoH5GKfteY39YDYEfzqXpmKDj4uUfEVb2gPSE8n5TgYTg6QDNMZSBsngCT2vrsbvyZf1K9G4FBzKWdRTyLnkpWp2iLTjCvRtb1W8Z2wf7pGg-154.12.235.41:14445,
			3YjKA58AQ4NPho7kKBCnBUzfUYpurEmyehRToFR7kTZ26JnuJSR2xQAhxzrBHmMwZ4EFCPbkFmde4QVK4GpUEBdxMgoBohHFnGN3hZCwnnD1VQ6e8PHLX2-154.38.161.58:14445,
			3YjKAkKCwWmoRhHxMnvoyvrg3ncUXcwk2heLuJbjVV2GxxLK5gfWb4A7sCse6ZMEbAu3odVaqguPcZxZRs8fAB2LnSezZjTF9MFMs1zm5de6g8nNm1gGyc-154.38.161.52:14445,
			3YjK6TzKc4beYhe8PQLtMfqFLicDS4PkhZJrb5TgzJQ7EBvRZjfBbo3dVu9SqFd2m6rupfrr5hZEHhNHyySjC5ZozjuqdJfTKEXYzJ4KYWFFmKcwTtiVGe-149.102.147.180:14445,
			3YjK6S5YcNWvGUP8YdDKYNXm4xKbXrhd4jkSeD7ZaPgo5tfCWAUPY16fZh7Dp2hCDz1rg35EBgrL2U67pSU3JjQYb4GAjKL2MqsdJ3HM5diQFZTbjZTVR7-154.12.235.16:14445,
			3YjKjhwUZ1ZskdmJ1n9rU7tKqkgwVQnNcKzaK6rgJPGvEXf4aqGS3boPaH5T2HwdkZYMHEW8YpHASLZTQR3DY3kgJwFBvpUMc2iKecmAp1xymqiyNh4NJ1-38.242.151.208:14445,
			3YjKAxUb4bbnK6mHPFQe1fmSoSssLBcQjgNEmTWgVWssnk4hkmr3Y8Yd9D1UNJ7PSvnfuwE5NB78G5XxcueqnCfDngfB1ZyuYiAiVgi7e9ZgTgtun8nLS9-149.102.147.66:14445,
			3YjKjNzpsUY2RqPqVifHHHtqDstMKExtTjQzT97wozZyDgn9HA3ckxu4dTYdsxw9isrqwMtGjMkPav344cfApLfsMrRBkZRR3n2mRZceY2xjiqnBDBfh3u-149.102.147.122:14445,
			3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa-3.142.100.245:14445`,
			targetKey:    "d60dfe8bcdb843f8c1bb2e0eebffebf54d8df17e5067e4d6cbeb08055b0aa9a2",
			includedNode: nil,
			ignoredNodes: make([]*Node, 0),
			topCount:     6,
			expectedResults: &NodeList{Nodes: []*Node{
				{IP: "3.18.240.77"},
				{IP: "38.242.151.208"},
				{IP: "154.38.160.234"},
				{IP: "154.38.161.58"},
				{IP: "154.12.235.12"},
				{IP: "38.242.151.234"},
			},
			},
		},
		"uc2": {
			nodes: `3YjK6TzKc4beYhe8PQLtMfqFLicDS4PkhZJrb5TgzJQ7EBvRZjfBbo3dVu9SqFd2m6rupfrr5hZEHhNHyySjC5ZozjuqdJfTKEXYzJ4KYWFFmKcwTtiVGe-149.102.147.180:14445,,
			3YjK6T7Cc9P29MdUnmZAxv1qYCGcRsNPh4mGmmXYyfVAtKCEjdBGegAMLiKcW47hi4gpGMWE9peM3DRFj7RWwST3x3EP1BrqKjDcdMJBc7xVAF9h1wDwaY-3.19.48.187:14445,
			3YjKEip9oXJMm8YFiUkhBgv8DbyAjWhzJMZaWqWs3KWhdvnKxQj77gYy2zohN5Fx8Yibzf3v1UbovnUvBTKq64caWDFLEVLMpLVVGFybNARM694t1Wz3h1-3.18.240.77:14445,
			3YjKAxUb4bbnK6mHPFQe1fmSoSssLBcQjgNEmTWgVWssnk4hkmr3Y8Yd9D1UNJ7PSvnfuwE5NB78G5XxcueqnCfDngfB1ZyuYiAiVgi7e9ZgTgtun8nLS9-149.102.147.66:14445,
			3YjKAnECn9CpGB3kqKEKiEKoisqpyWVYKzDNdQXzTZ97MH4hpg3nyZknEGb3vszNMccQNvGz51yQPtjs2BwWJpCYwPmzJJbs8s7Vy7qiUewaK5Ce5UGBxL-38.242.151.234:14445,
			3YjKAqro6htwGLENofxUwJ1V1QSsdcPGJ2v4A4zb8AphwVACizg8wFTq4jsJfCGhRznPYb86HqpYFVmGvgDiq6EASrHxFkbkX37JzfDGntvVGLinvZyqz4-154.12.235.19:14445,
			3YjKjB1PaGGdsWGoASNhBMs3XxGfFRV5qzqUFD31r8AwtAUcxYZXtnG7AGicRcbKS1LwD79dtRn1f9euS41wNeZNgcUM7GyHKeyU1RmNJbygethBimm9ei-3.130.254.15:14445,
			3YjKAtueoP1xMZBtH98xGKxDRmqTds58VHS9ev4uZSeTjRUeWyd1qfD99FNkJoKsQxZonjSaSdTNRPWTNqkGhXMxLZjhjzxTHryuJhMQRKJSs97mb2U4im-154.38.161.44:14445,
			3YjKEDhanRTxuP1WvEg4PFsq1ejnQuS1NhHRBZWWzxFq3dcUEwKeRdjRBrVU4b9BobvX6LesQRRNYpELhXccYzPmeNiLi69eFCGh4MZ1jK5zU24RSJo2DJ-149.102.147.113:14445,
			3YjKEC59G1hzBWWhGQivGqvJd27Hj7yyrqBjzBKbYHfKj8ycVTRQjGpMzgWZfucknYvcN1uwtrP1Fr7SMSPnmuBaeCejRfzFoeJPKZRLm76xVWcsAH8HZz-18.189.251.53:14445,
			3YjK6mqfZDNgbjSaTqPYDtkqsPm4XQ1vezUn7NYqSekpV3c9q9Qjndxaq3UcEmkBEih8iajd7iboPRbAjABqz7Ri6fCdXvETrhU1tFfUZHeJK8vNmZju6D-154.38.160.234:14445,
			3YjKA64Z5iLhRN7NEU1P444nMyjAsGtHR3RgdoXMV14evXh6iox3qWDudkbDnwFztk5h8NHQZPZxyFwyv6USyFyoxG8YT3qB1NLU4V1XFFEDkve4ie5gcE-154.12.235.12:14445,
			3YjKDzhn3GVoH5GKfteY39YDYEfzqXpmKDj4uUfEVb2gPSE8n5TgYTg6QDNMZSBsngCT2vrsbvyZf1K9G4FBzKWdRTyLnkpWp2iLTjCvRtb1W8Z2wf7pGg-154.12.235.41:14445,
			3YjKA58AQ4NPho7kKBCnBUzfUYpurEmyehRToFR7kTZ26JnuJSR2xQAhxzrBHmMwZ4EFCPbkFmde4QVK4GpUEBdxMgoBohHFnGN3hZCwnnD1VQ6e8PHLX2-154.38.161.58:14445,
			3YjKAkKCwWmoRhHxMnvoyvrg3ncUXcwk2heLuJbjVV2GxxLK5gfWb4A7sCse6ZMEbAu3odVaqguPcZxZRs8fAB2LnSezZjTF9MFMs1zm5de6g8nNm1gGyc-154.38.161.52:14445,
			3YjK6FkgW1eZvnVpcbeR5LyazNCmRNKg59yHHdpV52sAHp7c2Np9TVQeoCvZmRciQndXXbSZhVW81mEUunXd6apUbv8Df9Lk7sLRvMMdVd1nWQEYRGi1tu-95.111.253.200:14445,
			3YjK9rsqRn8jXC4T48J1tdeEGC2XxHw4YNdfLX8PvZyw53ykNdXYU7WSEPHSNkDyg3npdxAPG57vFqLs9yZzrzfrCSvt6uFaUTLV6ReXkVWdZ8SY6YfeHu-95.111.235.112:14445,
			3YjK6S5YcNWvGUP8YdDKYNXm4xKbXrhd4jkSeD7ZaPgo5tfCWAUPY16fZh7Dp2hCDz1rg35EBgrL2U67pSU3JjQYb4GAjKL2MqsdJ3HM5diQFZTbjZTVR7-154.12.235.16:14445,
			3YjKjhwUZ1ZskdmJ1n9rU7tKqkgwVQnNcKzaK6rgJPGvEXf4aqGS3boPaH5T2HwdkZYMHEW8YpHASLZTQR3DY3kgJwFBvpUMc2iKecmAp1xymqiyNh4NJ1-38.242.151.208:14445,
			3YjKjNzpsUY2RqPqVifHHHtqDstMKExtTjQzT97wozZyDgn9HA3ckxu4dTYdsxw9isrqwMtGjMkPav344cfApLfsMrRBkZRR3n2mRZceY2xjiqnBDBfh3u-149.102.147.122:14445,
			3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa-3.142.100.245:14445`,
			targetKey:    "ca2fec47d42d516b719b927f6748ff0929e814a60f41663e33af1c68e540a12b",
			includedNode: nil,
			ignoredNodes: make([]*Node, 0),
			topCount:     6,
			expectedResults: &NodeList{Nodes: []*Node{
				{IP: "154.38.161.58"},
				{IP: "154.38.160.234"},
				{IP: "38.242.151.208"},
				{IP: "3.18.240.77"},
				{IP: "154.12.235.12"},
				{IP: "149.102.147.122"},
			},
			},
		},
		"uc3": {
			nodes: `3YjKAxUb4bbnK6mHPFQe1fmSoSssLBcQjgNEmTWgVWssnk4hkmr3Y8Yd9D1UNJ7PSvnfuwE5NB78G5XxcueqnCfDngfB1ZyuYiAiVgi7e9ZgTgtun8nLS9-149.102.147.66:14445,
			3YjKEC59G1hzBWWhGQivGqvJd27Hj7yyrqBjzBKbYHfKj8ycVTRQjGpMzgWZfucknYvcN1uwtrP1Fr7SMSPnmuBaeCejRfzFoeJPKZRLm76xVWcsAH8HZz-18.189.251.53:14445,
			3YjKA58AQ4NPho7kKBCnBUzfUYpurEmyehRToFR7kTZ26JnuJSR2xQAhxzrBHmMwZ4EFCPbkFmde4QVK4GpUEBdxMgoBohHFnGN3hZCwnnD1VQ6e8PHLX2-154.38.161.58:14445,
			3YjK6mqfZDNgbjSaTqPYDtkqsPm4XQ1vezUn7NYqSekpV3c9q9Qjndxaq3UcEmkBEih8iajd7iboPRbAjABqz7Ri6fCdXvETrhU1tFfUZHeJK8vNmZju6D-154.38.160.234:14445,
			3YjKDzhn3GVoH5GKfteY39YDYEfzqXpmKDj4uUfEVb2gPSE8n5TgYTg6QDNMZSBsngCT2vrsbvyZf1K9G4FBzKWdRTyLnkpWp2iLTjCvRtb1W8Z2wf7pGg-154.12.235.41:14445,
			3YjKAkKCwWmoRhHxMnvoyvrg3ncUXcwk2heLuJbjVV2GxxLK5gfWb4A7sCse6ZMEbAu3odVaqguPcZxZRs8fAB2LnSezZjTF9MFMs1zm5de6g8nNm1gGyc-154.38.161.52:14445,
			3YjK6T7Cc9P29MdUnmZAxv1qYCGcRsNPh4mGmmXYyfVAtKCEjdBGegAMLiKcW47hi4gpGMWE9peM3DRFj7RWwST3x3EP1BrqKjDcdMJBc7xVAF9h1wDwaY-3.19.48.187:14445,
			3YjKjB1PaGGdsWGoASNhBMs3XxGfFRV5qzqUFD31r8AwtAUcxYZXtnG7AGicRcbKS1LwD79dtRn1f9euS41wNeZNgcUM7GyHKeyU1RmNJbygethBimm9ei-3.130.254.15:14445,
			3YjKEip9oXJMm8YFiUkhBgv8DbyAjWhzJMZaWqWs3KWhdvnKxQj77gYy2zohN5Fx8Yibzf3v1UbovnUvBTKq64caWDFLEVLMpLVVGFybNARM694t1Wz3h1-3.18.240.77:14445,
			3YjKjhwUZ1ZskdmJ1n9rU7tKqkgwVQnNcKzaK6rgJPGvEXf4aqGS3boPaH5T2HwdkZYMHEW8YpHASLZTQR3DY3kgJwFBvpUMc2iKecmAp1xymqiyNh4NJ1-38.242.151.208:14445,
			3YjKAtueoP1xMZBtH98xGKxDRmqTds58VHS9ev4uZSeTjRUeWyd1qfD99FNkJoKsQxZonjSaSdTNRPWTNqkGhXMxLZjhjzxTHryuJhMQRKJSs97mb2U4im-154.38.161.44:14445,
			3YjKEDhanRTxuP1WvEg4PFsq1ejnQuS1NhHRBZWWzxFq3dcUEwKeRdjRBrVU4b9BobvX6LesQRRNYpELhXccYzPmeNiLi69eFCGh4MZ1jK5zU24RSJo2DJ-149.102.147.113:14445,
			3YjK6FkgW1eZvnVpcbeR5LyazNCmRNKg59yHHdpV52sAHp7c2Np9TVQeoCvZmRciQndXXbSZhVW81mEUunXd6apUbv8Df9Lk7sLRvMMdVd1nWQEYRGi1tu-95.111.253.200:14445,
			3YjKAnECn9CpGB3kqKEKiEKoisqpyWVYKzDNdQXzTZ97MH4hpg3nyZknEGb3vszNMccQNvGz51yQPtjs2BwWJpCYwPmzJJbs8s7Vy7qiUewaK5Ce5UGBxL-38.242.151.234:14445,
			3YjK6TzKc4beYhe8PQLtMfqFLicDS4PkhZJrb5TgzJQ7EBvRZjfBbo3dVu9SqFd2m6rupfrr5hZEHhNHyySjC5ZozjuqdJfTKEXYzJ4KYWFFmKcwTtiVGe-149.102.147.180:14445,
			3YjK9rsqRn8jXC4T48J1tdeEGC2XxHw4YNdfLX8PvZyw53ykNdXYU7WSEPHSNkDyg3npdxAPG57vFqLs9yZzrzfrCSvt6uFaUTLV6ReXkVWdZ8SY6YfeHu-95.111.235.112:14445,
			3YjKA64Z5iLhRN7NEU1P444nMyjAsGtHR3RgdoXMV14evXh6iox3qWDudkbDnwFztk5h8NHQZPZxyFwyv6USyFyoxG8YT3qB1NLU4V1XFFEDkve4ie5gcE-154.12.235.12:14445,
			3YjKjNzpsUY2RqPqVifHHHtqDstMKExtTjQzT97wozZyDgn9HA3ckxu4dTYdsxw9isrqwMtGjMkPav344cfApLfsMrRBkZRR3n2mRZceY2xjiqnBDBfh3u-149.102.147.122:14445,
			3YjK6S5YcNWvGUP8YdDKYNXm4xKbXrhd4jkSeD7ZaPgo5tfCWAUPY16fZh7Dp2hCDz1rg35EBgrL2U67pSU3JjQYb4GAjKL2MqsdJ3HM5diQFZTbjZTVR7-154.12.235.16:14445,
			3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa-3.142.100.245:14445`,
			includedNode: &Node{ID: []byte("3YjKAqro6htwGLENofxUwJ1V1QSsdcPGJ2v4A4zb8AphwVACizg8wFTq4jsJfCGhRznPYb86HqpYFVmGvgDiq6EASrHxFkbkX37JzfDGntvVGLinvZyqz4"),
				IP: "154.12.235.19", Port: 14445},
			ignoredNodes: make([]*Node, 0),
			topCount:     6,
			targetKey:    "517b1003b805793f35f4242f481a48cd45e5431a2af6e10339cbcf97f7b1a27e",
			expectedResults: &NodeList{Nodes: []*Node{
				{IP: "18.189.251.53"},
				{IP: "149.102.147.113"},
				{IP: "154.12.235.19"},
				{IP: "154.38.161.44"},
				{IP: "95.111.253.200"},
				{IP: "3.19.48.187"},
			},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			nl1 := CreateNodeList(tc.nodes)
			hashTable, err := NewHashTable(&Options{IP: "3.142.100.245", Port: 14445, ID: []byte("3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa")})
			assert.NoError(t, err)

			tc.hashTable = hashTable

			for _, node := range nl1.Nodes {
				node.SetHashedID()

				idx := tc.hashTable.bucketIndex(tc.hashTable.self.HashedID, node.HashedID)
				bucket := tc.hashTable.routeTable[idx]
				bucket = append(bucket, node)
				tc.hashTable.routeTable[idx] = bucket

			}

			key, err := hex.DecodeString(tc.targetKey)
			assert.NoError(t, err)

			if tc.includedNode != nil {
				tc.includedNode.SetHashedID()
			}
			result := tc.hashTable.closestContactsWithInlcudingNode(tc.topCount, key, tc.ignoredNodes, tc.includedNode)
			for i, node := range result.Nodes {

				assert.Equal(t, node.IP, tc.expectedResults.Nodes[i].IP)
			}
		})
	}
}
