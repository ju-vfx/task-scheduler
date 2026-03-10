const TaskDetails = ({ item }: { item: any }) => {
  return (
    <>
      <tr>
        <td colSpan={5}>Command: "{item.task_command}"</td>
      </tr>
      {item.task_output != "" && (
        <tr>
          <td className="terminal-output" colSpan={5}>
            {item.task_output}
          </td>
        </tr>
      )}
    </>
  );
};

export default TaskDetails;
