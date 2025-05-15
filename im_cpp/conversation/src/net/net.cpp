#include "src/conversation/conversaion.h"
#include "net.h"

void setNonBlocking(int sockfd) {
    int opts = fcntl(sockfd, F_GETFL);
    if (opts < 0) {
        perror("fcntl(F_GETFL)");
        exit(EXIT_FAILURE);
    }
    opts = (opts | O_NONBLOCK);
    if (fcntl(sockfd, F_SETFL, opts) < 0) {
        perror("fcntl(F_SETFL)");
        exit(EXIT_FAILURE);
    }
}

int serveSocket() {
    int listen_fd, conn_fd;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

#ifdef USE_KQUEUE
    int kq, nev;
    struct kevent evSet, evList[MAX_EVENTS];

    // 创建 kqueue 实例
    kq = kqueue();
    if (kq == -1) {
        perror("kqueue");
        exit(EXIT_FAILURE);
    }
#else
    int epoll_fd, nfds;
    struct epoll_event ev, events[MAX_EVENTS];

    // 创建 epoll 实例
    epoll_fd = epoll_create1(0);
    if (epoll_fd < 0) {
        perror("epoll_create1");
        exit(EXIT_FAILURE);
    }
#endif

    // 创建监听套接字
    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd < 0) {
        perror("socket");
        exit(EXIT_FAILURE);
    }

    // 设置非阻塞模式
    setNonBlocking(listen_fd);

    // 绑定地址和端口
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    if (bind(listen_fd, (struct sockaddr*)&server_addr, sizeof(server_addr)) < 0) {
        perror("bind");
        close(listen_fd);
        exit(EXIT_FAILURE);
    }

    // 监听端口
    if (listen(listen_fd, SOMAXCONN) < 0) {
        perror("listen");
        close(listen_fd);
        exit(EXIT_FAILURE);
    }

#ifdef USE_KQUEUE
    // 添加监听套接字到 kqueue 实例
    EV_SET(&evSet, listen_fd, EVFILT_READ, EV_ADD, 0, 0, NULL);
    if (kevent(kq, &evSet, 1, NULL, 0, NULL) == -1) {
        perror("kevent: listen_fd");
        close(listen_fd);
        close(kq);
        exit(EXIT_FAILURE);
    }
#else
    // 添加监听套接字到 epoll 实例
    ev.events = EPOLLIN;
    ev.data.fd = listen_fd;
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, listen_fd, &ev) < 0) {
        perror("epoll_ctl: listen_fd");
        close(listen_fd);
        close(epoll_fd);
        exit(EXIT_FAILURE);
    }
#endif

    while (true) {
#ifdef USE_KQUEUE
        nev = kevent(kq, NULL, 0, evList, MAX_EVENTS, NULL);
        if (nev < 0) {
            perror("kevent");
            close(listen_fd);
            close(kq);
            exit(EXIT_FAILURE);
        }

        for (int i = 0; i < nev; ++i) {
            if (evList[i].ident == listen_fd) {
                // 处理新连接
                conn_fd = accept(listen_fd, (struct sockaddr*)&client_addr, &client_addr_len);
                if (conn_fd < 0) {
                    perror("accept");
                    continue;
                }
                setNonBlocking(conn_fd);
                EV_SET(&evSet, conn_fd, EVFILT_READ, EV_ADD, 0, 0, NULL);
                if (kevent(kq, &evSet, 1, NULL, 0, NULL) == -1) {
                    perror("kevent: conn_fd");
                    close(conn_fd);
                    continue;
                }
                std::cout << "New connection accepted" << std::endl;
            } else {
                // 处理客户端数据
                SocketServer server;
                if (!server.handleRequest(evList[i].ident)) {
                    close(evList[i].ident);
                }
            }
        }
#else
        nfds = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);
        if (nfds < 0) {
            perror("epoll_wait");
            close(listen_fd);
            close(epoll_fd);
            exit(EXIT_FAILURE);
        }

        for (int i = 0; i < nfds; ++i) {
            if (events[i].data.fd == listen_fd) {
                // 处理新连接
                conn_fd = accept(listen_fd, (struct sockaddr*)&client_addr, &client_addr_len);
                if (conn_fd < 0) {
                    perror("accept");
                    continue;
                }
                setNonBlocking(conn_fd);
                ev.events = EPOLLIN | EPOLLET;
                ev.data.fd = conn_fd;
                if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, conn_fd, &ev) < 0) {
                    perror("epoll_ctl: conn_fd");
                    close(conn_fd);
                    continue;
                }
                std::cout << "New connection accepted" << std::endl;
            } else {
                // 处理客户端数据
                SocketServer server;
                if (!server.handleRequest(events[i].data.fd)) {
                    close(events[i].data.fd);
                }
            }
        }
#endif
    }

#ifdef USE_KQUEUE
    close(kq);
#else
    close(epoll_fd);
#endif
    close(listen_fd);
    return 0;
}