grpcurl \
  -plaintext \
  -proto proto/todo.proto \
  -d '{"userId": "aaa"}' \
  127.0.0.1:50051 \
  NoticeService.Connect
