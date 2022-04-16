import { useCallback, useEffect, useRef, useState } from "react";
import prettyBytes from "pretty-bytes";

const Requests = () => {

  const [requests, setRequests] = useState([]);

  const [showMostRecent, setShowMostRecent] = useState(10);
  const [showOnlyText, setShowOnlyText] = useState(false);

  const [requestsServed, setRequestsServed] = useState(0);
  const [bytesSent, setBytesSent] = useState(0);

  let logRequest = useCallback((r) => {
    const date = new Date();
    const time =
      <>
        {date.getHours().toString().padStart(2, '0')}:
        {date.getMinutes().toString().padStart(2, '0')}:
        {date.getSeconds().toString().padStart(2, '0')}
        <small className="text-muted text-small">
          .{date.getMilliseconds().toString().padStart(3, '0')}
        </small>
      </>;
    setRequests(requests => [
      {
        ...r,
        date,
        time
      },
      ...requests
    ].slice(0, showMostRecent || undefined));
  }, [showMostRecent]);

  const websocket = useRef(null);

  useEffect(() => {
    const loc = window.location;
    let hostname = process.env.REACT_APP_HOSTNAME || loc.host;
    const ws = new WebSocket(
      `${loc.protocol === "https:" ? 'wss:' : 'ws:'}${hostname}/api/ws`
    );
    websocket.current = ws;
    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    websocket.current.onmessage = (e) => {
      const r = JSON.parse(e.data).data;
      if (!showOnlyText || r.content_type.startsWith("text/")) {
        logRequest(r);
      }
      setRequestsServed(requestsServed => requestsServed + 1);
      let requestBytes = parseInt(r.content_length);
      if (requestBytes > 0) {
        setBytesSent(bytesSent => bytesSent + requestBytes);
      }
    };
  }, [showOnlyText, logRequest]);

  function statusCodeColor(statusCode) {
    if (statusCode >= 200 && statusCode < 300) {
      return 'text-success';
    } else if (statusCode < 400) {
      return 'text-info';
    } else {
      return 'text-danger';
    }
  }

  function prettySize(size) {
    try {
      return prettyBytes(parseInt(size));
    } catch {
      return "-";
    }
  }

  function handleShowMostRecentChanged(e) {
    const v = e.target.value;
    setRequests(requests.slice(0, v || undefined));
    setShowMostRecent(v);
  }

  function handleShowOnlyTextChanged(e) {
    const v = e.target.checked;
    if (v) {
      setRequests(requests.filter(r => {
        return r.content_type.startsWith("text/");
      }));
    }
    setShowOnlyText(v);
  }

  return (
    <div>
      <h1 className="fw-normal">Requests</h1>
      <p>
        The table below shows requests in realtime as they come in.
      </p>
      <div className="row">
        <div className="col-10">
          <table className="table table-striped">
            <thead className="table-dark">
              <tr>
                <th>Time</th>
                <th>Client</th>
                <th>Request</th>
                <th>Response</th>
                <th>Type</th>
                <th className="text-end">Size</th>
              </tr>
            </thead>
            <tbody>
              {requests.length ?
                requests.map(request => (
                  <tr key={+ request.date}>
                    <td className="collapse-column">{request.time}</td>
                    <td className="collapse-column">
                      <img
                        src={`/static/img/${request.country_code}.png`}
                        className="flag"
                        title={request.country_name}
                        alt={request.country_name} />
                      {request.remote_addr}
                    </td>
                    <td className="truncate-column" style={{ maxWidth: '200px' }}>
                      <div className="truncate-fade">
                        <span className="badge bg-secondary">{request.method}</span>{' '}
                        <strong>{request.host}</strong>{request.path}
                      </div>
                    </td>
                    <td className={`collapse-column ${statusCodeColor(request.status_code)}`}>
                      <strong>{request.status}</strong>
                    </td>
                    <td className="collapse-column">{request.content_type || "-"}</td>
                    <td className="collapse-column text-end">{prettySize(request.content_length)}</td>
                  </tr>
                )) :
                <tr>
                  <td colSpan="6">
                    <p className="text-muted text-center p-5">No requests yet</p>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
        <div className="col-2">
          <div className="card">
            <div className="card-header">
              Display
            </div>
            <div className="card-body">
              <div className="mb-3">
                <label htmlFor="request_limit" className="form-label">Show most recent:</label>
                <select
                  className="form-select"
                  id="request_limit"
                  value={showMostRecent}
                  onChange={handleShowMostRecentChanged}>
                  <option value="0">all</option>
                  <option value="5">5</option>
                  <option value="10">10</option>
                  <option value="15">15</option>
                  <option value="20">20</option>
                </select>
              </div>
              <div className="mb-3">
                <div className="form-check">
                  <input
                    className="form-check-input"
                    id="text_only"
                    type="checkbox"
                    checked={showOnlyText}
                    onChange={handleShowOnlyTextChanged} />
                  <label className="form-check-label" htmlFor="text_only">
                    text/html only
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <br />
      <div className="row">
        <div className="col-3">
          <h1 className="card-title text-center fw-normal">
            {requestsServed}
            <small className="text-muted stat-label">request(s)</small>
          </h1>
        </div>
        <div className="col-3">
          <h1 className="card-title text-center fw-normal">
            {prettySize(bytesSent)}
            <small className="text-muted stat-label">content sent</small>
          </h1>
        </div>
        <div className="col-3"></div>
        <div className="col-3"></div>
      </div>
    </div>
  );
};

export default Requests;
