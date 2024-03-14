package kademlia

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (s *testSuite) TestHasBit() {
	for i := 0; i < 8; i++ {
		s.Equal(true, hasBit(byte(255), uint(i)))
	}
}

func (s *testSuite) newNodeList(n int) (*NodeList, error) {
	nl := &NodeList{}

	for i := 0; i < n; i++ {
		id, err := newRandomID()
		if err != nil {
			return nil, err
		}
		nl.AddNodes([]*Node{
			{
				ID:   id,
				IP:   s.IP,
				Port: s.bootstrapPort + i + 1,
			},
		})
	}
	return nl, nil
}

func (s *testSuite) TestAddNodes() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())
}

func (s *testSuite) TestDelNodes() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())

	ids := [][]byte{}
	for _, node := range nl.Nodes {
		ids = append(ids, node.ID)
	}
	for _, id := range ids {
		nl.DelNode(&Node{ID: id})
	}
	s.Equal(0, nl.Len())
}

func (s *testSuite) TestExists() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())

	for _, node := range nl.Nodes {
		s.Equal(true, nl.Exists(node))
	}
}

func CreateNodeList(data string) *NodeList {
	records := strings.Split(data, ",")
	nodelist := &NodeList{}

	for _, record := range records {
		parts := strings.Split(record, "-")
		if len(parts) != 2 {
			continue
		}

		nodeID := []byte(strings.TrimSpace(parts[0]))
		ipAndPort := strings.Split(parts[1], ":")
		if len(ipAndPort) != 2 {
			continue
		}

		ip := ipAndPort[0]
		var port int
		_, err := fmt.Sscanf(ipAndPort[1], "%d", &port)
		if err != nil {
			continue
		}

		node := &Node{
			ID:   nodeID,
			IP:   strings.TrimSpace(ip),
			Port: port,
		}
		node.SetHashedID()
		nodelist.Nodes = append(nodelist.Nodes, node)
	}

	return nodelist
}

