package meta

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	_ "go.etcd.io/etcd/client/v3/concurrency"
)

type Client struct{ cli *clientv3.Client }

func New(endpoints []string) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints, DialTimeout: 3 * time.Second})
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) keyObj(bucket, key string) string {
	return fmt.Sprintf("/titan/obj/%s/%s", bucket, key)
}

// CAS 创建对象元数据骨架（If-None-Match:*）
func (c *Client) PutObject(ctx context.Context, bucket, key, contentType string, size int64) (version string, err error) {
	k := c.keyObj(bucket, key)
	version = fmt.Sprintf("v-%d", time.Now().UnixNano())
	obj := Object{Bucket: bucket, Key: key, Version: version, Size: size, ContentType: contentType}
	b, _ := json.Marshal(obj)
	// 仅当不存在时写入
	txn := c.cli.Txn(ctx).If(clientv3.Compare(clientv3.Version(k), "=", 0)).Then(clientv3.OpPut(k, string(b)))
	resp, err := txn.Commit()
	if err != nil {
		return "", err
	}
	if !resp.Succeeded {
		return "", fmt.Errorf("precondition failed")
	}
	return version, nil
}

// 完成写入：更新 chunks + version hash（乐观并发）
func (c *Client) CommitObject(ctx context.Context, bucket, key string, chunks []ChunkRef) (string, error) {
	k := c.keyObj(bucket, key)
	get, err := c.cli.Get(ctx, k)
	if err != nil {
		return "", err
	}
	if len(get.Kvs) == 0 {
		return "", fmt.Errorf("not found")
	}
	var obj Object
	if err := json.Unmarshal(get.Kvs[0].Value, &obj); err != nil {
		return "", err
	}
	obj.Chunks = chunks
	// 新版本号=sha256(chunks)
	h := sha256.New()
	for _, c := range chunks {
		h.Write(c.ChunkID)
	}
	obj.Version = fmt.Sprintf("sha256:%x", h.Sum(nil)[:8])
	b, _ := json.Marshal(obj)
	txn := c.cli.Txn(ctx).If(clientv3.Compare(clientv3.ModRevision(k), "=", get.Kvs[0].ModRevision)).Then(clientv3.OpPut(k, string(b)))
	resp, err := txn.Commit()
	if err != nil {
		return "", err
	}
	if !resp.Succeeded {
		return "", fmt.Errorf("conflict")
	}
	return obj.Version, nil
}

// 这里应有：AllocatePlacement(k,m) 从 /titan/nodes/ 读取在线节点，做一致性哈希/反亲和，返回 data/parity nodes
