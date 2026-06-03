package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Jeremiah-Williams1/inference-orchestrator/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{client: client}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *models.Job) error {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("queue:%s", job.Type)
	err = q.client.LPush(ctx, key, jobBytes).Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *RedisQueue) Dequeue(ctx context.Context, jobType models.Type) (*models.Job, error) {
	key := fmt.Sprintf("queue:%s", jobType)

	timeout := 1 * time.Second
	result, err := q.client.BRPop(ctx, timeout, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	val := result[1]

	var job models.Job
	err = json.Unmarshal([]byte(val), &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (q *RedisQueue) SetResult(ctx context.Context, job *models.Job) error {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	ttl := 1 * time.Hour
	key := fmt.Sprintf("job:%s", job.ID)
	err = q.client.Set(ctx, key, jobBytes, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *RedisQueue) GetResult(ctx context.Context, jobID string) (*models.Job, error) {
	key := fmt.Sprintf("job:%s", jobID)
	val, err := q.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var job models.Job
	err = json.Unmarshal([]byte(val), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (q *RedisQueue) Depth(ctx context.Context, jobType models.Type) (int64, error) {
	key := fmt.Sprintf("queue:%s", jobType)

	length, err := q.client.LLen(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return length, nil
}
