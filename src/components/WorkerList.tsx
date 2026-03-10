const WorkerList = ({ allWorkers }: { allWorkers: any }) => {
  if (allWorkers.length < 1) {
    return <>No workers available</>;
  }
  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Host</th>
          <th scope="col">Status</th>
          <th scope="col">Last Seen</th>
          <th scope="col">Connected At</th>
        </tr>
      </thead>
      <tbody>
        {allWorkers.map((workerItem: any) => (
          <tr key={workerItem.id}>
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
