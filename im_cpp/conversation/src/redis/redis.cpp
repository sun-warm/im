#include <iostream>
#include <hiredis/hiredis.h>
#include "redis.h"
redisContext* context = NULL;
int initRedis() {
    // 连接到 Redis 服务器
    redisContext* context = redisConnect("127.0.0.1", 6379);
    if (context == NULL || context->err) {
        if (context) {
            std::cerr << "Error: " << context->errstr << std::endl;
            redisFree(context);
        } else {
            std::cerr << "Can't allocate redis context" << std::endl;
        }
        return 1;
    }

    // 设置键值对
    redisReply* reply = (redisReply*)redisCommand(context, "SET %s %s", "key", "value");
    if (reply == NULL) {
        std::cerr << "SET command failed" << std::endl;
        redisFree(context);
        return 1;
    }
    std::cout << "SET: " << reply->str << std::endl;
    freeReplyObject(reply);

    // 获取键值对
    reply = (redisReply*)redisCommand(context, "GET %s", "key");
    if (reply == NULL) {
        std::cerr << "GET command failed" << std::endl;
        redisFree(context);
        return 1;
    }
    if (reply->type == REDIS_REPLY_STRING) {
        std::cout << "GET: " << reply->str << std::endl;
    } else {
        std::cout << "GET: key not found" << std::endl;
    }
    freeReplyObject(reply);

    // 关闭连接
    redisFree(context);
    return 0;
}