func TestNodeListSortDetail(t *testing.T) {
	hashedComparator, err := hex.DecodeString("589ed5c2cca65b69c5229e7e17a85c9a2dc04fffde01b9870e8b6635bf0d073e")
	if err != nil {
		t.Fatal(err)
	}

	results := []*NodeList{}
	tests := []struct {
		name       string
		data       string
		comparator []byte
	}{
		{
			name:       "45",
			comparator: hashedComparator,
			data: `3YjK6T7Cc9P29MdUnmZAxv1qYCGcRsNPh4mGmmXYyfVAtKCEjdBGegAMLiKcW47hi4gpGMWE9peM3DRFj7RWwST3x3EP1BrqKjDcdMJBc7xVAF9h1wDwaY-3.19.48.187:14445,
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
		},
		{
			name:       "77",
			comparator: hashedComparator,
			data: `3YjKAxUb4bbnK6mHPFQe1fmSoSssLBcQjgNEmTWgVWssnk4hkmr3Y8Yd9D1UNJ7PSvnfuwE5NB78G5XxcueqnCfDngfB1ZyuYiAiVgi7e9ZgTgtun8nLS9-149.102.147.66:14445,
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
			3YjKAqro6htwGLENofxUwJ1V1QSsdcPGJ2v4A4zb8AphwVACizg8wFTq4jsJfCGhRznPYb86HqpYFVmGvgDiq6EASrHxFkbkX37JzfDGntvVGLinvZyqz4-154.12.235.19:14445,
			3YjK6TzKc4beYhe8PQLtMfqFLicDS4PkhZJrb5TgzJQ7EBvRZjfBbo3dVu9SqFd2m6rupfrr5hZEHhNHyySjC5ZozjuqdJfTKEXYzJ4KYWFFmKcwTtiVGe-149.102.147.180:14445,
			3YjK9rsqRn8jXC4T48J1tdeEGC2XxHw4YNdfLX8PvZyw53ykNdXYU7WSEPHSNkDyg3npdxAPG57vFqLs9yZzrzfrCSvt6uFaUTLV6ReXkVWdZ8SY6YfeHu-95.111.235.112:14445,
			3YjKA64Z5iLhRN7NEU1P444nMyjAsGtHR3RgdoXMV14evXh6iox3qWDudkbDnwFztk5h8NHQZPZxyFwyv6USyFyoxG8YT3qB1NLU4V1XFFEDkve4ie5gcE-154.12.235.12:14445,
			3YjKjNzpsUY2RqPqVifHHHtqDstMKExtTjQzT97wozZyDgn9HA3ckxu4dTYdsxw9isrqwMtGjMkPav344cfApLfsMrRBkZRR3n2mRZceY2xjiqnBDBfh3u-149.102.147.122:14445,
			3YjK6S5YcNWvGUP8YdDKYNXm4xKbXrhd4jkSeD7ZaPgo5tfCWAUPY16fZh7Dp2hCDz1rg35EBgrL2U67pSU3JjQYb4GAjKL2MqsdJ3HM5diQFZTbjZTVR7-154.12.235.16:14445,
			3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa-3.142.100.245:14445`,
		},
	}

	// Pre-compute the hashed IDs for the nodes.

	nl1 := CreateNodeList(tests[0].data)
	nl2 := CreateNodeList(tests[1].data)

	nl1.Comparator = tests[0].comparator
	nl2.Comparator = tests[1].comparator
	nl1.debug = true
	nl2.debug = true

	// check if both list have same nodes
	for i := 0; i < nl1.Len(); i++ {
		exists := false

		for j := 0; j < nl2.Len(); j++ {
			if nl1.Nodes[i].String() == nl2.Nodes[j].String() {
				exists = true
			}
		}

		if !exists {
			t.Fatalf("node doesn't exist in nl2 %s", nl1.Nodes[i].String())
		}
	}

	nl1.Sort()
	results = append(results, nl1)

	nl2.Sort()
	results = append(results, nl2)
	first := results[0]
	for i := 1; i < len(results); i++ {
		if len(first.Nodes) != len(results[i].Nodes) {
			t.Errorf("expected %d nodes but got %d", len(first.Nodes), len(results[i].Nodes))
		}

		for j := 0; j < len(first.Nodes); j++ {
			if !bytes.Equal(first.Nodes[j].ID, results[i].Nodes[j].ID) {
				t.Errorf("expected ID at index %d to be %v but got %v", j, string(first.Nodes[j].ID), string(results[i].Nodes[j].ID))
			}
		}
	}
}

func TestNodeListSort(t *testing.T) {
	hashedComparator, err := hex.DecodeString("a6cc846e5019e845fa58a5af1c16baf96d7921de13375d05e745f266520734ba")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		data       string
		comparator []byte
	}{
		{
			name:       "A",
			comparator: hashedComparator,
			data: `3YjK6T7Cc9P29MdUnmZAxv1qYCGcRsNPh4mGmmXYyfVAtKCEjdBGegAMLiKcW47hi4gpGMWE9peM3DRFj7RWwST3x3EP1BrqKjDcdMJBc7xVAF9h1wDwaY-3.19.48.187:14445,
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
		},
		{
			name:       "B",
			comparator: hashedComparator,
			data: `3YjKAxUb4bbnK6mHPFQe1fmSoSssLBcQjgNEmTWgVWssnk4hkmr3Y8Yd9D1UNJ7PSvnfuwE5NB78G5XxcueqnCfDngfB1ZyuYiAiVgi7e9ZgTgtun8nLS9-149.102.147.66:14445,
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
			3YjKAqro6htwGLENofxUwJ1V1QSsdcPGJ2v4A4zb8AphwVACizg8wFTq4jsJfCGhRznPYb86HqpYFVmGvgDiq6EASrHxFkbkX37JzfDGntvVGLinvZyqz4-154.12.235.19:14445,
			3YjK6TzKc4beYhe8PQLtMfqFLicDS4PkhZJrb5TgzJQ7EBvRZjfBbo3dVu9SqFd2m6rupfrr5hZEHhNHyySjC5ZozjuqdJfTKEXYzJ4KYWFFmKcwTtiVGe-149.102.147.180:14445,
			3YjK9rsqRn8jXC4T48J1tdeEGC2XxHw4YNdfLX8PvZyw53ykNdXYU7WSEPHSNkDyg3npdxAPG57vFqLs9yZzrzfrCSvt6uFaUTLV6ReXkVWdZ8SY6YfeHu-95.111.235.112:14445,
			3YjKA64Z5iLhRN7NEU1P444nMyjAsGtHR3RgdoXMV14evXh6iox3qWDudkbDnwFztk5h8NHQZPZxyFwyv6USyFyoxG8YT3qB1NLU4V1XFFEDkve4ie5gcE-154.12.235.12:14445,
			3YjKjNzpsUY2RqPqVifHHHtqDstMKExtTjQzT97wozZyDgn9HA3ckxu4dTYdsxw9isrqwMtGjMkPav344cfApLfsMrRBkZRR3n2mRZceY2xjiqnBDBfh3u-149.102.147.122:14445,
			3YjK6S5YcNWvGUP8YdDKYNXm4xKbXrhd4jkSeD7ZaPgo5tfCWAUPY16fZh7Dp2hCDz1rg35EBgrL2U67pSU3JjQYb4GAjKL2MqsdJ3HM5diQFZTbjZTVR7-154.12.235.16:14445,
			3YjK6RMKAqFPD99e6EMt514wLRGqoFnKWhqQzpxRZowTHZPVFAuRB2fSLU2eucNyYBrFowFWBQ9hm3x4JBdCSesMJD1vcFYywxohmXET3NwCm3dBTBijxa-3.142.100.245:14445`,
		},
	}

	// Pre-compute the hashed IDs for the nodes.

	nl1 := CreateNodeList(tests[0].data)
	nl2 := CreateNodeList(tests[1].data)

	nl1.Comparator = tests[0].comparator
	nl2.Comparator = tests[1].comparator
	nl1.debug = true
	nl2.debug = true

	nl1.Sort()
	nl1.TopN(6)

	nl2.Sort()
	nl2.TopN(6)

	assert.Equal(t, nl1.String(), nl2.String())
}
