const WorkerList = ({ allWorkers }) => {
  if (allWorkers.length < 1) {
    return <>No workers available</>;
  }
  return (
    <table className="table">
      <thead>
        <tr href="#">
          <th scope="col">Host</th>
          <th scope="col">Status</th>
          <th scope="col">Last Seen</th>
          <th scope="col">Connected At</th>
        </tr>
      </thead>
      <tbody>
        {allWorkers.map((workerItem) => (
          <tr>
            <td>{workerItem.host}</td>
            <td>{workerItem.status}</td>
            <td>{workerItem.last_seen_at}</td>
            <td>{workerItem.connected_at}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default WorkerList;
