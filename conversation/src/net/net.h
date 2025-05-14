#ifndef NET_H
#define NET_H

#include <sys/socket.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <fcntl.h>
#include <iostream>
#include <string>
#include <cstring>

#if defined(__APPLE__) || defined(__FreeBSD__)
#include <sys/event.h>
#define USE_KQUEUE
#else
#include <sys/epoll.h>
#define USE_EPOLL
#endif

#define MAX_EVENTS 10
#define PORT 8080

// 设置非阻塞模式
void setNonBlocking(int sockfd);

// 服务器主函数
int serveSocket();

#endif // NET_H