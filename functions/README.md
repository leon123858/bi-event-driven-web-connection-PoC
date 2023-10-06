# deploy

```shell

gcloud functions deploy get-todo-list \
--gen2 \
--region=asia-east1 \
--runtime=go121 \
--source=./ \
--entry-point=GetTodoList \
--trigger-topic=get-todo-list

gcloud functions deploy add-todo-item \
--gen2 \
--region=asia-east1 \
--runtime=go121 \
--source=./ \
--entry-point=AddTodoItem \
--trigger-topic=add-todo-item

gcloud functions deploy remove-todo-item \
--gen2 \
--region=asia-east1 \
--runtime=go121 \
--source=./ \
--entry-point=RemoveTodoItem \
--trigger-topic=remove-todo-item

gcloud functions deploy update-todo-item \
--gen2 \
--region=asia-east1 \
--runtime=go121 \
--source=./ \
--entry-point=UpdateTodoItem \
--trigger-topic=update-todo-item

```
