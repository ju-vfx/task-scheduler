import TaskDetails from "./TaskDetails";

const TaskItem = ({
  item,
  onSelectTask,
  selectedTask,
}: {
  item: any;
  onSelectTask: Function;
  selectedTask: any;
}) => {
  let className = "";
  if (item.task_status === "Finished") {
    className = "status-finished";
  } else if (item.task_status === "Error") {
    className = "status-error";
  }

  return (
    <>
      <tr
        role="button"
        onClick={() => {
          onSelectTask(item.task_id);
        }}
        className={selectedTask === item.task_id ? "table-primary" : ""}
      >
        <td>{item.task_name}</td>
        <td className={className}>{item.task_status}</td>
        <td>{item.task_created_at}</td>
        <td>{item.task_finished_at}</td>
      </tr>
      {selectedTask === item.task_id && <TaskDetails item={item} />}
    </>
  );
};

export default TaskItem;
