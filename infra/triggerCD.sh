# 定义存储触发器名称的列表
triggers=("create-todolist-trigger" "delete-todolist-trigger" "get-todolist-trigger" "http-trigger" "update-todolist-trigger" "web-trigger" "websocket-trigger")

# 循环触发每个触发器并行执行
for trigger in "${triggers[@]}"; do
  gcloud builds triggers run $trigger --region=asia-east1 --branch=main &
done

# 等待所有后台作业完成
wait
