package storage

import (
	"context"
	"encoding/json"
	"github.com/Frank-Way/note-go-rest-service/internal/note"
	"github.com/Frank-Way/note-go-rest-service/internal/note/nerror"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

var _ note.Storage = &redisStorage{}

type redisStorage struct {
	client *redis.Client
	logger *logrus.Logger
}

const nextIdKey = ".nextId2"

func NewRedisStorage(host, port, password string, db int, logger *logrus.Logger) (note.Storage, error) {
	addr := net.JoinHostPort(host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB: " + addr
		return nil, storeErr
	}
	return &redisStorage{
		client: client,
		logger: logger,
	}, nil
}

type userAggregate struct {
	Login   string `json:"login"`
	NoteIds []int  `json:"notes_ids"`
}

func (rs *redisStorage) Save(ctx context.Context, n note.Note) (string, error) {
	rs.logger.Info("save note to redis")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return "", storeErr
	}
	rs.logger.Debugf("get new id from redis for note: %v", n)
	nextIdStr, err := rs.client.Get(nextIdKey).Result()
	if err != nil {
		rs.logger.Debugf("error during getting next id: %v", err)
		rs.logger.Debug("set next id to 1")
		nextIdStr = "1"
		err := rs.client.Set(nextIdKey, "2", 0).Err()
		if err != nil {
			return "", err
		}
	}
	rs.logger.Debugf("parse new id %s for note: %v", nextIdStr, n)
	nextId, err := strconv.Atoi(nextIdStr)
	if err != nil {
		rs.logger.Debugf("error during parsing next id: %v", err)
		return "", err
	}
	rs.logger.Debugf("set new id %d to note %v", nextId, n)
	n.Id = int(nextId)
	rs.logger.Debug("incr id in redis")
	rs.client.Incr(nextIdKey)
	rs.logger.Debugf("marshaling note %v", n)
	bytes, err := json.Marshal(n)
	if err != nil {
		rs.logger.Debugf("error during marshaling note: %v", err)
		return "", err
	}
	rs.logger.Debugf("save note in redis: %v", n)
	if err := rs.client.Set(strconv.Itoa(int(n.Id)), bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving note: %v", err)
		return "", err
	}
	rs.logger.Debugf("get user's note aggregate for login: %q", n.Author)
	aggrStr, err := rs.client.Get(n.Author).Result()
	if err != nil {
		rs.logger.Debugf("error during getting aggregate: %v", err)
		aggrStr = "{\"login\":\"" + n.Author + "\",\"notes_ids\":[]}"
	}
	rs.logger.Debugf("unmarshal aggregate %q", aggrStr)
	var aggr userAggregate
	err = json.Unmarshal([]byte(aggrStr), &aggr)
	if err != nil {
		return "", err
	}
	rs.logger.Debugf("append note %v to aggregate %v", n, aggr)
	aggr.NoteIds = append(aggr.NoteIds, n.Id)
	rs.logger.Debugf("marshaling aggregate %v", aggr)
	bytes, err = json.Marshal(aggr)
	if err != nil {
		rs.logger.Debugf("error during marshaling aggregate: %v", err)
		return "", err
	}
	rs.logger.Debugf("save aggregate to redis: %v", aggr)
	if err = rs.client.Set(aggr.Login, bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving aggregate: %v", err)
		return "", err
	}
	return strconv.Itoa(int(n.Id)), nil
}

func (rs *redisStorage) GetById(ctx context.Context, id int) (note.Note, error) {
	rs.logger.Info("get note from redis")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return note.Note{}, storeErr
	}
	noteStr, err := rs.client.Get(strconv.Itoa(int(id))).Result()
	if err != nil {
		rs.logger.Debugf("error during getting note: %v", err)
		return note.Note{}, err
	}
	rs.logger.Debugf("unmarshal note %q", noteStr)
	var n note.Note
	err = json.Unmarshal([]byte(noteStr), &n)
	if err != nil {
		rs.logger.Debugf("error during unmarshaling note: %v", err)
		return note.Note{}, err
	}
	return n, err
}

func (rs *redisStorage) GetAll(ctx context.Context, login string) (note.Notes, error) {
	rs.logger.Info("get notes from redis")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return note.Notes{}, storeErr
	}
	rs.logger.Debugf("get aggregate by login %q", login)
	aggrStr, err := rs.client.Get(login).Result()
	if err != nil {
		rs.logger.Debugf("error during getting aggregate %v", err)
		return nil, err
	}
	rs.logger.Debugf("unmarshal aggregate %q", aggrStr)
	var aggr userAggregate
	err = json.Unmarshal([]byte(aggrStr), &aggr)
	if err != nil {
		rs.logger.Debugf("error during unmarshaling aggregate: %v", err)
		return note.Notes{}, err
	}
	rs.logger.Debugf("get notes by ids: %v", aggr)
	var ns []note.Note
	for i, nId := range aggr.NoteIds {
		n, err := rs.GetById(ctx, nId)
		rs.logger.Debugf("%d note.Id %d note %v", i, nId, n)
		if err != nil {
			return note.Notes{}, err
		}
		ns = append(ns, n)
	}
	return ns, err
}

func (rs *redisStorage) Update(ctx context.Context, n note.Note) error {
	rs.logger.Info("get notes from redis")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return storeErr
	}
	rs.logger.Debugf("marshaling note: %v", n)
	bytes, err := json.Marshal(n)
	if err != nil {
		rs.logger.Debugf("error during marshaling note: %v", err)
		return err
	}
	rs.logger.Debugf("save note to redis %v", n)
	if err = rs.client.Set(strconv.Itoa(int(n.Id)), bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving note: %v", err)
		return err
	}
	return nil
}

func (rs *redisStorage) Delete(ctx context.Context, id int) error {
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		storeErr := nerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return storeErr
	}
	rs.logger.Debugf("get note from redis by id %d", id)
	n, err := rs.GetById(ctx, id)
	if err != nil {
		rs.logger.Debugf("error during getting note: %v", err)
		return err
	}
	rs.logger.Debugf("got note %v", n)
	rs.logger.Debugf("get user's note aggregate by login %q", n.Author)
	aggrStr, err := rs.client.Get(n.Author).Result()
	if err != nil {
		rs.logger.Debugf("error during getting aggregate: %v", err)
		return err
	}
	rs.logger.Debugf("unmarshal aggregate %q", aggrStr)
	var aggr userAggregate
	err = json.Unmarshal([]byte(aggrStr), &aggr)
	if err != nil {
		return err
	}
	rs.logger.Debugf("delete note from aggregate %v", aggr)
	newAggr := userAggregate{
		Login:   n.Author,
		NoteIds: []int{},
	}
	rs.logger.Debugf("%v %d", aggr, n.Id)
	for _, nId := range aggr.NoteIds {
		rs.logger.Debugf("%v %d", newAggr, nId)
		if nId != n.Id {
			newAggr.NoteIds = append(newAggr.NoteIds, nId)
		}
	}
	rs.logger.Debugf("marshaling aggregate %v", newAggr)
	bytes, err := json.Marshal(newAggr)
	if err != nil {
		rs.logger.Debugf("error during marshaling aggregate: %v", err)
		return err
	}
	rs.logger.Debugf("save aggregate to redis %v", newAggr)
	if err = rs.client.Set(newAggr.Login, bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving aggregate: %v", err)
		return err
	}
	if err = rs.client.Del(string(rune(id))).Err(); err != nil {
		return err
	}
	return nil
}

func (rs *redisStorage) DeleteAll(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
