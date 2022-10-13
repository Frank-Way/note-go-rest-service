package storage

import (
	"context"
	"encoding/json"
	"github.com/Frank-Way/note-go-rest-service/internal/user"
	"github.com/Frank-Way/note-go-rest-service/internal/user/uerror"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

var _ user.Storage = &redisStorage{}

type redisStorage struct {
	client *redis.Client
	logger *logrus.Logger
}

const nextIdKey = ".nextId1"

func NewRedisStorage(host, port, password string, db int, logger *logrus.Logger) (user.Storage, error) {
	addr := net.JoinHostPort(host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		storeErr := uerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB: " + addr
		return nil, storeErr
	}
	return &redisStorage{
		client: client,
		logger: logger,
	}, nil
}

func (rs *redisStorage) Save(ctx context.Context, user user.User) (string, error) {
	rs.logger.Info("save user to redis storage")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		rs.logger.Debug("No connection to Redis DB")
		storeErr := uerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return "", storeErr
	}
	rs.logger.Debug("generate password hash")
	if err := user.GeneratePasswordHash(); err != nil {
		rs.logger.Debugf("error during password hashing: %v", err)
		return "", err
	}
	rs.logger.Debug("get new id from redis")
	nextIdStr, err := rs.client.Get(nextIdKey).Result()
	if err != nil {
		rs.logger.Debugf("error during getting next id: %v", err)
		rs.logger.Debug("set next id to 1")
		nextIdStr = "1"
		err = rs.client.Set(nextIdKey, "2", 0).Err()
		if err != nil {
			return "", err
		}
	}
	rs.logger.Debug("parse new id")
	nextId, err := strconv.Atoi(nextIdStr)
	if err != nil {
		rs.logger.Debugf("error during parsing next id: %v", err)
		return "", err
	}
	rs.logger.Debug("set new id to user")
	user.Id = int(nextId)
	rs.logger.Debug("incr id in redis")
	rs.client.Incr(nextIdKey)
	rs.logger.Debug("marshaling user")
	bytes, err := json.Marshal(user)
	if err != nil {
		rs.logger.Debugf("error during marshaling user: %v", err)
		return "", err
	}
	rs.logger.Debug("save user in redis")
	if err = rs.client.Set(user.Login, bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving user: %v", err)
		return "", err
	}
	return user.Login, nil
}

func (rs *redisStorage) GetByLogin(ctx context.Context, login string) (user.User, error) {
	rs.logger.Info("get user from redis")
	rs.logger.Tracef("get user by login: %s", login)
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		rs.logger.Debug("No connection to Redis DB")
		storeErr := uerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return user.User{}, storeErr
	}
	rs.logger.Debug("check if user exists")
	u, err := rs.findUserByLogin(login)
	if err != nil {
		rs.logger.Debugf("error during getting user by login: %v", err)
		sErr := uerror.ErrorNotFound
		sErr.Err = err
		return user.User{}, sErr
	}
	rs.logger.Tracef("user: %v", u)
	return u, nil
}

func (rs *redisStorage) GetById(ctx context.Context, id int) (user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (rs *redisStorage) GetAll(ctx context.Context) (user.Users, error) {
	//TODO implement me
	panic("implement me")
}

func (rs *redisStorage) Update(ctx context.Context, user user.User) error {
	rs.logger.Info("update user in redis")
	rs.logger.Debug("check if redis available")
	if err := rs.client.Ping().Err(); err != nil {
		rs.logger.Debug("No connection to Redis DB")
		storeErr := uerror.ErrorStorage
		storeErr.DeveloperMessage = "No connection to Redis DB"
		return storeErr
	}
	rs.logger.Debug("get user")
	_, err := rs.findUserByLogin(user.Login)
	if err != nil {
		rs.logger.Debugf("error during getting user from redis: %v", err)
		return err
	}
	rs.logger.Debug("marshaling user")
	bytes, err := json.Marshal(user)
	if err != nil {
		rs.logger.Debugf("error during marshaling user: %v", err)
		return err
	}
	rs.logger.Debug("save user in redis")
	if err := rs.client.Set(user.Login, bytes, 0).Err(); err != nil {
		rs.logger.Debugf("error during saving user: %v", err)
		return err
	}
	return nil
}

func (rs *redisStorage) DeleteByLogin(ctx context.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (rs *redisStorage) DeleteById(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (rs *redisStorage) findUserByLogin(login string) (user.User, error) {
	rs.logger.Tracef("find user by login: %s", login)
	uStr, err := rs.client.Get(login).Result()
	if err != nil {
		rs.logger.Tracef("error in redis: %s", login)
		sErr := uerror.ErrorStorage
		sErr.Err = err
		return user.User{}, sErr
	}
	var u user.User
	err = json.Unmarshal([]byte(uStr), &u)
	if err != nil {
		sErr := uerror.ErrorStorage
		sErr.Err = err
		return user.User{}, sErr
	}
	rs.logger.Tracef("unmarshaled user: %v", u)
	return u, err
}
