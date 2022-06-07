package provider

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
)

func (conn *EtcdConnection) putKeyWithRetries(key string, val string, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	_, err := conn.Client.Put(ctx, key, val)
	if err != nil {
		etcdErr, ok := err.(rpctypes.EtcdError)
		if !ok {
			return err
		}
		
		if etcdErr.Code() != codes.Unavailable || retries <= 0 {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.putKeyWithRetries(key, val, retries - 1)
	}
	return nil
}

func (conn *EtcdConnection) PutKey(key string, val string) error {
	return conn.putKeyWithRetries(key, val, conn.Retries)
}

func (conn *EtcdConnection) getKeyWithRetries(key string, retries int) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	getRes, err := conn.Client.Get(ctx, key)

	if err != nil {
		etcdErr, ok := err.(rpctypes.EtcdError)
		if !ok {
			return "", false, err
		}
		
		if etcdErr.Code() != codes.Unavailable || retries <= 0 {
			return "", false, err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.getKeyWithRetries(key, retries - 1)
	}

	if len(getRes.Kvs) == 0 {
		return "", false, nil
	}

	return string(getRes.Kvs[0].Value), true, nil
}

func (conn *EtcdConnection) GetKey(key string) (string, bool, error) {
	return conn.getKeyWithRetries(key, conn.Retries)
}

func (conn *EtcdConnection) deleteKeyWithRetries(key string, retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	_, err := conn.Client.Delete(ctx, key)
	if err != nil {
		etcdErr, ok := err.(rpctypes.EtcdError)
		if !ok {
			return err
		}
		
		if etcdErr.Code() != codes.Unavailable || retries <= 0 {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteKeyWithRetries(key, retries - 1)
	}

	return nil
}

func (conn *EtcdConnection) DeleteKey(key string) error {
	return conn.deleteKeyWithRetries(key, conn.Retries)
}

func (conn *EtcdConnection) deleteKeyRangeWithRetries(key string, rangeEnd string , retries int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conn.Timeout)*time.Second)
	defer cancel()

	_, err := conn.Client.Delete(ctx, key, clientv3.WithRange(rangeEnd))
	if err != nil {
		etcdErr, ok := err.(rpctypes.EtcdError)
		if !ok {
			return err
		}
		
		if etcdErr.Code() != codes.Unavailable || retries <= 0 {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return conn.deleteKeyRangeWithRetries(key, rangeEnd, retries - 1)
	}

	return nil
}

func (conn *EtcdConnection) DeleteKeyRange(key string, rangeEnd string) error {
	return conn.deleteKeyRangeWithRetries(key, rangeEnd, conn.Retries)
}