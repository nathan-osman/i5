import { useEffect, useState } from "react";

const RequestTicker = () => {

  const [request, setRequests] = useState([]);

  useEffect(() => {
    const loc = window.location;
    let hostname = process.env.REACT_APP_HOSTNAME || loc.host;
    let websocket = new WebSocket(
      `${loc.protocol === "https:" ? 'wss:' : 'ws:'}${hostname}/api/ws`
    );
    websocket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setRequests(r => [
        ...r,
        {
          ...data.data,
          time: new Date()
        }
      ]);
    };
    return () => {
      websocket.close();
    };
  }, []);

  return (
    <div>
      <h1>Requests</h1>
      <p>
        The table below shows requests in realtime as they come in.
      </p>

      <table className="table table-striped">
        <thead>
          <tr>
            <th>Time</th>
            <th>Client</th>
            <th>Method</th>
            <th>Path</th>
          </tr>
        </thead>
        <tbody>
          {request.map(request => (
            <tr key={+ request.time}>
              <td>{request.time.toISOString()}</td>
              <td>{request.remote_addr}</td>
              <td>
                <span className="badge bg-secondary">{request.method}</span></td>
              <td><strong>{request.host}</strong>{request.path}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default RequestTicker;
