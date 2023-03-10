/*
 * @Description:
 * @Author: mali
 * @Date: 2023-03-02 13:47:05
 * @LastEditTime: 2023-03-10 14:41:11
 * @LastEditors: VSCode
 * @Reference:
 */
package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/laravelGo/core/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var once sync.Once

var Client *MongoClient

type MongoClient struct {
	Client  *mongo.Database
	Context context.Context
}

/**
 * @Author: mali
 * @Func:
 * @Description: 连接mongo
 * @Param:
 * @Return:
 * @Example:
 * @param {string} host
 * @param {uint64} port
 * @param {string} username
 * @param {string} password
 * @param {string} dbName
 * @param {*} timeOut
 * @param {uint64} maxNum
 */
func ConnectMongo(host string, port uint64, username string, password string, dbName string, timeOut, maxNum uint64) {
	once.Do(func() {
		Client = NewClient(host, port, username, password, dbName, timeOut, maxNum)
	})
}

/**
 * @Author: mali
 * @Func:
 * @Description: 新建客户端
 * @Param:
 * @Return:
 * @Example:
 * @param {string} host
 * @param {uint64} port
 * @param {string} username
 * @param {string} password
 * @param {string} dbName
 * @param {*} timeOut
 * @param {uint64} maxNum
 */
func NewClient(host string, port uint64, username string, password string, dbName string, timeOut, maxNum uint64) *MongoClient {
	clientOptions := options.Client()
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, host, port)
	clientOptions.ApplyURI(url)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut))
	defer cancel()
	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	clientOptions.SetMaxPoolSize(maxNum)
	// 发起链接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.ErrorString("ConnectToMongoDB", "connect", err.Error())

	}
	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		logger.ErrorString("ConnectToMongoDB", "ping检测", err.Error())
	}
	// 返回 client
	return &MongoClient{
		Client:  client.Database(dbName),
		Context: ctx,
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 插入一条数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} document
 * @param {...*options.InsertOneOptions} opts
 */
func (c *MongoClient) InsertOne(ctx context.Context, collection string, document interface{}, opts ...*options.InsertOneOptions) (interface{}, bool) {
	insertResult, err := c.Client.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		logger.ErrorString("Mongo", "InsertOne", err.Error())
		return nil, false
	}
	return insertResult.InsertedID, true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 批量插入出入数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {[]interface{}} document
 */
func (c *MongoClient) InsertMany(ctx context.Context, collection string, document []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, bool) {
	insertResult, err := c.Client.Collection(collection).InsertMany(ctx, document, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "InsertMany", err.Error())
		return nil, false
	}
	return insertResult.InsertedIDs, true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 删除单条数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} filter
 * @param {...*options.DeleteOptions} opts
 */
func (c *MongoClient) DeleteOne(ctx context.Context, collection string, filter interface{}, opts ...*options.DeleteOptions) (int64, bool) {
	deleteResult, err := c.Client.Collection(collection).DeleteOne(ctx, filter, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "DeleteOne", err.Error())
		return 0, false
	}
	return deleteResult.DeletedCount, true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 批量删除
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} document
 * @param {...*options.DeleteOptions} opts
 */
func (c *MongoClient) DeleteMany(ctx context.Context, collection string, document interface{}, opts ...*options.DeleteOptions) (int64, bool) {
	deleteResult, err := c.Client.Collection(collection).DeleteMany(ctx, document, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "DeleteMany", err.Error())
		return 0, false
	}
	return deleteResult.DeletedCount, true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 更新用户数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} id
 * @param {interface{}} update
 * @param {...*options.UpdateOptions} opts
 */
func (c *MongoClient) UpdateByID(ctx context.Context, collection string, id interface{}, update interface{}, opts ...*options.UpdateOptions) bool {
	updateResult, err := c.Client.Collection(collection).UpdateByID(ctx, id, update, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "UpdateByID", err.Error())
		return false
	}
	//未匹配到更新元素
	if updateResult.MatchedCount == 0 {
		logger.ErrorString("Mongo", "UpdateByID", "未匹配到更新元素")
		return false
	}
	return true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 更新单条数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} filter
 * @param {interface{}} update
 * @param {...*options.UpdateOptions} opts
 */
func (c *MongoClient) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) bool {
	updateResult, err := c.Client.Collection(collection).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "UpdateOne", err.Error())
		return false
	}
	//未匹配到更新元素
	if updateResult.MatchedCount == 0 {
		logger.ErrorString("Mongo", "UpdateOne", "未匹配到更新元素")
		return false
	}
	return true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 批量更新
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} filter
 * @param {interface{}} update
 * @param {...*options.UpdateOptions} opts
 */
func (c *MongoClient) UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) bool {
	updateResult, err := c.Client.Collection(collection).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "UpdateMany", err.Error())
		return false
	}
	//未匹配到更新元素
	if updateResult.MatchedCount == 0 {
		logger.ErrorString("Mongo", "UpdateMany", "未匹配到更新元素")
		return false
	}
	return true
}

/**
 * @Author: mali
 * @Func:
 * @Description: 查询数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} filter
 * @param {...*options.FindOptions} opts
 */
func (c *MongoClient) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) []interface{} {
	findResult, err := c.Client.Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		logger.ErrorString("Mongo", "Find", err.Error())
	}
	results := []interface{}{}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for findResult.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem interface{}
		err = findResult.Decode(&elem)
		if err != nil {
			logger.ErrorString("Mongo", "Find-result-decode", err.Error())
		}

		results = append(results, &elem)
	}

	if err = findResult.Err(); err != nil {
		logger.ErrorString("Mongo", "Find-result-err", err.Error())
	}

	// Close the cursor once finished
	findResult.Close(context.TODO())
	return results
}

/**
 * @Author: mali
 * @Func:
 * @Description: 查询单挑数据
 * @Param:
 * @Return:
 * @Example:
 * @param {context.Context} ctx
 * @param {string} collection
 * @param {interface{}} filter
 * @param {...*options.FindOneOptions} opts
 */
func (c *MongoClient) FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) interface{} {
	findResult := c.Client.Collection(collection).FindOne(ctx, filter, opts...)
	var result interface{}
	findResult.Decode(result)
	return result
}
