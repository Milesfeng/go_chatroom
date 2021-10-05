# Go_Chatroom
使用Golang實作即時多人聊天室

##  使用技術&框架
- 使用[gorilla/websocket](https://www.gorillatoolkit.org/)技術，實作聊天室通訊
- 使用[Echo](https://echo.labstack.com/)框架
- 使用Docker架設 MySQL & 部署

##  部署
```
docker build -t go_chatroom . --no-cache
docker run -it --name web -p 5000:5000 --net=net -d go_chatroom
```

